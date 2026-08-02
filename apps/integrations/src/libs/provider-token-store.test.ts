import { expect, test } from 'bun:test'

import {
	ProviderTokensNotFoundError,
	type TokenQuery,
	createProviderTokenStores,
} from './provider-token-store.ts'

interface CapturedQuery {
	readonly text: string
	readonly values: readonly unknown[]
}

function createQuery(results: Array<readonly Record<string, unknown>[]>) {
	const calls: CapturedQuery[] = []
	const query: TokenQuery = async <Row extends Record<string, unknown>>(
		strings: TemplateStringsArray,
		...values: readonly unknown[]
	): Promise<Row[]> => {
		calls.push({ text: strings.join('$'), values })
		return (results.shift() ?? []) as Row[]
	}
	return { calls, query }
}

test('StreamElements reads enabled tokens through the service join', async () => {
	const fake = createQuery([[{ accessToken: 'access', refreshToken: 'refresh' }]])
	const { streamElements } = createProviderTokenStores(fake.query)

	expect(await streamElements.getTokens('channel-1')).toEqual({
		accessToken: 'access',
		refreshToken: 'refresh',
	})
	expect(fake.calls[0]?.text).toContain('FROM channels_integrations AS ci')
	expect(fake.calls[0]?.text).toContain('JOIN integrations AS i')
	expect(fake.calls[0]?.text).toContain("i.service = 'STREAMELEMENTS'")
	expect(fake.calls[0]?.text).toContain('ci.enabled = TRUE')
	expect(fake.calls[0]?.values).toEqual(['channel-1'])
})

test('StreamElements updates both tokens on one enabled matching row', async () => {
	const fake = createQuery([[{ id: 'integration-1' }]])
	const { streamElements } = createProviderTokenStores(fake.query)

	await streamElements.updateTokens('channel-1', {
		accessToken: 'next-access',
		refreshToken: 'next-refresh',
	})

	expect(fake.calls[0]?.text).toContain('UPDATE channels_integrations AS ci')
	expect(fake.calls[0]?.text).toContain('"accessToken" = $')
	expect(fake.calls[0]?.text).toContain('"refreshToken" = $')
	expect(fake.calls[0]?.text).toContain("i.service = 'STREAMELEMENTS'")
	expect(fake.calls[0]?.text).toContain('ci.enabled = TRUE')
	expect(fake.calls[0]?.values).toEqual(['next-access', 'next-refresh', 'channel-1'])
})

test('Streamlabs reads and atomically updates its enabled dedicated record', async () => {
	const fake = createQuery([
		[{ accessToken: 'access', refreshToken: 'refresh' }],
		[{ id: 'integration-2' }],
	])
	const { streamLabs } = createProviderTokenStores(fake.query)

	expect(await streamLabs.getTokens('channel-2')).toEqual({
		accessToken: 'access',
		refreshToken: 'refresh',
	})
	await streamLabs.updateTokens('channel-2', {
		accessToken: 'next-access',
		refreshToken: 'next-refresh',
	})

	expect(fake.calls[0]?.text).toContain('FROM channels_integrations_streamlabs')
	expect(fake.calls[0]?.text).toContain('enabled = TRUE')
	expect(fake.calls[1]?.text).toContain('UPDATE channels_integrations_streamlabs')
	expect(fake.calls[1]?.text).toContain('access_token = $')
	expect(fake.calls[1]?.text).toContain('refresh_token = $')
	expect(fake.calls[1]?.text).toContain('enabled = TRUE')
	expect(fake.calls[1]?.values).toEqual(['next-access', 'next-refresh', 'channel-2'])
})

test('rejects token persistence when no enabled integration matches', async () => {
	const fake = createQuery([[]])
	const { streamLabs } = createProviderTokenStores(fake.query)

	await expect(streamLabs.updateTokens('channel-3', {
		accessToken: 'access',
		refreshToken: 'refresh',
	})).rejects.toBeInstanceOf(ProviderTokensNotFoundError)
})
