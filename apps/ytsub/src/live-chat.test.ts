import { expect, test } from 'bun:test'
import { YTNodes } from 'youtubei.js'

import {
	LiveChatManager,
	type LiveChatSession,
	type LiveChatSource,
	StreamOfflineError,
} from './live-chat.ts'
import { RedisBindingOwnership, type RedisClient } from './locks.ts'
import { shouldIgnoreYoutubeBotMessage } from './message.ts'

import type { RetryScheduler } from './live-chat-scheduler.ts'
import type { ChannelBinding } from './message.ts'
import type { YoutubeStream } from './streams.ts'
import type { Helpers } from 'youtubei.js'

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

class FakeLiveChat implements LiveChatSession {
	startCalls = 0
	stopCalls = 0
	selectLiveChatCalls = 0
	#startListener: (() => void) | undefined
	#chatUpdateListener: ((action: Helpers.YTNode) => void) | undefined
	#errorListener: ((error: Error) => void) | undefined
	#endListener: (() => void) | undefined

	onStart(listener: () => void): void {
		this.#startListener = listener
	}

	onChatUpdate(listener: (action: Helpers.YTNode) => void): void {
		this.#chatUpdateListener = listener
	}

	onError(listener: (error: Error) => void): void {
		this.#errorListener = listener
	}

	onEnd(listener: () => void): void {
		this.#endListener = listener
	}

	selectLiveChat(): void {
		this.selectLiveChatCalls += 1
	}

	start(): void {
		this.startCalls += 1
		this.#startListener?.()
	}

	stop(): void {
		this.stopCalls += 1
	}

	emitError(error: Error): void {
		this.#errorListener?.(error)
	}

	emitChatUpdate(action: Helpers.YTNode): void {
		this.#chatUpdateListener?.(action)
	}

	emitEnd(): void {
		this.#endListener?.()
	}
}

class FakeRetryScheduler implements RetryScheduler {
	readonly delays: number[] = []
	readonly #tasks: Array<{ cancelled: boolean; callback: () => void }> = []

	schedule(delay: number, callback: () => void): () => void {
		const task = { cancelled: false, callback }
		this.delays.push(delay)
		this.#tasks.push(task)
		return () => {
			task.cancelled = true
		}
	}

	runNext(): void {
		const task = this.#tasks.shift()
		if (!task) {
			throw new Error('Expected a scheduled retry')
		}
		if (!task.cancelled) {
			task.callback()
		}
	}
}

function binding(platformChannelId = 'UCbroadcaster'): ChannelBinding {
	return {
		id: 'binding-1',
		channelId: 'channel-1',
		platformChannelId,
		botPlatformId: 'UCbot',
		userId: 'broadcaster-user-1',
		enabled: true,
	}
}

function stream(): YoutubeStream {
	return {
		videoId: 'video-1',
		broadcasterName: 'Broadcaster',
		title: 'Live title',
		viewers: 42,
		startedAt: new Date('2026-07-30T12:00:00.000Z'),
	}
}

function bus() {
	return {
		ChatMessages: { publish: async (): Promise<void> => {} },
		Parser: { ProcessMessageAsCommand: { publish: async (): Promise<void> => {} } },
	}
}

function deferred<T>(): { readonly promise: Promise<T>; resolve(value: T): void } {
	let resolvePromise: ((value: T) => void) | undefined
	const promise = new Promise<T>((resolve) => {
		resolvePromise = resolve
	})

	return {
		promise,
		resolve(value: T): void {
			if (!resolvePromise) {
				throw new Error('Deferred promise was not initialized')
			}
			resolvePromise(value)
		},
	}
}

