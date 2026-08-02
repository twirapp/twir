import { expect, test } from 'bun:test'

import type { ProviderTokenStore, ProviderTokens } from './provider-token-store.ts'
import { StreamElementsClient, StreamElementsRefreshError } from './streamelements-client.ts'

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
	_provider: 'streamelements',
	_channelID: string,
	callback: (signal: AbortSignal) => Promise<T>
) => callback(new AbortController().signal)

test('refresh rereads under lock, persists rotated tokens, and uses form encoding', async () => {
	const fake = createStore({ accessToken: 'old-access', refreshToken: 'old-refresh' })
	let request: Request | undefined
	const client = new StreamElementsClient({
		channelID: 'channel-1',
		tokens: { accessToken: 'old-access', refreshToken: 'old-refresh' },
		tokenStore: fake.store,
		clientID: 'client',
		clientSecret: 'secret',
		lock: directLock,
		fetch: async (input, init) => {
			request = new Request(input, init)
			return Response.json({ access_token: 'new-access', refresh_token: 'new-refresh' })
		},
	})

	await client.refresh()

	expect(request?.url).toBe('https://api.streamelements.com/oauth2/token')
	expect(request?.headers.get('content-type')).toBe('application/x-www-form-urlencoded')
	expect(await request?.text()).toBe(
		'grant_type=refresh_token&client_id=client&client_secret=secret&refresh_token=old-refresh'
	)
	expect(fake.updates).toEqual([{ accessToken: 'new-access', refreshToken: 'new-refresh' }])
	expect(client.tokens).toEqual({ accessToken: 'new-access', refreshToken: 'new-refresh' })
})

test('refresh uses the old refresh token when rotation omits it', async () => {
	const fake = createStore({ accessToken: 'old-access', refreshToken: 'old-refresh' })
	const client = new StreamElementsClient({
		channelID: 'channel-1',
		tokens: { accessToken: 'old-access', refreshToken: 'old-refresh' },
		tokenStore: fake.store,
		clientID: 'client',
		clientSecret: 'secret',
		lock: directLock,
		fetch: async () => Response.json({ access_token: 'new-access' }),
	})

	await client.refresh()

	expect(fake.updates).toEqual([{ accessToken: 'new-access', refreshToken: 'old-refresh' }])
})

test('refresh skips provider request when another process already rotated credentials', async () => {
	const fake = createStore({ accessToken: 'old-access', refreshToken: 'old-refresh' })
	fake.set({ accessToken: 'other-access', refreshToken: 'other-refresh' })
	let fetches = 0
	const client = new StreamElementsClient({
		channelID: 'channel-1',
		tokens: { accessToken: 'old-access', refreshToken: 'old-refresh' },
		tokenStore: fake.store,
		clientID: 'client',
		clientSecret: 'secret',
		lock: directLock,
		fetch: async () => {
			fetches += 1
			return Response.json({})
		},
	})

	await client.refresh()

	expect(fetches).toBe(0)
	expect(client.tokens).toEqual({ accessToken: 'other-access', refreshToken: 'other-refresh' })
})

test('refresh keeps memory unchanged and sanitizes provider failures', async () => {
	const fake = createStore({ accessToken: 'old-access', refreshToken: 'old-refresh' })
	const client = new StreamElementsClient({
		channelID: 'channel-1',
		tokens: { accessToken: 'old-access', refreshToken: 'old-refresh' },
		tokenStore: fake.store,
		clientID: 'client',
		clientSecret: 'secret',
		lock: directLock,
		fetch: async () => new Response('secret provider body', { status: 401 }),
	})

	await expect(client.refresh()).rejects.toBeInstanceOf(StreamElementsRefreshError)
	await expect(client.refresh()).rejects.not.toThrow('secret provider body')
	expect(fake.updates).toHaveLength(0)
	expect(client.tokens).toEqual({ accessToken: 'old-access', refreshToken: 'old-refresh' })
})

test('refresh changes memory only after persistence succeeds', async () => {
	const store: ProviderTokenStore = {
		async getTokens() {
			return { accessToken: 'old-access', refreshToken: 'old-refresh' }
		},
		async updateTokens() {
			throw new Error('database unavailable')
		},
	}
	const client = new StreamElementsClient({
		channelID: 'channel-1',
		tokens: { accessToken: 'old-access', refreshToken: 'old-refresh' },
		tokenStore: store,
		clientID: 'client',
		clientSecret: 'secret',
		lock: directLock,
		fetch: async () => Response.json({ access_token: 'new-access', refresh_token: 'new-refresh' }),
	})

	await expect(client.refresh()).rejects.toBeInstanceOf(StreamElementsRefreshError)
	expect(client.tokens).toEqual({ accessToken: 'old-access', refreshToken: 'old-refresh' })
})
