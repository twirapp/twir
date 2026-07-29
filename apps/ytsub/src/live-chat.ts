import { YTNodes } from 'youtubei.js'

import type { Helpers } from 'youtubei.js'

import { defaultRetryScheduler } from './live-chat-scheduler.ts'
import {
	areYoutubeBindingsEqual,
	normalizeYoutubeTextMessage,
	shouldIgnoreYoutubeBotMessage,
	toYoutubeTextChatMessage,
} from './message.ts'

import type { ChannelBinding, YoutubeTextChatMessage } from './message.ts'
import type { BindingOwnershipManager } from './locks.ts'
import type { RetryScheduler } from './live-chat-scheduler.ts'

const RETRY_BASE_MS = 1_000
const RETRY_MAX_MS = 120_000

export class StreamOfflineError extends Error {
	constructor(platformChannelId: string) {
		super(`YouTube stream ${platformChannelId} is offline`)
		this.name = 'StreamOfflineError'
	}
}

type Bus = {
	readonly ChatMessages: {
		publish(message: ReturnType<typeof normalizeYoutubeTextMessage>): Promise<void>
	}
	readonly Parser: {
		readonly ProcessMessageAsCommand: {
			publish(message: ReturnType<typeof normalizeYoutubeTextMessage>): Promise<void>
		}
	}
}

export interface LiveChatSession {
	onStart(listener: () => void): void
	onChatUpdate(listener: (action: Helpers.YTNode) => void): void
	onError(listener: (error: Error) => void): void
	onEnd(listener: () => void): void
	selectLiveChat(): void
	start(): void
	stop(): void
}

export interface LiveChatSource {
	resolve(binding: ChannelBinding): Promise<{
		readonly session: LiveChatSession
		readonly broadcasterName: string
	}>
}

export interface LiveChatManagerOptions {
	readonly ensureChatter: (
		binding: ChannelBinding,
		message: YoutubeTextChatMessage
	) => Promise<ChannelBinding>
	readonly retryScheduler?: RetryScheduler
	readonly ownership?: BindingOwnershipManager
}

interface ActiveSession {
	readonly generation: number
	readonly session: LiveChatSession
}

interface StartingSession {
	readonly generation: number
	readonly promise: Promise<void>
}

function logAsyncError(event: string, error: unknown, bindingId: string): void {
	console.error(event, { bindingId, error })
}

export class LiveChatManager {
	readonly #bindings = new Map<string, ChannelBinding>()
	readonly #sessions = new Map<string, ActiveSession>()
	readonly #retryTimers = new Map<string, () => void>()
	readonly #starting = new Map<string, StartingSession>()
	readonly #attempts = new Map<string, number>()
	readonly #generations = new Map<string, number>()
	#closed = false
	#reconciling = false
	#pendingReconcile: readonly ChannelBinding[] | null = null
	readonly #removeOwnershipListener: (() => void) | undefined