test('LiveChatManager stops a stale session and starts its replacement when a binding changes during startup', async () => {
	const stale = new FakeLiveChat()
	const replacement = new FakeLiveChat()
	const firstResult = deferred<{ readonly session: LiveChatSession; readonly broadcasterName: string }>()
	const results = [firstResult.promise, Promise.resolve({ session: replacement, broadcasterName: 'Replacement' })]
	const source: LiveChatSource = {
		resolve(): Promise<{ readonly session: LiveChatSession; readonly broadcasterName: string }> {
			const result = results.shift()
			if (!result) {
				throw new Error('Unexpected live chat resolution')
			}
			return result
		},
	}
	const manager = new LiveChatManager(source, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
	})

	const firstSubscription = manager.subscribe(binding('UCold'))
	await manager.subscribe(binding('UCnew'))
	firstResult.resolve({ session: stale, broadcasterName: 'Stale' })
	await firstSubscription

	expect(replacement.startCalls).toBe(1)
	expect(replacement.selectLiveChatCalls).toBe(1)
	expect(stale.startCalls).toBe(0)
	expect(stale.stopCalls).toBe(1)
})

test('LiveChatManager reconciles its sessions to the authoritative binding snapshot', async () => {
	const removed = new FakeLiveChat()
	const added = new FakeLiveChat()
	const results = [
		Promise.resolve({ session: removed, broadcasterName: 'Removed' }),
		Promise.resolve({ session: added, broadcasterName: 'Added' }),
	]
	const source: LiveChatSource = {
		resolve(): Promise<{ readonly session: LiveChatSession; readonly broadcasterName: string }> {
			const result = results.shift()
			if (!result) {
				throw new Error('Unexpected live chat resolution')
			}
			return result
		},
	}
	const manager = new LiveChatManager(source, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
	})
	const removedBinding = binding('UCremoved')
	const addedBinding = { ...binding('UCadded'), id: 'binding-2' }

	await manager.subscribe(removedBinding)
	await manager.reconcile([addedBinding])

	expect(removed.stopCalls).toBe(1)
	expect(added.startCalls).toBe(1)
})

test('LiveChatManager preserves retry backoff until a chat update proves the session healthy', async () => {
	const first = new FakeLiveChat()
	const second = new FakeLiveChat()
	const scheduler = new FakeRetryScheduler()
	const results = [
		Promise.resolve({ session: first, broadcasterName: 'Broadcaster' }),
		Promise.resolve({ session: second, broadcasterName: 'Broadcaster' }),
	]
	const source: LiveChatSource = {
		resolve(): Promise<{ readonly session: LiveChatSession; readonly broadcasterName: string }> {
			const result = results.shift()
			if (!result) {
				throw new Error('Unexpected live chat resolution')
			}
			return result
		},
	}
	const manager = new LiveChatManager(source, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
		retryScheduler: scheduler,
	})

	await manager.subscribe(binding())
	first.emitError(new Error('first session failed'))
	scheduler.runNext()
	await Promise.resolve()
	second.emitError(new Error('second session failed'))

	expect(scheduler.delays).toEqual([1_000, 2_000])
})

test('shouldIgnoreYoutubeBotMessage suppresses bot replies but keeps broadcaster messages', () => {
	expect(shouldIgnoreYoutubeBotMessage(binding(), 'UCbot')).toBe(true)
	expect(shouldIgnoreYoutubeBotMessage(binding('UCbot'), 'UCbot')).toBe(false)
	expect(shouldIgnoreYoutubeBotMessage(binding(), 'UCviewer')).toBe(false)
})

function chatUpdateAction(authorId = 'UCviewer'): Helpers.YTNode {
	const item = {
		is: (node: unknown) => node === YTNodes.LiveChatTextMessage,
		id: 'message-1',
		message: { toString: () => 'hello' },
		author: { id: authorId, name: 'viewer', badges: [] },
	}

	return {
		is: (node: unknown) => node === YTNodes.AddChatItemAction,
		item,
	} as unknown as Helpers.YTNode
}

function flushMicrotasks(): Promise<void> {
	return new Promise((resolve) => setTimeout(resolve, 0))
}

