import { expect, test } from 'bun:test'

import type { StreamlabsIntegration } from '../libs/db.ts'
import type { StreamLabsConnectionHandle } from './streamlabs.ts'
import { createStreamLabsStore } from './streamlabs.ts'

function integration(overrides: Partial<StreamlabsIntegration> = {}): StreamlabsIntegration {
	return {
		id: 'integration-1',
		enabled: true,
		channel_id: 'channel-1',
		access_token: 'access',
		refresh_token: 'refresh',
		username: 'streamer',
		avatar: 'avatar',
		...overrides,
	}
}

test('keeps one connection per channel and closes the replaced connection', async () => {
	const destroyed: string[] = []
	let created = 0
	const store = createStreamLabsStore({
		createConnection: async () => {
			created += 1
			const id = `connection-${created}`
			return { destroy: () => { destroyed.push(id) } }
		},
	})

	await store.addIntegration(integration())
	await store.addIntegration(integration({ id: 'integration-2', access_token: 'next' }))

	expect(store.connections.size).toBe(1)
	expect(destroyed).toEqual(['connection-1'])
})

test('serializes concurrent replacements and retains only the newest socket', async () => {
	let releaseFirst: (() => void) | undefined
	const firstBlocked = new Promise<void>((resolve) => { releaseFirst = resolve })
	const handles: Array<StreamLabsConnectionHandle & { name: string; destroyed: boolean }> = []
	const store = createStreamLabsStore({
		createConnection: async (record) => {
			if (record.id === 'first') await firstBlocked
			const handle = {
				name: record.id,
				destroyed: false,
				destroy() { this.destroyed = true },
			}
			handles.push(handle)
			return handle
		},
	})

	const first = store.addIntegration(integration({ id: 'first' }))
	const second = store.addIntegration(integration({ id: 'second', access_token: 'next' }))
	releaseFirst?.()
	await Promise.all([first, second])

	expect(handles.map((handle) => [handle.name, handle.destroyed])).toEqual([
		['first', true],
		['second', false],
	])
})

test('ignores disabled records and provider auth failures never disable integration state', async () => {
	const enabled = integration()
	const store = createStreamLabsStore({
		createConnection: async () => { throw new Error('provider auth failed') },
	})

	await store.addIntegration(integration({ enabled: false }))
	expect(store.connections.size).toBe(0)
	await expect(store.addIntegration(enabled)).rejects.toThrow('provider auth failed')
	expect(enabled.enabled).toBe(true)
	expect(store.connections.size).toBe(0)
})

test('removal destroys the connection and closeAll drains every channel', async () => {
	let destroyed = 0
	const store = createStreamLabsStore({
		createConnection: async () => ({ destroy: () => { destroyed += 1 } }),
	})

	await store.addIntegration(integration())
	await store.removeIntegration('channel-1')
	expect(destroyed).toBe(1)

	await store.addIntegration(integration())
	await store.addIntegration(integration({ id: 'integration-2', channel_id: 'channel-2' }))
	await store.closeAll()
	expect(destroyed).toBe(3)
	expect(store.connections.size).toBe(0)
})

test('rereads authoritative state inside queue so stale Add cannot undo Remove', async () => {
	let resolveInitial: ((record: StreamlabsIntegration) => void) | undefined
	const initial = new Promise<StreamlabsIntegration>((resolve) => { resolveInitial = resolve })
	let loads = 0
	let created = 0
	const store = createStreamLabsStore({
		async loadIntegrationByID() {
			loads += 1
			if (loads === 1) return initial
			return null
		},
		createConnection: async () => {
			created += 1
			return { destroy: () => undefined }
		},
	})

	const staleAdd = store.addIntegrationByID('integration-1')
	await store.removeIntegration('channel-1')
	resolveInitial?.(integration())
	await staleAdd

	expect(loads).toBe(2)
	expect(created).toBe(0)
	expect(store.connections.size).toBe(0)
})

test('shutdown prevents delayed discovery from creating a connection', async () => {
	let resolveDiscovery: ((record: StreamlabsIntegration) => void) | undefined
	const discovery = new Promise<StreamlabsIntegration>((resolve) => { resolveDiscovery = resolve })
	let created = 0
	const store = createStreamLabsStore({
		loadIntegrationByID: async () => discovery,
		createConnection: async () => {
			created += 1
			return { destroy: () => undefined }
		},
	})

	const add = store.addIntegrationByID('integration-1')
	await store.closeAll()
	resolveDiscovery?.(integration())
	await add
	expect(created).toBe(0)
})
