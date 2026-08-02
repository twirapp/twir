import { expect, test } from 'bun:test'

import type { ProviderTokenStore, ProviderTokens } from './provider-token-store.ts'
import {
	StreamLabsClient,
	StreamLabsClientError,
	type StreamLabsFetch,
} from './streamlabs-client.ts'

function createStore(tokens: ProviderTokens) {
	let persisted = tokens
	const updates: ProviderTokens[] = []
	const store: ProviderTokenStore = {
		async getTokens() {
			return persisted
		},
		async updateTokens(_channelID, next) {
			updates.push(next)
			persisted = next
		},
	}
	return { store, updates, set: (next: ProviderTokens) => { persisted = next } }
}

const directLock = async <T>(
	_provider: 'streamlabs',
	_channelID: string,
	callback: (signal: AbortSignal) => Promise<T>
) => callback(new AbortController().signal)

function options(
	fetch: StreamLabsFetch,
	store = createStore({ accessToken: 'old-access', refreshToken: 'old-refresh' })
) {
	return {
		client: new StreamLabsClient({
			channelID: 'channel-1',
			tokens: { accessToken: 'old-access', refreshToken: 'old-refresh' },
			tokenStore: store.store,
			clientID: 'client',
			clientSecret: 'secret',
			redirectURI: 'https://twir.test/dashboard/integrations/streamlabs',
			lock: directLock,
			fetch,
		}),
		store,
	}
}

test('gets a socket token without refreshing successful credentials', async () => {
	const requests: Request[] = []
	const fake = options(async (input, init) => {
		requests.push(new Request(input, init))
		return Response.json({ socket_token: 'socket-1' })
	})

	await expect(fake.client.getSocketToken()).resolves.toEqual({ socketToken: 'socket-1' })
	expect(requests).toHaveLength(1)
	expect(requests[0]?.url).toBe('https://streamlabs.com/api/v2.0/socket/token')
	expect(requests[0]?.headers.get('authorization')).toBe('Bearer old-access')
	expect(fake.store.updates).toHaveLength(0)
})

test('first 401 rereads, refreshes with a form, persists, and retries once', async () => {
	const requests: Request[] = []
	const fake = options(async (input, init) => {
		const request = new Request(input, init)
		requests.push(request)
		if (request.url.endsWith('/socket/token') && requests.length === 1) {
			return new Response('expired access token', { status: 401 })
		}
		if (request.url === 'https://streamlabs.com/api/v2.0/token') {
			return Response.json({ access_token: 'new-access', refresh_token: 'new-refresh' })
		}
		return Response.json({ socket_token: 'socket-2' })
	})

	await expect(fake.client.getSocketToken()).resolves.toEqual({ socketToken: 'socket-2' })
	expect(requests).toHaveLength(3)
	const refresh = requests[1]!
	expect(refresh.headers.get('content-type')).toBe('application/x-www-form-urlencoded')
	expect(await refresh.text()).toBe(
		'grant_type=refresh_token&client_id=client&client_secret=secret&refresh_token=old-refresh&redirect_uri=https%3A%2F%2Ftwir.test%2Fdashboard%2Fintegrations%2Fstreamlabs'
	)
	expect(requests[2]?.headers.get('authorization')).toBe('Bearer new-access')
	expect(fake.store.updates).toEqual([{ accessToken: 'new-access', refreshToken: 'new-refresh' }])
	expect(fake.client.tokens).toEqual({ accessToken: 'new-access', refreshToken: 'new-refresh' })
})

test('uses tokens rotated by another process and skips provider refresh', async () => {
	const store = createStore({ accessToken: 'old-access', refreshToken: 'old-refresh' })
	store.set({ accessToken: 'other-access', refreshToken: 'other-refresh' })
	const requests: Request[] = []
	const fake = options(async (input, init) => {
		const request = new Request(input, init)
		requests.push(request)
		if (requests.length === 1) return new Response(null, { status: 401 })
		return Response.json({ socket_token: 'socket-2' })
	}, store)

	await expect(fake.client.getSocketToken()).resolves.toEqual({ socketToken: 'socket-2' })
	expect(requests).toHaveLength(2)
	expect(requests[1]?.headers.get('authorization')).toBe('Bearer other-access')
	expect(store.updates).toHaveLength(0)
})

test('preserves an omitted refresh token and changes memory only after persistence', async () => {
	const fake = options(async (input) => {
		const url = new Request(input).url
		if (url.endsWith('/socket/token') && fake.store.updates.length === 0) {
			return new Response(null, { status: 401 })
		}
		if (url === 'https://streamlabs.com/api/v2.0/token') {
			return Response.json({ access_token: 'new-access' })
		}
		return Response.json({ socket_token: 'socket-2' })
	})

	await expect(fake.client.getSocketToken()).resolves.toEqual({ socketToken: 'socket-2' })
	expect(fake.store.updates).toEqual([{ accessToken: 'new-access', refreshToken: 'old-refresh' }])

	const failingStore: ProviderTokenStore = {
		async getTokens() {
			return { accessToken: 'old-access', refreshToken: 'old-refresh' }
		},
		async updateTokens() {
			throw new Error('database unavailable')
		},
	}
	let socketRequests = 0
	const failing = new StreamLabsClient({
		channelID: 'channel-1',
		tokens: { accessToken: 'old-access', refreshToken: 'old-refresh' },
		tokenStore: failingStore,
		clientID: 'client',
		clientSecret: 'secret',
		redirectURI: 'https://twir.test/callback',
		lock: directLock,
		fetch: async (input) => {
			if (new Request(input).url.endsWith('/socket/token')) {
				socketRequests += 1
				return new Response(null, { status: 401 })
			}
			return Response.json({ access_token: 'new-access', refresh_token: 'new-refresh' })
		},
	})
	await expect(failing.getSocketToken()).rejects.toBeInstanceOf(StreamLabsClientError)
	expect(failing.tokens).toEqual({ accessToken: 'old-access', refreshToken: 'old-refresh' })
	expect(socketRequests).toBe(1)
})

test('second 401 stops and all provider failures are sanitized', async () => {
	const bodies = ['first secret body', 'refresh secret body', 'second secret body']
	for (const body of bodies) {
		let socketRequests = 0
		const fake = options(async (input) => {
			const request = new Request(input)
			if (request.url === 'https://streamlabs.com/api/v2.0/token') {
				return body === 'refresh secret body'
					? new Response(body, { status: 500 })
					: Response.json({ access_token: 'new-access' })
			}
			socketRequests += 1
			if (body === 'first secret body') return new Response(body, { status: 500 })
			return new Response(body, { status: 401 })
		})

		let error: unknown
		try {
			await fake.client.getSocketToken()
		} catch (caught) {
			error = caught
		}
		expect(error).toBeInstanceOf(StreamLabsClientError)
		expect(String(error)).not.toContain(body)
		expect(String(error)).not.toContain('secret')
		if (body === 'second secret body') expect(socketRequests).toBe(2)
	}
})

test('bounds provider responses and attaches a request deadline', async () => {
	let signal: AbortSignal | null | undefined
	const fake = options(async (_input, init) => {
		signal = init?.signal
		return new Response('{}', {
			headers: { 'content-length': String(1024 * 1024 + 1) },
		})
	})

	await expect(fake.client.getSocketToken()).rejects.toBeInstanceOf(StreamLabsClientError)
	expect(signal).toBeInstanceOf(AbortSignal)
})
