const LOCK_TTL_MS = 30_000
const ACQUIRE_OPERATION_TIMEOUT_MS = 5_000
const RENEWAL_INTERVAL_MS = 10_000
const RENEWAL_OPERATION_TIMEOUT_MS = 5_000
const LEASE_WATCHDOG_MS = 25_000
const ACQUIRE_ATTEMPTS = 1_200
const RETRY_DELAY_MS = 25

const RENEW_IF_OWNER_SCRIPT = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0`

const RELEASE_IF_OWNER_SCRIPT = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0`

export type OAuthProvider = 'streamelements' | 'streamlabs'

export interface RedisCommands {
	send(command: string, args: readonly string[]): Promise<unknown>
}

export interface OAuthRefreshLockOptions {
	readonly redis?: RedisCommands
	readonly createOwner?: () => string
	readonly acquireAttempts?: number
	readonly acquireOperationTimeoutMs?: number
	readonly retryDelayMs?: number
	readonly renewalIntervalMs?: number
	readonly renewalOperationTimeoutMs?: number
	readonly leaseWatchdogMs?: number
	readonly sleep?: (milliseconds: number) => Promise<void>
}

export class OAuthRefreshLockUnavailableError extends Error {
	constructor(options?: ErrorOptions) {
		super('OAuth refresh lock is unavailable', options)
		this.name = 'OAuthRefreshLockUnavailableError'
	}
}

export class OAuthRefreshLockLostError extends Error {
	constructor() {
		super('OAuth refresh lock was lost')
		this.name = 'OAuthRefreshLockLostError'
	}
}

function delay(milliseconds: number): Promise<void> {
	return new Promise((resolve) => setTimeout(resolve, milliseconds))
}

function isSuccessfulRedisInteger(value: unknown): boolean {
	return value === 1 || value === '1'
}

async function withDeadline<T>(
	operation: Promise<T>,
	milliseconds: number,
	timeoutError: Error = new OAuthRefreshLockLostError()
): Promise<T> {
	let timeout: ReturnType<typeof setTimeout> | undefined
	const deadline = new Promise<never>((_, reject) => {
		timeout = setTimeout(() => reject(timeoutError), milliseconds)
	})
	try {
		return await Promise.race([operation, deadline])
	} finally {
		if (timeout !== undefined) {
			clearTimeout(timeout)
		}
	}
}

export async function withOAuthRefreshLock<T>(
	provider: OAuthProvider,
	channelID: string,
	callback: (signal: AbortSignal) => Promise<T>,
	options: OAuthRefreshLockOptions = {}
): Promise<T> {
	const redis = options.redis ?? (await import('./redis.ts')).client
	const owner = (options.createOwner ?? (() => crypto.randomUUID()))()
	const key = `twir:integration-token-refresh:${provider}:${channelID}`
	const acquireAttempts = options.acquireAttempts ?? ACQUIRE_ATTEMPTS
	const acquireOperationTimeoutMs = options.acquireOperationTimeoutMs ?? ACQUIRE_OPERATION_TIMEOUT_MS
	const sleep = options.sleep ?? delay
	const retryDelayMs = options.retryDelayMs ?? RETRY_DELAY_MS

	let acquired = false
	for (let attempt = 0; attempt < acquireAttempts; attempt += 1) {
		let result: unknown
		try {
			result = await withDeadline(
				redis.send('SET', [key, owner, 'NX', 'PX', String(LOCK_TTL_MS)]),
				acquireOperationTimeoutMs,
				new OAuthRefreshLockUnavailableError()
			)
		} catch (cause) {
			throw new OAuthRefreshLockUnavailableError({ cause })
		}
		if (result === 'OK') {
			acquired = true
			break
		}
		if (attempt + 1 < acquireAttempts) {
			await sleep(retryDelayMs)
		}
	}
	if (!acquired) {
		throw new OAuthRefreshLockUnavailableError()
	}

	const renewalIntervalMs = options.renewalIntervalMs ?? RENEWAL_INTERVAL_MS
	const renewalOperationTimeoutMs = options.renewalOperationTimeoutMs ?? RENEWAL_OPERATION_TIMEOUT_MS
	const leaseWatchdogMs = options.leaseWatchdogMs ?? LEASE_WATCHDOG_MS
	const abortController = new AbortController()
	const callbackPromise = Promise.resolve().then(() => callback(abortController.signal))
	let lostError: OAuthRefreshLockLostError | undefined
	let rejectLost: (error: OAuthRefreshLockLostError) => void = () => undefined
	const lost = new Promise<never>((_, reject) => {
		rejectLost = reject
	})
	const loseLease = () => {
		if (lostError) return
		lostError = new OAuthRefreshLockLostError()
		abortController.abort(lostError)
		rejectLost(lostError)
	}

	let stopped = false
	let renewalPromise: Promise<void> | undefined
	let watchdog = setTimeout(loseLease, leaseWatchdogMs)
	const resetWatchdog = () => {
		if (stopped) return
		clearTimeout(watchdog)
		watchdog = setTimeout(loseLease, leaseWatchdogMs)
	}
	const renewal = setInterval(() => {
		if (renewalPromise || lostError || stopped) return
		renewalPromise = withDeadline(
			redis.send('EVAL', [RENEW_IF_OWNER_SCRIPT, '1', key, owner, String(LOCK_TTL_MS)]),
			renewalOperationTimeoutMs
		).then((result) => {
			if (stopped) return
			if (isSuccessfulRedisInteger(result)) {
				resetWatchdog()
			} else {
				loseLease()
			}
		}).catch(() => {
			if (!stopped) loseLease()
		}).finally(() => {
			renewalPromise = undefined
		})
	}, renewalIntervalMs)

	let value: T | undefined
	let operationError: Error | undefined
	try {
		value = await Promise.race([callbackPromise, lost])
		if (lostError) {
			throw new OAuthRefreshLockLostError()
		}
	} catch (error) {
		operationError = error instanceof Error
			? error
			: new Error('OAuth refresh callback failed', { cause: error })
		if (lostError) {
			await callbackPromise.catch(() => undefined)
		}
	} finally {
		stopped = true
		clearInterval(renewal)
		clearTimeout(watchdog)
		abortController.abort()
		if (renewalPromise) {
			await renewalPromise
		}
		try {
			await withDeadline(
				redis.send('EVAL', [RELEASE_IF_OWNER_SCRIPT, '1', key, owner]),
				RENEWAL_OPERATION_TIMEOUT_MS,
				new Error('OAuth refresh lock release timed out')
			)
		} catch (releaseError) {
			if (operationError === undefined) {
				operationError = new Error('Failed to release OAuth refresh lock', { cause: releaseError })
			}
		}
	}

	if (operationError) {
		if (operationError instanceof OAuthRefreshLockLostError) {
			throw new OAuthRefreshLockLostError()
		}
		throw new Error('OAuth refresh operation failed', { cause: operationError })
	}
	return value as T
}
