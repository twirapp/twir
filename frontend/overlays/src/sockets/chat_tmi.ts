/* eslint-disable @typescript-eslint/ban-ts-comment */
import { Message } from '@twir/frontend-chat';
import { Client } from 'tmi.js';
import { Ref, onUnmounted, unref, watch } from 'vue';

import { makeMessageChunks } from '../components/chat_tmi_helpers.js';

type MakeOptional<Type, Key extends keyof Type> = Omit<Type, Key> &
  Partial<Pick<Type, Key>>;

export type ChatMessage = Omit<
	MakeOptional<Message, 'isItalic' | 'isAnnounce'>,
	'createdAt' | 'internalId'
> & { messageHideTimeout?: number, messageShowDelay?: number; }

export const knownBots = new Set([
	'moobot',
	'fossabot',
	'wizebot',
	'twirapp',
	'nightbot',
	'streamlabs',
	'streamelements',
]);

type MaybePromise = any | Promise<any>

export type ChatSettings = {
	channelName: string,
	channelId: string
	onMessage: (message: ChatMessage) => MaybePromise
	onRemoveMessage?: (msgId: string) => MaybePromise
	onRemoveMessageByUser?: (userName: string) => MaybePromise
	onChatClear?: () => void
}

// const { sevenTvEmotes, bttvEmotes, ffzEmotes } = useThirdPartyEmotes(channelName, channelId);
export const useTmiChat = (opts: Ref<ChatSettings>) => {
	let client: Client | null = null;

	onUnmounted(async () => {
		destroy();
	});

	const createMessage = (opts: ChatMessage) => {
		const internalId = crypto.randomUUID();

		return {
			...opts,
			isItalic: opts.isItalic ?? false,
			createdAt: new Date(),
			internalId,
			isAnnounce: opts.isAnnounce ?? false,
		};
	};

	async function destroy() {
		if (client) {
			await client.disconnect();
			client.removeAllListeners();
			client = null;
		}
	}

	async function create(channel: string) {
		await destroy();

		client = new Client({
			connection: {
				secure: true,
				reconnect: true,
			},
			channels: [],
		});

		client.on('message', (_, tags, message) => {
			opts.value.onMessage(createMessage({
				id: tags.id,
				type: 'message',
				chunks: makeMessageChunks(message, tags.emotes),
				sender: tags.username!,
				senderColor: tags.color,
				senderDisplayName: tags['display-name'],
				badges: tags.badges as Record<string, string> | undefined,
				isItalic: tags['message-type'] === 'action',
			}));
		});

		// @ts-ignore
		client.on('usernotice', (msgId, channel, tags, msg) => {
			if(msgId === 'announcement') {
				opts.value.onMessage(createMessage({
					id: msgId,
					type: 'message',
					// @ts-ignore
					chunks: makeMessageChunks(msg, tags.emotes),
					// @ts-ignore
					sender: tags.login,
					// @ts-ignore
					senderColor: tags.color,
					senderDisplayName: tags['display-name'],
					// @ts-ignore
					badges: tags.badges as Record<string, string> | undefined,
					isItalic: tags['message-type'] === 'action',
					isAnnounce: true,
				}));
			}
		});

		client.on('messagedeleted', (_channel, _username, _msgText, userState) => {
			const msgId = userState['target-msg-id'];
			if (msgId) {
				opts.value.onRemoveMessage?.(msgId);
			}
		});

		client.on('timeout', (_channel, username) => {
			opts.value.onRemoveMessageByUser?.(username);
		});

		client.on('ban', (_channel, username) => {
			opts.value.onRemoveMessageByUser?.(username);
		});

		client.on('clearchat', () => {
			if (opts.value.onChatClear) {
				opts.value.onChatClear();
			}
		});

		client.on('connecting', () => {
			opts.value.onMessage(createMessage({
				type: 'system',
				chunks: [{
					type: 'text',
					value: 'Connecting to servers...',
				}],
				messageHideTimeout: 5,
			}));
		});

		client.on('connected', async () => {
			opts.value.onMessage(createMessage({
				type: 'system',
				chunks: [{ type: 'text', value: 'Connected' }],
				messageHideTimeout: 6,
			}));

			await client!.join(channel);

			opts.value.onMessage(createMessage({
				type: 'system',
				chunks: [{ type: 'text', value: `Joined channel ${channel}` }],
				messageHideTimeout: 7,
			}));
		});

		await client.connect();
	}

	watch(() => opts.value.channelName, (v) => {
		const name = unref(v);
		if (!name) return;

		create(name);
	});

	return {
		destroy,
	};
};