test('LiveChatManager drops chat updates from a session replaced after an error', async () => {
	const stale = new FakeLiveChat()
	const replacement = new FakeLiveChat()
	const scheduler = new FakeRetryScheduler()
	const results = [
		Promise.resolve({ session: stale, broadcasterName: 'Broadcaster' }),
		Promise.resolve({ session: replacement, broadcasterName: 'Broadcaster' }),
	]
	const source: LiveChatSource = {
		resolve(): Promise<{ readonly session: LiveChatSession; readonly broadcasterName: string }> {
			const result = results.shift()
			if (!result) {
				throw new Error('Unexpected live chat resolution')
			}
			return result
		},
	}
	const published: unknown[] = []
	const fakeBus = {
		ChatMessages: { publish: async (message: unknown): Promise<void> => { published.push(message) } },
		Parser: { ProcessMessageAsCommand: { publish: async (): Promise<void> => {} } },
	}
	const manager = new LiveChatManager(source, fakeBus, {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
		retryScheduler: scheduler,
	})

	await manager.subscribe(binding())
	stale.emitError(new Error('stale session failed'))
	scheduler.runNext()
	await flushMicrotasks()

	stale.emitChatUpdate(chatUpdateAction())
	await flushMicrotasks()
	expect(published).toHaveLength(0)

	replacement.emitChatUpdate(chatUpdateAction())
	await flushMicrotasks()
	expect(published).toHaveLength(1)
})

test('LiveChatManager reconcile stops removed sessions even when an addition hangs', async () => {
	const removed = new FakeLiveChat()
	const hanging = deferred<{ readonly session: LiveChatSession; readonly broadcasterName: string }>()
	const results = [
		Promise.resolve({ session: removed, broadcasterName: 'Removed' }),
		hanging.promise,
	]
	const source: LiveChatSource = {
		resolve(): Promise<{ readonly session: LiveChatSession; readonly broadcasterName: string }> {
			const result = results.shift()
			if (!result) {
				throw new Error('Unexpected live chat resolution')
			}
			return result
		},
	}
	const manager = new LiveChatManager(source, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
	})

	await manager.subscribe(binding('UCremoved'))
	const reconcilePromise = manager.reconcile([{ ...binding('UCadded'), id: 'binding-2' }])
	await flushMicrotasks()

	expect(removed.stopCalls).toBe(1)

	hanging.resolve({ session: new FakeLiveChat(), broadcasterName: 'Added' })
	await reconcilePromise
})

test('LiveChatManager coalesces reconcile snapshots that arrive while one is in flight', async () => {
	const hanging = deferred<{ readonly session: LiveChatSession; readonly broadcasterName: string }>()
	const settled = new FakeLiveChat()
	const results = [hanging.promise, Promise.resolve({ session: settled, broadcasterName: 'Settled' })]
	const source: LiveChatSource = {
		resolve(): Promise<{ readonly session: LiveChatSession; readonly broadcasterName: string }> {
			const result = results.shift()
			if (!result) {
				throw new Error('Unexpected live chat resolution')
			}
			return result
		},
	}
	const manager = new LiveChatManager(source, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
	})

	const firstReconcile = manager.reconcile([binding('UCfirst')])
	const secondReconcile = manager.reconcile([{ ...binding('UCsecond'), id: 'binding-2' }])

	hanging.resolve({ session: new FakeLiveChat(), broadcasterName: 'First' })
	await firstReconcile
	await secondReconcile

	expect(settled.startCalls).toBe(1)
})

test('LiveChatManager stops a session after its Redis ownership is lost', async () => {
	const redis = new FakeRedis()
	const ownership = new RedisBindingOwnership(redis, { replicaId: 'replica-1' })
	const liveChat = new FakeLiveChat()
	const source: LiveChatSource = {
		resolve: async () => ({ session: liveChat, broadcasterName: 'Broadcaster' }),
	}
	const manager = new LiveChatManager(source, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
		ownership,
	})

	await manager.subscribe(binding())
	redis.values.set('ytsub:lock:binding-1', 'replica-2')
	await ownership.renew()
	await flushMicrotasks()

	expect(liveChat.stopCalls).toBe(1)
	await manager.close()
})

test('LiveChatManager acquires a released binding during a later reconcile', async () => {
	const redis = new FakeRedis()
	const firstOwnership = new RedisBindingOwnership(redis, { replicaId: 'replica-1' })
	const secondOwnership = new RedisBindingOwnership(redis, { replicaId: 'replica-2' })
	const firstSession = new FakeLiveChat()
	const secondSession = new FakeLiveChat()
	const first = new LiveChatManager({ resolve: async () => ({ session: firstSession, broadcasterName: 'First' }) }, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
		ownership: firstOwnership,
	})
	const second = new LiveChatManager({ resolve: async () => ({ session: secondSession, broadcasterName: 'Second' }) }, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
		ownership: secondOwnership,
	})
	const channelBinding = binding()

	await first.subscribe(channelBinding)
	await second.reconcile([channelBinding])
	expect(secondSession.startCalls).toBe(0)

	await first.unsubscribe(channelBinding.id)
	await second.reconcile([channelBinding])

	expect(secondSession.startCalls).toBe(1)
	await second.close()
})

