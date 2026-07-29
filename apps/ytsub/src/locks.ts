const LOCK_PREFIX = 'ytsub:lock:'
const LOCK_TTL_MS = 30_000
const RENEWAL_INTERVAL_MS = 10_000
const RELEASE_IF_OWNER_SCRIPT = "if redis.call('GET', KEYS[1]) == ARGV[1] then return redis.call('DEL', KEYS[1]) else return 0 end"
const RENEW_IF_OWNER_SCRIPT = "if redis.call('GET', KEYS[1]) == ARGV[1] then redis.call('PSETEX', KEYS[1], ARGV[2], ARGV[1]) return 1 else return 0 end"

export interface RedisClient {
	set(key: string, value: string, ...options: readonly string[]): Promise<string | null>
	send(command: string, args: readonly string[]): Promise<unknown>
	close(): void
}

export interface BindingOwnershipManager {
	tryAcquire(bindingId: string): Promise<boolean>
	release(bindingId: string): Promise<void>
	onLostOwnership(listener: (bindingId: string) => void): () => void
	close(): Promise<void>
}

export interface RedisBindingOwnershipOptions {
	readonly replicaId?: string
	readonly ttlMs?: number
	readonly renewalIntervalMs?: number
}

export class RedisBindingOwnership implements BindingOwnershipManager {
	readonly #ownedBindingIds = new Set<string>()
	readonly #lostOwnershipListeners = new Set<(bindingId: string) => void>()
	readonly #replicaId: string
	readonly #ttlMs: number
	readonly #cancelRenewal: () => void
	#closed = false

	constructor(
		private readonly redis: RedisClient,
		options: RedisBindingOwnershipOptions = {}
	) {
		this.#replicaId = options.replicaId ?? Bun.env.HOSTNAME ?? crypto.randomUUID()
		this.#ttlMs = options.ttlMs ?? LOCK_TTL_MS
		const renewalIntervalMs = options.renewalIntervalMs ?? RENEWAL_INTERVAL_MS
		const timer = setInterval(() => {
			void this.renew().catch((error: unknown) => {
				console.error('youtube.lock-renewal.failed', { error })
			})
		}, renewalIntervalMs)
		this.#cancelRenewal = () => clearInterval(timer)
	}

	async tryAcquire(bindingId: string): Promise<boolean> {
		if (this.#closed) {
			return false
		}
		const result = await this.redis.set(this.#key(bindingId), this.#replicaId, 'PX', String(this.#ttlMs), 'NX')
		if (result !== 'OK') {
			return false
		}
		this.#ownedBindingIds.add(bindingId)
		return true
	}

	async renew(): Promise<void> {
		for (const bindingId of this.#ownedBindingIds) {
			try {
				const result = await this.redis.send('EVAL', [
					RENEW_IF_OWNER_SCRIPT,
					'1',
					this.#key(bindingId),
					this.#replicaId,
					String(this.#ttlMs),
				])
				if (result === 1 || result === '1') {
					continue
				}
			} catch (error) {
				const renewalError = error instanceof Error ? error : new Error('Redis lock renewal failed')
				console.error('youtube.lock-renewal.failed', { bindingId, error: renewalError })
			}
			this.#ownedBindingIds.delete(bindingId)
			for (const listener of this.#lostOwnershipListeners) {
				listener(bindingId)
			}
		}
	}

	async release(bindingId: string): Promise<void> {
		if (!this.#ownedBindingIds.has(bindingId)) {
			return
		}
		await this.redis.send('EVAL', [RELEASE_IF_OWNER_SCRIPT, '1', this.#key(bindingId), this.#replicaId])
		this.#ownedBindingIds.delete(bindingId)
	}

	onLostOwnership(listener: (bindingId: string) => void): () => void {
		this.#lostOwnershipListeners.add(listener)
		return () => this.#lostOwnershipListeners.delete(listener)
	}

	async close(): Promise<void> {
		if (this.#closed) {
			return
		}
		this.#closed = true
		this.#cancelRenewal()
		for (const bindingId of [...this.#ownedBindingIds]) {
			try {
				await this.release(bindingId)
			} catch (error) {
				const releaseError = error instanceof Error ? error : new Error('Redis lock release failed')
				console.error('youtube.lock-release.failed', { bindingId, error: releaseError })
			}
		}
		this.#lostOwnershipListeners.clear()
		this.redis.close()
	}

	#key(bindingId: string): string {
		return `${LOCK_PREFIX}${bindingId}`
	}
}
