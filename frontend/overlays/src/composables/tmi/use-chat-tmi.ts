import { type ChatUserstate, Client } from 'tmi.js'
import { type Ref, unref, watch } from 'vue'

import { useMessageHelpers } from './use-message-helpers.js'
import { useThirdPartyEmotes } from './use-third-party-emotes.js'

import type { Message } from '@twir/frontend-chat'

declare module 'tmi.js' {
	export interface Events {
		usernotice: (msgId: string, channel: string, tags: ChatUserstate, msg: string) => void
	}
}

type MakeOptional<Type, Key extends keyof Type> = Omit<Type, Key> &
  Partial<Pick<Type, Key>>

export type ChatMessage = Omit<
	MakeOptional<Message, 'isItalic' | 'isAnnounce'>,
	'createdAt' | 'internalId'
> & {
	senderId?: string
	rawMessage?: string
	messageHideTimeout?: number
	messageShowDelay?: number
}

export const knownBots = new Set([
	'moobot',
	'fossabot',
	'wizebot',
	'twirapp',
	'nightbot',
	'streamlabs',
	'streamelements',
])

type MaybePromise = any | Promise<any>

export interface ChatSettings {
	channelName: string
	channelId: string
	emotes: ThirdPartyEmotesOptions
	onMessage: (message: ChatMessage) => MaybePromise
	onRemoveMessage?: (msgId: string) => MaybePromise
	onRemoveMessageByUser?: (userName: string) => MaybePromise
	onChatClear?: () => void
}

export interface ThirdPartyEmotesOptions {
	/**
	 * @default false
	 */
	isSmaller?: boolean
	sevenTv?: boolean
	bttv?: boolean
	ffz?: boolean
}

export function useChatTmi(options: Ref<ChatSettings>) {
	const { destroy: destroyThirdPartyEmotes } = useThirdPartyEmotes(options)
	const { makeMessageChunks } = useMessageHelpers()

	let client: Client | null = null

	function createMessage(chatMessage: ChatMessage) {
		const internalId = crypto.randomUUID()

		return {
			...chatMessage,
			isItalic: chatMessage.isItalic ?? false,
			createdAt: new Date(),
			internalId,
			isAnnounce: chatMessage.isAnnounce ?? false,
		}
	}

	function messageChunks(message: string, tags: ChatUserstate) {
		return makeMessageChunks(message, {
			isSmaller: options.value.emotes.isSmaller ?? false,
			emotesList: tags.emotes ?? {},
		})
	}

	async function destroy() {
		if (!client) return

		await client.disconnect()
		client.removeAllListeners()
		client = null
		destroyThirdPartyEmotes()
	}

	async function create(channel: string) {
		await destroy()

		client = new Client({
			connection: {
				secure: true,
				reconnect: true,
			},
			channels: [],
		})

		client.on('message', (_channel, tags, message) => {
			options.value.onMessage(createMessage({
				id: tags.id,
				type: 'message',
				rawMessage: message,
				chunks: messageChunks(message, tags),
				sender: tags.username,
				senderId: tags['user-id']!,
				senderColor: tags.color,
				senderDisplayName: tags['display-name'],
				badges: tags.badges as Record<string, string> | undefined,
				isItalic: tags['message-type'] === 'action',
			}))
		})

		client.on('usernotice', (msgId, _channel, tags, message) => {
			if (msgId !== 'announcement') return

			options.value.onMessage(createMessage({
				id: msgId,
				type: 'message',
				rawMessage: message,
				chunks: messageChunks(message, tags),
				sender: tags.login,
				senderId: tags['user-id'],
				senderColor: tags.color,
				senderDisplayName: tags['display-name'],
				badges: tags.badges as Record<string, string> | undefined,
				isItalic: tags['message-type'] === 'action',
				isAnnounce: true,
			}))
		})

		client.on('messagedeleted', (_channel, _username, _msgText, userState) => {
			const msgId = userState['target-msg-id']
			if (msgId) {
				options.value.onRemoveMessage?.(msgId)
			}
		})

		client.on('timeout', (_channel, username) => {
			options.value.onRemoveMessageByUser?.(username)
		})

		client.on('ban', (_channel, username) => {
			options.value.onRemoveMessageByUser?.(username)
		})

		client.on('clearchat', () => {
			if (options.value.onChatClear) {
				options.value.onChatClear()
			}
		})

		client.on('connecting', () => {
			options.value.onMessage(createMessage({
				type: 'system',
				chunks: [{
					type: 'text',
					value: 'Connecting to servers...',
				}],
				messageHideTimeout: 5,
			}))
		})

		client.on('connected', async () => {
			options.value.onMessage(createMessage({
				type: 'system',
				chunks: [{ type: 'text', value: 'Connected' }],
				messageHideTimeout: 6,
			}))

			await client!.join(channel)

			options.value.onMessage(createMessage({
				type: 'system',
				chunks: [{ type: 'text', value: `Joined channel ${channel}` }],
				messageHideTimeout: 7,
			}))
		})

		await client.connect()
	}

	watch(() => options.value.channelName, (v) => {
		const name = unref(v)
		if (!name) return

		create(name)
	})

	return {
		destroy,
	}
}
