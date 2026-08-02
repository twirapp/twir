import type { StreamElementsIntegration } from '../libs/db.ts'

export interface StreamElementsConnectionHandle {
	destroy(): void | Promise<void>
}

export interface StreamElementsStoreOptions {
	readonly loadIntegrationByID?: (
		integrationID: string
	) => StreamElementsIntegration | null | Promise<StreamElementsIntegration | null>
	readonly createConnection: (
		integration: StreamElementsIntegration
	) => StreamElementsConnectionHandle | Promise<StreamElementsConnectionHandle>
}

export interface StreamElementsStore {
	readonly connections: ReadonlyMap<string, StreamElementsConnectionHandle>
	addIntegration(integration: StreamElementsIntegration): Promise<void>
	addIntegrationByID(integrationID: string): Promise<void>
	removeIntegration(channelID: string): Promise<void>
	closeAll(): Promise<void>
}

export function runLifecycleOperation(
	operation: () => Promise<void>,
	onError: (error: unknown) => void = (error) => {
		console.error('Integration lifecycle operation failed', error)
	},
): void {
	void operation().catch(onError)
}

export function createStreamElementsStore(options: StreamElementsStoreOptions): StreamElementsStore {
	const connections = new Map<string, StreamElementsConnectionHandle>()
	const operations = new Map<string, Promise<void>>()
	let closed = false

	const enqueue = (channelID: string, operation: () => Promise<void>): Promise<void> => {
		const previous = operations.get(channelID) ?? Promise.resolve()
		const next = previous.catch(() => undefined).then(operation)
		operations.set(channelID, next)
		return next.finally(() => {
			if (operations.get(channelID) === next) operations.delete(channelID)
		})
	}

	const removeConnection = async (channelID: string): Promise<void> => {
		const existing = connections.get(channelID)
		if (!existing) return
		connections.delete(channelID)
		await existing.destroy()
	}

	return {
		connections,
		async addIntegration(integration) {
			if (closed || !integration.enabled || !integration.accessToken || !integration.refreshToken) return
			await enqueue(integration.channelId, async () => {
				if (closed) return
				await removeConnection(integration.channelId)
				const connection = await options.createConnection(integration)
				connections.set(integration.channelId, connection)
			})
		},
		async addIntegrationByID(integrationID) {
			const loader = options.loadIntegrationByID
			if (!loader) {
				throw new Error('StreamElements integration loader is not configured')
			}
			if (closed) return
			const discovered = await loader(integrationID)
			if (closed || !discovered) return
			await enqueue(discovered.channelId, async () => {
				if (closed) return
				const authoritative = await loader(integrationID)
				if (!authoritative
					|| authoritative.channelId !== discovered.channelId
					|| !authoritative.enabled
					|| !authoritative.accessToken
					|| !authoritative.refreshToken) {
					await removeConnection(discovered.channelId)
					return
				}
				await removeConnection(authoritative.channelId)
				const connection = await options.createConnection(authoritative)
				connections.set(authoritative.channelId, connection)
			})
		},
		async removeIntegration(channelID) {
			await enqueue(channelID, () => removeConnection(channelID))
		},
		async closeAll() {
			closed = true
			const channelIDs = new Set([...connections.keys(), ...operations.keys()])
			await Promise.all([...channelIDs].map((channelID) =>
				enqueue(channelID, () => removeConnection(channelID))
			))
		},
	}
}