test('LiveChatManager publishes stream online once and keeps it online across transient retries', async () => {
	const first = new FakeLiveChat()
	const second = new FakeLiveChat()
	const scheduler = new FakeRetryScheduler()
	const online: YoutubeStream[] = []
	const offline: YoutubeStream[] = []
	const results = [
		Promise.resolve({ session: first, broadcasterName: 'Broadcaster', stream: stream() }),
		Promise.resolve({ session: second, broadcasterName: 'Broadcaster', stream: stream() }),
	]
	const manager = new LiveChatManager({
		resolve(): Promise<{ readonly session: LiveChatSession; readonly broadcasterName: string; readonly stream: YoutubeStream }> {
			const result = results.shift()
			if (!result) {
				throw new Error('Unexpected live chat resolution')
			}
			return result
		},
	}, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
		retryScheduler: scheduler,
		onStreamOnline: async (_binding, activeStream): Promise<void> => { online.push(activeStream) },
		onStreamOffline: async (_binding, activeStream): Promise<void> => { offline.push(activeStream) },
	})

	await manager.subscribe(binding())
	first.emitError(new Error('temporary failure'))
	scheduler.runNext()
	await flushMicrotasks()

	expect(online).toHaveLength(1)
	expect(offline).toHaveLength(0)
})

test('LiveChatManager publishes stream offline when a live chat ends', async () => {
	const liveChat = new FakeLiveChat()
	const offline: YoutubeStream[] = []
	const manager = new LiveChatManager({ resolve: async () => ({ session: liveChat, broadcasterName: 'Broadcaster', stream: stream() }) }, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
		onStreamOffline: async (_binding, activeStream): Promise<void> => { offline.push(activeStream) },
	})

	await manager.subscribe(binding())
	liveChat.emitEnd()
	await flushMicrotasks()

	expect(offline).toHaveLength(1)
})

test('LiveChatManager publishes stream offline when a retry resolves offline', async () => {
	const liveChat = new FakeLiveChat()
	const scheduler = new FakeRetryScheduler()
	const offline: YoutubeStream[] = []
	const results = [
		Promise.resolve({ session: liveChat, broadcasterName: 'Broadcaster', stream: stream() }),
		Promise.reject(new StreamOfflineError('UCbroadcaster')),
	]
	const manager = new LiveChatManager({
		resolve(): Promise<{ readonly session: LiveChatSession; readonly broadcasterName: string; readonly stream?: YoutubeStream }> {
			const result = results.shift()
			if (!result) {
				throw new Error('Unexpected live chat resolution')
			}
			return result
		},
	}, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
		retryScheduler: scheduler,
		onStreamOffline: async (_binding, activeStream): Promise<void> => { offline.push(activeStream) },
	})

	await manager.subscribe(binding())
	liveChat.emitError(new Error('temporary failure'))
	scheduler.runNext()
	await flushMicrotasks()

	expect(offline).toHaveLength(1)
})

test('LiveChatManager publishes stream offline when unsubscribed', async () => {
	const liveChat = new FakeLiveChat()
	const offline: YoutubeStream[] = []
	const manager = new LiveChatManager({ resolve: async () => ({ session: liveChat, broadcasterName: 'Broadcaster', stream: stream() }) }, bus(), {
		ensureChatter: async (channelBinding): Promise<ChannelBinding> => channelBinding,
		onStreamOffline: async (_binding, activeStream): Promise<void> => { offline.push(activeStream) },
	})

	const channelBinding = binding()
	await manager.subscribe(channelBinding)
	await manager.unsubscribe(channelBinding.id)
	await flushMicrotasks()

	expect(offline).toHaveLength(1)
})
