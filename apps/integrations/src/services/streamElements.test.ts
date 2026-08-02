import { expect, test } from 'bun:test'

import type { Donate } from '../utils/onDonation.ts'
import {
	StreamElementsConnection,
	type StreamElementsSocket,
	type StreamElementsSocketFactory,
	type StreamElementsSocketOptions,
} from './streamElements.ts'

class FakeSocket implements StreamElementsSocket {
	readonly handlers = new Map<string, Array<(payload?: unknown) => void>>()
	readonly emitted: Array<{ event: string; payload: unknown }> = []
	closed = false

	on(event: string, handler: (payload?: unknown) => void): StreamElementsSocket {
		const handlers = this.handlers.get(event) ?? []
		handlers.push(handler)
		this.handlers.set(event, handlers)
		return this
	}

	emit(event: string, payload: unknown): StreamElementsSocket {
		this.emitted.push({ event, payload })
		return this
	}

	close(): void {
		this.closed = true
	}

	trigger(event: string, payload?: unknown): void {
		for (const handler of this.handlers.get(event) ?? []) handler(payload)
	}
}

function setup() {
	const sockets: FakeSocket[] = []
	const factoryCalls: Array<{ url: string; options: StreamElementsSocketOptions }> = []
	const socketFactory: StreamElementsSocketFactory = (url, options) => {
		factoryCalls.push({ url, options })
		const socket = new FakeSocket()
		sockets.push(socket)
		return socket
	}
	const donations: Donate[] = []
	const claimed: string[] = []
	let accessToken = 'access-1'
	let refreshes = 0
	const connection = new StreamElementsConnection({
		channelID: 'channel-1',
		client: {
			get tokens() {
				return { accessToken, refreshToken: 'refresh' }
			},
			async refresh() {
				refreshes += 1
				accessToken = 'access-2'
			},
		},
		socketFactory,
		onDonation: async (donation) => { donations.push(donation) },
		claimDonation: async (_provider, eventID) => {
			claimed.push(eventID)
			return eventID !== 'duplicate'
		},
		schedule: (callback) => {
			callback()
			return 1
		},
		clearSchedule: () => undefined,
		random: () => 0.5,
	})
	connection.connect()
	return { connection, sockets, factoryCalls, donations, claimed, refreshes: () => refreshes }
}

test('connects with websocket-only transport and authenticates after connect', () => {
	const fake = setup()

	expect(fake.factoryCalls).toEqual([{
		url: 'https://realtime.streamelements.com',
		options: { transports: ['websocket'], reconnection: false, forceNew: true },
	}])
	fake.sockets[0]?.trigger('connect')
	expect(fake.sockets[0]?.emitted).toEqual([{
		event: 'authenticate',
		payload: { method: 'oauth2', token: 'access-1' },
	}])
})

test('normalizes tips, deduplicates stable IDs, and passes through missing IDs', async () => {
	const fake = setup()
	const socket = fake.sockets[0]!
	socket.trigger('event', { type: 'subscriber', data: {} })
	socket.trigger('event', {
		type: 'tip',
		data: { tipId: 'tip-1', username: 'Alice', amount: 12.5, currency: 'USD', message: 'hi' },
	})
	socket.trigger('event', {
		type: 'tip',
		data: { tipId: 'duplicate', username: 'Bob', amount: 2, currency: 'EUR', message: '' },
	})
	socket.trigger('event', {
		type: 'tip',
		data: { username: '', amount: '3.00', currency: 'GBP' },
	})
	await Bun.sleep(0)

	expect(fake.claimed).toEqual(['tip-1', 'duplicate', ''])
	expect(fake.donations).toEqual([
		{
			twitchUserId: 'channel-1', amount: 12.5, currency: 'USD', message: 'hi', userName: 'Alice',
		},
		{
			twitchUserId: 'channel-1', amount: '3.00', currency: 'GBP', message: null, userName: null,
		},
	])
})