	constructor(
		private readonly source: LiveChatSource,
		private readonly bus: Bus,
		private readonly options: LiveChatManagerOptions
	) {
		this.#removeOwnershipListener = options.ownership?.onLostOwnership((bindingId) => {
			void this.#handleOwnershipLoss(bindingId).catch((error: unknown) => {
				logAsyncError('youtube.ownership-loss.failed', error, bindingId)
			})
		})
	}

	async subscribe(binding: ChannelBinding): Promise<void> {
		if (this.#closed) {
			return
		}

		const existing = this.#bindings.get(binding.id)
		if (existing && areYoutubeBindingsEqual(existing, binding)) {
			await this.#start(binding, this.#generation(binding.id))
			return
		}
		if (!existing && this.options.ownership && !(await this.options.ownership.tryAcquire(binding.id))) {
			return
		}

		const generation = this.#nextGeneration(binding.id)
		this.#attempts.delete(binding.id)
		this.#stopSession(binding.id)
		this.#bindings.set(binding.id, binding)
		await this.#start(binding, generation)
	}

	async unsubscribe(bindingId: string): Promise<void> {
		if (!this.#bindings.has(bindingId)) {
			return
		}
		this.#forgetBinding(bindingId)
		await this.options.ownership?.release(bindingId)
	}

	#forgetBinding(bindingId: string): void {
		this.#nextGeneration(bindingId)
		this.#bindings.delete(bindingId)
		this.#attempts.delete(bindingId)
		this.#stopSession(bindingId)
	}

	async reconcile(bindings: readonly ChannelBinding[]): Promise<void> {
		if (this.#reconciling) {
			this.#pendingReconcile = bindings
			return
		}
		this.#reconciling = true
		try {
			let current = bindings
			for (;;) {
				const desired = new Map(current.map((binding) => [binding.id, binding]))
				for (const bindingId of [...this.#bindings.keys()]) {
					if (!desired.has(bindingId)) {
						await this.unsubscribe(bindingId)
					}
				}
				for (const binding of current) {
					await this.subscribe(binding)
				}
				if (this.#pendingReconcile === null) {
					break
				}
				current = this.#pendingReconcile
				this.#pendingReconcile = null
			}
		} finally {
			this.#reconciling = false
		}
	}

	async close(): Promise<void> {
		if (this.#closed) {
			return
		}
		this.#closed = true
		for (const bindingId of [...this.#bindings.keys()]) {
			this.#forgetBinding(bindingId)
		}
		this.#attempts.clear()
		this.#removeOwnershipListener?.()
		await this.options.ownership?.close()
	}

	async #start(binding: ChannelBinding, generation: number): Promise<void> {
		if (!this.#isCurrent(binding, generation)
			|| this.#sessions.has(binding.id)
			|| this.#retryTimers.has(binding.id)) {
			return
		}
		const starting = this.#starting.get(binding.id)
		if (starting?.generation === generation) {
			await starting.promise
			return
		}

		const promise = this.#startSession(binding, generation)
		this.#starting.set(binding.id, { generation, promise })
		try {
			await promise
		} finally {
			if (this.#starting.get(binding.id)?.generation === generation) {
				this.#starting.delete(binding.id)
			}
		}
	}

	async #startSession(binding: ChannelBinding, generation: number): Promise<void> {
		try {
			const resolved = await this.source.resolve(binding)
			if (!this.#isCurrent(binding, generation)) {
				resolved.session.stop()
				return
			}

			resolved.session.onStart(() => {
				if (this.#isCurrent(binding, generation)) {
					resolved.session.selectLiveChat()
				}
			})
			resolved.session.onChatUpdate((action) => {
				void this.#handleChatUpdate(binding, generation, resolved.session, resolved.broadcasterName, action)
					.catch((error: unknown) => logAsyncError('youtube.chat-update.failed', error, binding.id))
			})
			resolved.session.onError((error) => {
				this.#finishWithRetry(binding, generation, resolved.session, error)
			})
			resolved.session.onEnd(() => {
				this.#finishWithRetry(binding, generation, resolved.session, new StreamOfflineError(binding.platformChannelId))
			})
			if (!this.#isCurrent(binding, generation)) {
				resolved.session.stop()
				return
			}

			this.#sessions.set(binding.id, { generation, session: resolved.session })
			resolved.session.start()
			console.info(`started YouTube live chat for ${binding.platformChannelId}`)
		} catch (error) {
			if (this.#isCurrent(binding, generation)) {
				this.#stopSession(binding.id, generation)
				if (error instanceof StreamOfflineError) {
					this.#scheduleRetry(binding, generation, error, true)
				} else {
					this.#scheduleRetry(binding, generation, error instanceof Error ? error : new Error('YouTube live chat startup failed'), false)
				}
			}
		}
	}

	async #handleChatUpdate(
		binding: ChannelBinding,
		generation: number,
		session: LiveChatSession,
		broadcasterName: string,
		action: Helpers.YTNode
	): Promise<void> {
		if (!this.#isCurrentSession(binding, generation, session)) {
			return
		}
		this.#attempts.delete(binding.id)
		if (!action.is(YTNodes.AddChatItemAction)) {
			console.debug(`ignored YouTube live chat action ${action.type}`)
			return
		}
		if (!action.item.is(YTNodes.LiveChatTextMessage)) {
			console.debug(`ignored YouTube live chat item ${action.item.type}`)
			return
		}

		const sourceMessage = toYoutubeTextChatMessage(action.item)
		if (shouldIgnoreYoutubeBotMessage(binding, sourceMessage.author.id)) {
			return
		}
		const chatterBinding = await this.options.ensureChatter(binding, sourceMessage)
		if (!this.#isCurrentSession(binding, generation, session)) {
			return
		}
		const message = normalizeYoutubeTextMessage(chatterBinding, broadcasterName, sourceMessage)

		console.info(`received YouTube live chat message ${message.message_id} from ${message.sender_display_name} in ${binding.platformChannelId}`)

		await this.bus.ChatMessages.publish(message)
		if (!this.#isCurrentSession(binding, generation, session)) {
			return
		}
		await this.bus.Parser.ProcessMessageAsCommand.publish(message)
	}

	#finishWithRetry(binding: ChannelBinding, generation: number, session: LiveChatSession, error: Error): void {
		const active = this.#sessions.get(binding.id)
		if (!this.#isCurrent(binding, generation) || active?.generation !== generation || active.session !== session) {
			return
		}
		this.#stopSession(binding.id, generation)
		this.#scheduleRetry(binding, generation, error, error instanceof StreamOfflineError)
	}

	#scheduleRetry(binding: ChannelBinding, generation: number, error: Error, quiet = false): void {
		if (!this.#isCurrent(binding, generation) || this.#retryTimers.has(binding.id)) {
			return
		}
		const attempt = (this.#attempts.get(binding.id) ?? 0) + 1
		this.#attempts.set(binding.id, attempt)
		const delay = Math.min(RETRY_BASE_MS * 2 ** (attempt - 1), RETRY_MAX_MS)
		if (quiet) {
			console.debug(`YouTube stream offline for ${binding.platformChannelId}; retrying in ${delay}ms`)
		} else {
			console.warn(`YouTube live chat unavailable for ${binding.platformChannelId}; retrying in ${delay}ms`, error)
		}
		const cancel = (this.options.retryScheduler ?? defaultRetryScheduler).schedule(delay, () => {
			this.#retryTimers.delete(binding.id)
			void this.#start(binding, generation)
				.catch((retryError: unknown) => logAsyncError('youtube.retry-start.failed', retryError, binding.id))
		})
		this.#retryTimers.set(binding.id, cancel)
	}

	#stopSession(bindingId: string, generation?: number): void {
		if (generation === undefined || this.#generation(bindingId) === generation) {
			const cancelRetry = this.#retryTimers.get(bindingId)
			if (cancelRetry) {
				cancelRetry()
				this.#retryTimers.delete(bindingId)
			}
		}
		const active = this.#sessions.get(bindingId)
		if (active && (generation === undefined || active.generation === generation)) {
			active.session.stop()
			this.#sessions.delete(bindingId)
		}
	}

	#isCurrent(binding: ChannelBinding, generation: number): boolean {
		return !this.#closed
			&& this.#generation(binding.id) === generation
			&& this.#bindings.get(binding.id) === binding
	}

	#isCurrentSession(binding: ChannelBinding, generation: number, session: LiveChatSession): boolean {
		return this.#isCurrent(binding, generation)
			&& this.#sessions.get(binding.id)?.session === session
	}

	#generation(bindingId: string): number {
		return this.#generations.get(bindingId) ?? 0
	}

	#nextGeneration(bindingId: string): number {
		const generation = this.#generation(bindingId) + 1
		this.#generations.set(bindingId, generation)
		return generation
	}

	async #handleOwnershipLoss(bindingId: string): Promise<void> {
		if (!this.#bindings.has(bindingId)) {
			return
		}
		this.#forgetBinding(bindingId)
	}
}
