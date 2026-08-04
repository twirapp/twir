import { expect, test } from 'bun:test'

import {
	OAuthRefreshLockLostError,
	type RedisCommands,
	withOAuthRefreshLock,
} from './oauth-lock.ts'

class FakeRedis implements RedisCommands {
	readonly calls: Array<{ command: string; args: readonly string[] }> = []
	renewResult: unknown = 1

	async send(command: string, args: readonly string[]): Promise<unknown> {
		this.calls.push({ command, args })
		if (command === 'SET') {
			return 'OK'
		}
		if (args.length === 5) {
			return this.renewResult
		}
		return 1
	}
}

test('uses the shared OAuth lock key and releases only for its owner', async () => {
	const redis = new FakeRedis()

	const result = await withOAuthRefreshLock(
		'streamelements',
		'channel-1',
		async () => 'done',
		{ redis, createOwner: () => 'opaque-owner' }
	)

	expect(result).toBe('done')
	expect(redis.calls[0]).toEqual({
		command: 'SET',
		args: [
			'twir:integration-token-refresh:streamelements:channel-1',
			'opaque-owner',
			'NX',
			'PX',
			'30000',
		],
	})
	expect(redis.calls.at(-1)?.command).toBe('EVAL')
	expect(redis.calls.at(-1)?.args.slice(1)).toEqual([
		'1',
		'twir:integration-token-refresh:streamelements:channel-1',
		'opaque-owner',
	])
})

test('bounds acquisition retries', async () => {
	const calls: Array<{ command: string; args: readonly string[] }> = []
	const redis: RedisCommands = {
		async send(command, args) {
			calls.push({ command, args })
			return null
		},
	}

	await expect(withOAuthRefreshLock('streamlabs', 'channel-2', async () => undefined, {
		redis,
		acquireAttempts: 3,
		retryDelayMs: 0,
	})).rejects.toThrow('OAuth refresh lock is unavailable')
	expect(calls).toHaveLength(3)
})

test('times out a stuck acquisition without invoking the callback', async () => {
	let callbackCalled = false
	const redis: RedisCommands = {
		async send() {
			return new Promise<never>(() => undefined)
		},
	}

	await expect(withOAuthRefreshLock('streamelements', 'channel-stuck', async () => {
		callbackCalled = true
	}, {
		redis,
		acquireAttempts: 1,
		acquireOperationTimeoutMs: 5,
	})).rejects.toThrow('OAuth refresh lock is unavailable')
	expect(callbackCalled).toBe(false)
})

test('aborts and fails the callback when renewal loses ownership', async () => {
	const redis = new FakeRedis()
	redis.renewResult = 0
	let aborted = false

	const operation = withOAuthRefreshLock(
		'streamlabs',
		'channel-3',
		async (signal) => {
			await new Promise<void>((resolve) => {
				signal.addEventListener('abort', () => {
					aborted = true
					resolve()
				}, { once: true })
			})
		},
		{
			redis,
			renewalIntervalMs: 1,
			renewalOperationTimeoutMs: 20,
			leaseWatchdogMs: 50,
		}
	)

	await expect(operation).rejects.toBeInstanceOf(OAuthRefreshLockLostError)
	expect(aborted).toBe(true)
})