test('unauthorized refreshes once, rebuilds with new token, and ignores stale callbacks', async () => {
	const fake = setup()
	const first = fake.sockets[0]!
	first.trigger('unauthorized')
	await Bun.sleep(0)

	expect(fake.refreshes()).toBe(1)
	expect(first.closed).toBe(true)
	expect(fake.sockets).toHaveLength(2)
	fake.sockets[1]?.trigger('connect')
	expect(fake.sockets[1]?.emitted[0]).toEqual({
		event: 'authenticate', payload: { method: 'oauth2', token: 'access-2' },
	})

	first.trigger('disconnect', 'network')
	expect(fake.sockets).toHaveLength(2)
	fake.sockets[1]?.trigger('unauthorized')
	await Bun.sleep(0)
	expect(fake.refreshes()).toBe(1)
})

test('confirmed authentication starts a new future refresh cycle', async () => {
	const fake = setup()
	fake.sockets[0]?.trigger('unauthorized')
	await Bun.sleep(0)
	fake.sockets[1]?.trigger('connect')
	fake.sockets[1]?.trigger('authenticated')
	fake.sockets[1]?.trigger('unauthorized')
	await Bun.sleep(0)

	expect(fake.refreshes()).toBe(2)
	expect(fake.sockets).toHaveLength(3)
})

test('does not publish a claimed tip after the connection is destroyed', async () => {
	const sockets: FakeSocket[] = []
	const donations: Donate[] = []
	let resolveClaim: ((claimed: boolean) => void) | undefined
	const claim = new Promise<boolean>((resolve) => { resolveClaim = resolve })
	const connection = new StreamElementsConnection({
		channelID: 'channel-1',
		client: {
			tokens: { accessToken: 'access', refreshToken: 'refresh' },
			async refresh() { return undefined },
		},
		socketFactory: () => {
			const socket = new FakeSocket()
			sockets.push(socket)
			return socket
		},
		onDonation: async (donation) => { donations.push(donation) },
		claimDonation: async () => claim,
	})
	connection.connect()
	sockets[0]?.trigger('event', {
		type: 'tip',
		data: { tipId: 'tip-1', username: 'Alice', amount: 1, currency: 'USD' },
	})
	await Bun.sleep(0)
	connection.destroy()
	resolveClaim?.(true)
	await Bun.sleep(0)

	expect(donations).toHaveLength(0)
})

test('network disconnect reconnects without refresh and destroy cancels stale work', () => {
	const callbacks: Array<() => void> = []
	const sockets: FakeSocket[] = []
	const connection = new StreamElementsConnection({
		channelID: 'channel-1',
		client: {
			tokens: { accessToken: 'access', refreshToken: 'refresh' },
			async refresh() { throw new Error('must not refresh') },
		},
		socketFactory: () => {
			const socket = new FakeSocket()
			sockets.push(socket)
			return socket
		},
		onDonation: async () => undefined,
		claimDonation: async () => true,
		schedule: (callback) => { callbacks.push(callback); return callbacks.length },
		clearSchedule: () => undefined,
		random: () => 0.5,
	})
	connection.connect()
	sockets[0]?.trigger('disconnect', 'network')
	expect(callbacks).toHaveLength(1)
	connection.destroy()
	callbacks[0]?.()
	expect(sockets).toHaveLength(1)
})

test('network reconnect delay is jittered and capped at thirty seconds', () => {
	const callbacks: Array<() => void> = []
	const delays: number[] = []
	const sockets: FakeSocket[] = []
	let refreshes = 0
	const connection = new StreamElementsConnection({
		channelID: 'channel-1',
		client: {
			tokens: { accessToken: 'access', refreshToken: 'refresh' },
			async refresh() { refreshes += 1 },
		},
		socketFactory: () => {
			const socket = new FakeSocket()
			sockets.push(socket)
			return socket
		},
		onDonation: async () => undefined,
		claimDonation: async () => true,
		schedule: (callback, milliseconds) => {
			callbacks.push(callback)
			delays.push(milliseconds)
			return callbacks.length
		},
		clearSchedule: () => undefined,
		random: () => 0.99,
	})
	connection.connect()

	for (let attempt = 0; attempt < 8; attempt += 1) {
		sockets.at(-1)?.trigger('disconnect', 'network')
		callbacks.shift()?.()
	}

	expect(Math.max(...delays)).toBeLessThanOrEqual(30_000)
	expect(refreshes).toBe(0)
	connection.destroy()
})
