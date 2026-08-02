import { expect, test } from 'bun:test'

import type { Donate } from '../utils/onDonation.ts'
import {
	StreamLabsConnection,
	type StreamLabsSocket,
	type StreamLabsSocketFactory,
	type StreamLabsSocketOptions,
} from './streamLabs.ts'

class FakeSocket implements StreamLabsSocket {
	readonly handlers = new Map<string, Array<(payload?: unknown) => void>>()
	closed = false

	on(event: string, handler: (payload?: unknown) => void): StreamLabsSocket {
		const handlers = this.handlers.get(event) ?? []
		handlers.push(handler)
		this.handlers.set(event, handlers)
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
	const factoryCalls: Array<{ url: string; options: StreamLabsSocketOptions }> = []
	const socketFactory: StreamLabsSocketFactory = (url, options) => {
		factoryCalls.push({ url, options })
		const socket = new FakeSocket()
		sockets.push(socket)
		return socket
	}
	const donations: Donate[] = []
	const claims: Array<{ provider: string; eventID: string }> = []
	const errors: unknown[] = []
	const connection = new StreamLabsConnection({
		channelID: 'channel-1',
		socketToken: 'socket token',
		socketFactory,
		onDonation: async (donation) => { donations.push(donation) },
		claimDonation: async (provider, eventID) => {
			claims.push({ provider, eventID })
			return eventID !== 'duplicate'
		},
		schedule: (callback) => { callback(); return 1 },
		clearSchedule: () => undefined,
		random: () => 0.5,
		onError: (error) => { errors.push(error) },
	})
	connection.connect()
	return { connection, sockets, factoryCalls, donations, claims, errors }
}

test('connects to Streamlabs with websocket-only isolated socket options', () => {
	const fake = setup()
	expect(fake.factoryCalls).toEqual([{
		url: 'https://sockets.streamlabs.com?token=socket+token',
		options: { transports: ['websocket'], reconnection: false, forceNew: true },
	}])
})

test('normalizes every donation message and deduplicates stable message IDs', async () => {
	const fake = setup()
	const socket = fake.sockets[0]!
	socket.trigger('event', { type: 'follow', message: [] })
	socket.trigger('event', {
		type: 'donation',
		event_id: 'event-1',
		message: [
			{ _id: 'message-1', amount: 12.5, currency: 'USD', message: 'hello', from: 'Alice' },
			{ _id: 'duplicate', amount: 2, currency: 'EUR', message: '', from: 'Bob' },
			{ amount: '3.00', currency: 'GBP', from: '' },
		],
	})
	await Bun.sleep(0)

	expect(fake.claims).toEqual([
		{ provider: 'streamlabs', eventID: 'message-1' },
		{ provider: 'streamlabs', eventID: 'duplicate' },
		{ provider: 'streamlabs', eventID: 'event-1:2' },
	])
	expect(fake.donations).toEqual([
		{
			twitchUserId: 'channel-1', amount: 12.5, currency: 'USD', message: 'hello', userName: 'Alice',
		},
		{
			twitchUserId: 'channel-1', amount: '3.00', currency: 'GBP', message: null, userName: null,
		},
	])
})

test('handles malformed and failing donation messages independently', async () => {
	const fake = setup()
	let calls = 0
	const donations: Donate[] = []
	const connection = new StreamLabsConnection({
		channelID: 'channel-1',
		socketToken: 'token',
		socketFactory: () => fake.sockets[0]!,
		claimDonation: async () => true,
		onDonation: async (donation) => {
			calls += 1
			if (calls === 1) throw new Error('first failed')
			donations.push(donation)
		},
		onError: (error) => { fake.errors.push(error) },
	})
	connection.connect()
	fake.sockets[0]?.trigger('event', {
		type: 'donation',
		message: [
			{ _id: 'bad', amount: null, currency: 'USD' },
			{ _id: 'first', amount: 1, currency: 'USD', from: 'A' },
			{ _id: 'second', amount: 2, currency: 'USD', from: 'B' },
		],
	})
	await Bun.sleep(0)

	expect(calls).toBe(2)
	expect(donations).toHaveLength(1)
	expect(donations[0]?.userName).toBe('B')
	expect(fake.errors).toHaveLength(1)
	connection.destroy()
})

test('network failures use capped reconnect and destroy cancels stale callbacks', () => {
	const callbacks: Array<() => void> = []
	const delays: number[] = []
	const sockets: FakeSocket[] = []
	const connection = new StreamLabsConnection({
		channelID: 'channel-1',
		socketToken: 'token',
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
		sockets.at(-1)?.trigger('connect_error', new Error('network'))
		callbacks.shift()?.()
	}
	expect(Math.max(...delays)).toBeLessThanOrEqual(30_000)

	sockets.at(-1)?.trigger('disconnect', 'network')
	const stale = callbacks.shift()!
	const count = sockets.length
	connection.destroy()
	stale()
	expect(sockets).toHaveLength(count)
})

test('does not publish a claimed donation after removal', async () => {
	const socket = new FakeSocket()
	const donations: Donate[] = []
	let resolveClaim: ((value: boolean) => void) | undefined
	const claim = new Promise<boolean>((resolve) => { resolveClaim = resolve })
	const connection = new StreamLabsConnection({
		channelID: 'channel-1',
		socketToken: 'token',
		socketFactory: () => socket,
		onDonation: async (donation) => { donations.push(donation) },
		claimDonation: async () => claim,
	})
	connection.connect()
	socket.trigger('event', {
		type: 'donation', message: [{ _id: 'one', amount: 1, currency: 'USD', from: 'A' }],
	})
	await Bun.sleep(0)
	connection.destroy()
	resolveClaim?.(true)
	await Bun.sleep(0)
	expect(donations).toHaveLength(0)
})
