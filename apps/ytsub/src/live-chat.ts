import { YTNodes } from 'youtubei.js'

import type { ChatMessageBadge } from '@twir/bus-core'
import type { Helpers, Innertube, YT } from 'youtubei.js'

import { ensureYoutubeChatter } from './db.ts'
import { normalizeYoutubeTextMessage } from './message.ts'

import type { ChannelBinding, YoutubeTextChatMessage } from './message.ts'

const RETRY_BASE_MS = 1_000
const RETRY_MAX_MS = 120_000

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

export class LiveChatManager {
	readonly #bindings = new Map<string, ChannelBinding>()
	readonly #sessions = new Map<string, YT.LiveChat>()
	readonly #retryTimers = new Map<string, ReturnType<typeof setTimeout>>()
	readonly #starting = new Set<string>()
	readonly #attempts = new Map<string, number>()

	constructor(
		private readonly yt: Innertube,
		private readonly bus: Bus
	) {}

	async subscribe(binding: ChannelBinding): Promise<void> {
		const existing = this.#bindings.get(binding.id)
		if (existing?.platformChannelId !== binding.platformChannelId) {
			this.#stopSession(binding.id)
		}
		this.#bindings.set(binding.id, binding)
		await this.#start(binding)
	}

	unsubscribe(bindingId: string): void {
		this.#bindings.delete(bindingId)
		this.#attempts.delete(bindingId)
		this.#stopSession(bindingId)
	}

	async reconcile(): Promise<void> {
		await Promise.all([...this.#bindings.values()].map((binding) => this.#start(binding)))
	}

	close(): void {
		for (const bindingId of this.#bindings.keys()) {
			this.#stopSession(bindingId)
		}
		this.#bindings.clear()
	}

	async #start(binding: ChannelBinding): Promise<void> {
		if (this.#sessions.has(binding.id) || this.#starting.has(binding.id)) {
			return
		}
		this.#starting.add(binding.id)

		try {
			const endpoint = await this.yt.resolveURL(
				`https://www.youtube.com/channel/${binding.platformChannelId}/live`
			)
			const info = await this.yt.getInfo(endpoint)
			const liveChat = info.getLiveChat()
			const broadcasterName = info.basic_info.channel?.name ?? info.basic_info.author ?? binding.platformChannelId

			liveChat.on('chat-update', (action) => {
				void this.#handleChatUpdate(binding, broadcasterName, action)
			})
			liveChat.on('error', (error) => this.#finishWithRetry(binding, liveChat, error))
			liveChat.on('end', () => this.#finishWithRetry(binding, liveChat, new Error('live chat ended')))

			this.#sessions.set(binding.id, liveChat)
			this.#attempts.delete(binding.id)
			liveChat.start()
			console.info(`started YouTube live chat for ${binding.platformChannelId}`)
		} catch (error) {
			if (error instanceof Error) {
				this.#scheduleRetry(binding, error)
			} else {
				this.#scheduleRetry(binding, new Error('YouTube live chat startup failed'))
			}
		} finally {
			this.#starting.delete(binding.id)
		}
	}

	async #handleChatUpdate(
		binding: ChannelBinding,
		broadcasterName: string,
		action: Helpers.YTNode
	): Promise<void> {
		if (!action.is(YTNodes.AddChatItemAction)) {
			console.debug(`ignored YouTube live chat action ${action.type}`)
			return
		}
		if (!action.item.is(YTNodes.LiveChatTextMessage)) {
			console.debug(`ignored YouTube live chat item ${action.item.type}`)
			return
		}

		const sourceMessage = this.#toTextMessage(action.item)
		const chatterBinding = await ensureYoutubeChatter(binding, sourceMessage)
		const message = normalizeYoutubeTextMessage(chatterBinding, broadcasterName, sourceMessage)

		await this.bus.ChatMessages.publish(message)
		await this.bus.Parser.ProcessMessageAsCommand.publish(message)
	}

	#toTextMessage(item: YTNodes.LiveChatTextMessage): YoutubeTextChatMessage {
		const badges = item.author.badges.flatMap((badge): ChatMessageBadge[] => {
			if (!badge.is(YTNodes.LiveChatAuthorBadge)) {
				return []
			}
			return [{
				id: badge.icon_type,
				set_id: badge.style ?? '',
				text: badge.label ?? badge.tooltip ?? '',
			}]
		})

		return {
			id: item.id,
			text: item.message.toString(),
			author: {
				id: item.author.id,
				name: item.author.name,
				isModerator: item.author.is_moderator ?? false,
				badges,
			},
		}
	}

	#finishWithRetry(binding: ChannelBinding, liveChat: YT.LiveChat, error: Error): void {
		if (this.#sessions.get(binding.id) !== liveChat) {
			return
		}
		liveChat.stop()
		this.#sessions.delete(binding.id)
		this.#scheduleRetry(binding, error)
	}

	#scheduleRetry(binding: ChannelBinding, error: Error): void {
		if (!this.#bindings.has(binding.id) || this.#retryTimers.has(binding.id)) {
			return
		}
		const attempt = (this.#attempts.get(binding.id) ?? 0) + 1
		this.#attempts.set(binding.id, attempt)
		const delay = Math.min(RETRY_BASE_MS * 2 ** (attempt - 1), RETRY_MAX_MS)
		console.warn(`YouTube live chat unavailable for ${binding.platformChannelId}; retrying in ${delay}ms`, error)
		const timer = setTimeout(() => {
			this.#retryTimers.delete(binding.id)
			const activeBinding = this.#bindings.get(binding.id)
			if (activeBinding) {
				void this.#start(activeBinding)
			}
		}, delay)
		this.#retryTimers.set(binding.id, timer)
	}

	#stopSession(bindingId: string): void {
		const retryTimer = this.#retryTimers.get(bindingId)
		if (retryTimer) {
			clearTimeout(retryTimer)
			this.#retryTimers.delete(bindingId)
		}
		const liveChat = this.#sessions.get(bindingId)
		if (liveChat) {
			liveChat.stop()
			this.#sessions.delete(bindingId)
		}
	}
}
