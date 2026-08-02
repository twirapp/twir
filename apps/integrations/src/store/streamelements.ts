import type { StreamElementsIntegration } from '../libs/db.ts'

export interface StreamElementsConnectionHandle {
	destroy(): void | Promise<void>
}

export interface StreamElementsStoreOptions {
	readonly createConnection: (
		integration: StreamElementsIntegration
	) => StreamElementsConnectionHandle | Promise<StreamElementsConnectionHandle>
}

export interface StreamElementsStore {
	readonly connections: ReadonlyMap<string, StreamElementsConnectionHandle>
	addIntegration(integration: StreamElementsIntegration): Promise<void>
	removeIntegration(channelID: string): Promise<void>
	closeAll(): Promise<void>
}

export function createStreamElementsStore(options: StreamElementsStoreOptions): StreamElementsStore {
	const connections = new Map<string, StreamElementsConnectionHandle>()
	const operations = new Map<string, Promise<void>>()

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
			if (!integration.enabled || !integration.accessToken || !integration.refreshToken) return
			await enqueue(integration.channelId, async () => {
				await removeConnection(integration.channelId)
				const connection = await options.createConnection(integration)
				connections.set(integration.channelId, connection)
			})
		},
		async removeIntegration(channelID) {
			await enqueue(channelID, () => removeConnection(channelID))
		},
		async closeAll() {
			const channelIDs = new Set([...connections.keys(), ...operations.keys()])
			await Promise.all([...channelIDs].map((channelID) =>
				enqueue(channelID, () => removeConnection(channelID))
			))
		},
	}
}
