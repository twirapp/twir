import type { StreamlabsIntegration } from '../libs/db.ts'

export interface StreamLabsConnectionHandle {
	destroy(): void | Promise<void>
}

export interface StreamLabsStoreOptions {
	readonly loadIntegrationByID?: (
		integrationID: string
	) => StreamlabsIntegration | null | Promise<StreamlabsIntegration | null>
	readonly createConnection: (
		integration: StreamlabsIntegration
	) => StreamLabsConnectionHandle | Promise<StreamLabsConnectionHandle>
}

export interface StreamLabsStore {
	readonly connections: ReadonlyMap<string, StreamLabsConnectionHandle>
	addIntegration(integration: StreamlabsIntegration): Promise<void>
	addIntegrationByID(integrationID: string): Promise<void>
	removeIntegration(channelID: string): Promise<void>
	closeAll(): Promise<void>
}

export function createStreamLabsStore(options: StreamLabsStoreOptions): StreamLabsStore {
	const connections = new Map<string, StreamLabsConnectionHandle>()
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
			if (closed
				|| !integration.enabled
				|| !integration.access_token
				|| !integration.refresh_token) {
				return
			}
			await enqueue(integration.channel_id, async () => {
				if (closed) return
				await removeConnection(integration.channel_id)
				const connection = await options.createConnection(integration)
				connections.set(integration.channel_id, connection)
			})
		},
		async addIntegrationByID(integrationID) {
			const loader = options.loadIntegrationByID
			if (!loader) throw new Error('Streamlabs integration loader is not configured')
			if (closed) return
			const discovered = await loader(integrationID)
			if (closed || !discovered) return
			await enqueue(discovered.channel_id, async () => {
				if (closed) return
				const authoritative = await loader(integrationID)
				if (!authoritative
					|| authoritative.channel_id !== discovered.channel_id
					|| !authoritative.enabled
					|| !authoritative.access_token
					|| !authoritative.refresh_token) {
					await removeConnection(discovered.channel_id)
					return
				}
				await removeConnection(authoritative.channel_id)
				const connection = await options.createConnection(authoritative)
				connections.set(authoritative.channel_id, connection)
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
