import { expect, test } from 'bun:test'
import { YTNodes } from 'youtubei.js'

import {
	LiveChatManager,
	type LiveChatSession,
	type LiveChatSource,
} from './live-chat.ts'
import { shouldIgnoreYoutubeBotMessage } from './message.ts'

import type { RetryScheduler } from './live-chat-scheduler.ts'
import type { ChannelBinding } from './message.ts'
import type { Helpers } from 'youtubei.js'

class FakeLiveChat implements LiveChatSession {
	startCalls = 0
	stopCalls = 0
	selectLiveChatCalls = 0
	#startListener: (() => void) | undefined
	#chatUpdateListener: ((action: Helpers.YTNode) => void) | undefined
	#errorListener: ((error: Error) => void) | undefined

	onStart(listener: () => void): void {
		this.#startListener = listener
	}

	onChatUpdate(listener: (action: Helpers.YTNode) => void): void {
		this.#chatUpdateListener = listener
	}

	onError(listener: (error: Error) => void): void {
		this.#errorListener = listener
	}

	onEnd(): void {}

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
