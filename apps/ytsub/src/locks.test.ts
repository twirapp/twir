import { expect, test } from 'bun:test'

import { RedisBindingOwnership, type RedisClient } from './locks.ts'

class FakeRedis implements RedisClient {
	readonly values = new Map<string, string>()

	async set(key: string, value: string, ...options: readonly string[]): Promise<string | null> {
		if (options.includes('NX') && this.values.has(key)) {
			return null
		}
		if (options.includes('XX') && !this.values.has(key)) {
			return null
		}
		this.values.set(key, value)
		return 'OK'
	}

	async get(key: string): Promise<string | null> {
		return this.values.get(key) ?? null
	}

	async del(...keys: readonly string[]): Promise<number> {
		let deleted = 0
		for (const key of keys) {
			if (this.values.delete(key)) {
				deleted += 1
			}
		}
		return deleted
	}

	async send(command: string, args: readonly string[]): Promise<unknown> {
		if (command !== 'EVAL') {
			throw new Error(`Unexpected Redis command: ${command}`)
		}
		const key = args[2]
		const expectedReplicaId = args[3]
		if (key === undefined || expectedReplicaId === undefined) {
			throw new Error('Missing EVAL key arguments')
		}
		if (this.values.get(key) !== expectedReplicaId) {
			return 0
		}
		if (args.length === 5) {
			return 1
		}
		this.values.delete(key)
		return 1
	}

	close(): void {}
}

test('RedisBindingOwnership grants a binding to only one replica', async () => {
	const redis = new FakeRedis()
	const first = new RedisBindingOwnership(redis, { replicaId: 'replica-1' })
	const second = new RedisBindingOwnership(redis, { replicaId: 'replica-2' })

	expect(await first.tryAcquire('binding-1')).toBe(true)
	expect(await second.tryAcquire('binding-1')).toBe(false)

	await first.close()
	await second.close()
})

test('RedisBindingOwnership releases only locks held by its own replica', async () => {
	const redis = new FakeRedis()
	const owner = new RedisBindingOwnership(redis, { replicaId: 'replica-1' })

	await owner.tryAcquire('binding-1')
	redis.values.set('ytsub:lock:binding-1', 'replica-2')
	await owner.release('binding-1')

	expect(redis.values.get('ytsub:lock:binding-1')).toBe('replica-2')
	await owner.close()
})
