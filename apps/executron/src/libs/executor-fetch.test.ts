import { describe, expect, mock, test } from 'bun:test'

import type { FetchResponse } from './host-fetch'

const hostFetchMock = mock(
	async (_url: string, _options?: RequestInit): Promise<FetchResponse> => ({
		status: 200,
		statusText: 'OK',
		headers: { 'content-type': 'application/json' },
		body: '{"ok":true}',
	})
)

mock.module('./host-fetch', () => ({ hostFetch: hostFetchMock }))

const { executeCode } = await import('./executor')

describe('executeCode fetch options passthrough', () => {
	test('passes method, headers and body to hostFetch', async () => {
		hostFetchMock.mockClear()

		const execution = await executeCode(
			`
				const res = await fetch('https://example.com/api', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						'X-Api-Key': 'secret-key',
					},
					body: JSON.stringify({ hello: 'world' }),
				});
				const data = await res.json();
				return JSON.stringify({ status: res.status, ok: res.ok, data });
			`,
			'test-channel',
			new Map()
		)

		expect(execution.error).toBe('')
		expect(JSON.parse(execution.result)).toEqual({
			status: 200,
			ok: true,
			data: { ok: true },
		})

		expect(hostFetchMock).toHaveBeenCalledTimes(1)
		const [url, opts] = hostFetchMock.mock.calls[0] as unknown as [string, RequestInit]
		expect(url).toBe('https://example.com/api')
		expect(opts.method).toBe('POST')
		expect(opts.headers).toEqual({
			'Content-Type': 'application/json',
			'X-Api-Key': 'secret-key',
		})
		expect(opts.body).toBe('{"hello":"world"}')
		expect(opts.signal).toBeInstanceOf(AbortSignal)
	})

	test('passes headers given as array of tuples', async () => {
		hostFetchMock.mockClear()

		const execution = await executeCode(
			`
				await fetch('https://example.com/api', {
					headers: [['X-First', '1'], ['X-Second', '2']],
				});
				return 'done';
			`,
			'test-channel',
			new Map()
		)

		expect(execution.error).toBe('')
		const [, opts] = hostFetchMock.mock.calls[0] as unknown as [string, RequestInit]
		expect(opts.headers).toEqual({ 'X-First': '1', 'X-Second': '2' })
	})

	test('passes misc request options through', async () => {
		hostFetchMock.mockClear()

		const execution = await executeCode(
			`
				await fetch('https://example.com/api', {
					credentials: 'include',
					mode: 'cors',
					referrer: 'https://example.com/page',
				});
				return 'done';
			`,
			'test-channel',
			new Map()
		)

		expect(execution.error).toBe('')
		const [, opts] = hostFetchMock.mock.calls[0] as unknown as [string, RequestInit]
		expect(opts.credentials).toBe('include')
		expect(opts.mode).toBe('cors')
		expect(opts.referrer).toBe('https://example.com/page')
	})

	test('works without options at all', async () => {
		hostFetchMock.mockClear()

		const execution = await executeCode(
			`
				const res = await fetch('https://example.com/api');
				return res.status;
			`,
			'test-channel',
			new Map()
		)

		expect(execution.error).toBe('')
		expect(execution.result).toBe('200')

		const [, opts] = hostFetchMock.mock.calls[0] as unknown as [string, RequestInit]
		expect(opts.method).toBeUndefined()
		expect(opts.headers).toBeUndefined()
		expect(opts.body).toBeUndefined()
	})
})
