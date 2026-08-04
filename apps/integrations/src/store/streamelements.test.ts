import { expect, test } from 'bun:test'

import type { StreamElementsIntegration } from '../libs/db.ts'
import type { StreamElementsConnectionHandle } from './streamelements.ts'
import { createStreamElementsStore, runLifecycleOperation } from './streamelements.ts'

function integration(overrides: Partial<StreamElementsIntegration> = {}): StreamElementsIntegration {
	return {
		id: 'integration-1',
		channelId: 'channel-1',
		enabled: true,
		accessToken: 'access',
		refreshToken: 'refresh',
		...overrides,
	}
}

test('adds one connection per channel and replaces the previous connection', async () => {
	const destroyed: string[] = []
	let created = 0
	const store = createStreamElementsStore({
		createConnection: async () => {
			created += 1
			const id = `connection-${created}`
			return { destroy: () => { destroyed.push(id) } }
		},
	})

	await store.addIntegration(integration())
	await store.addIntegration(integration({ id: 'integration-2', accessToken: 'next' }))

	expect(store.connections.size).toBe(1)
	expect(created).toBe(2)
	expect(destroyed).toEqual(['connection-1'])
})

test('serializes concurrent replacements and leaves only the newest connection live', async () => {
	let releaseFirst: (() => void) | undefined
	const firstBlocked = new Promise<void>((resolve) => { releaseFirst = resolve })
	const handles: Array<StreamElementsConnectionHandle & { name: string; destroyed: boolean }> = []
	const store = createStreamElementsStore({
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
	const second = store.addIntegration(integration({ id: 'second', accessToken: 'next' }))
	releaseFirst?.()
	await Promise.all([first, second])

	expect(handles.map((handle) => [handle.name, handle.destroyed])).toEqual([
		['first', true],
		['second', false],
	])
	expect(store.connections.size).toBe(1)
})

test('ignores disabled or incomplete records, removes by channel, and closes all', async () => {
	let destroyed = 0
	const store = createStreamElementsStore({
		createConnection: async () => ({ destroy: () => { destroyed += 1 } }),
	})

	await store.addIntegration(integration({ enabled: false }))
	await store.addIntegration(integration({ channelId: 'channel-2', accessToken: null }))
	expect(store.connections.size).toBe(0)

	await store.addIntegration(integration())
	await store.removeIntegration('channel-1')
	expect(destroyed).toBe(1)
	expect(store.connections.size).toBe(0)

	await store.addIntegration(integration({ channelId: 'channel-1' }))
	await store.addIntegration(integration({ id: 'integration-2', channelId: 'channel-2' }))
	await store.closeAll()
	expect(destroyed).toBe(3)
	expect(store.connections.size).toBe(0)
})

test('rechecks authoritative state inside the channel queue before a stale Add can recreate', async () => {
	let resolveInitial: ((record: StreamElementsIntegration) => void) | undefined
	const initial = new Promise<StreamElementsIntegration>((resolve) => { resolveInitial = resolve })
	let loads = 0
	let created = 0
	const store = createStreamElementsStore({
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

test('lifecycle callback wrapper handles rejected operations', async () => {
	let caught: unknown
	runLifecycleOperation(
		async () => { throw new Error('failed lifecycle') },
		(error) => { caught = error },
	)
	await Bun.sleep(0)

	expect(caught).toBeInstanceOf(Error)
})

test('shutdown prevents an Add still waiting for its discovery query', async () => {
	let resolveDiscovery: ((record: StreamElementsIntegration) => void) | undefined
	const discovery = new Promise<StreamElementsIntegration>((resolve) => { resolveDiscovery = resolve })
	let created = 0
	const store = createStreamElementsStore({
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
	expect(store.connections.size).toBe(0)
})
