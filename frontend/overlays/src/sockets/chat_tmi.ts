/* eslint-disable @typescript-eslint/ban-ts-comment */
import { Client } from 'tmi.js';
import { Ref, onUnmounted, ref, unref, watch } from 'vue';

import { useThirdPartyEmotes } from './chat_tmi_emotes.js';
import { makeMessageChunks } from './chat_tmi_helpers.js';

export type MessageChunk = {
	type: 'text' | 'emote' | '7tv_emote' | 'bttv_emote' | 'ffz_emote';
	value: string;
}

type MakeOptional<Type, Key extends keyof Type> = Omit<Type, Key> &
  Partial<Pick<Type, Key>>;

type Message = {
	internalId: string,
	id?: string,
	type: string,
	chunks: MessageChunk[],
	sender?: string,
	senderColor?: string,
	senderDisplayName?: string
	badges?: Record<string, string>,
	isItalic: boolean;
	createdAt: Date;
};

type AddMessageOpts = Omit<
	MakeOptional<Message, 'isItalic'>,
	'createdAt' | 'internalId'
> & { messageTimeout?: number }

export const useTmiChat = (
	channelName: Ref<string>,
	channelId: Ref<string>,
	messageTimeout: Ref<number>,
) => {
	let client: Client | null = null;
	const messages = ref<Message[]>([]);
	useThirdPartyEmotes(channelName, channelId);

	onUnmounted(async () => {
		destroy();
	});

	function addMessage(opts: AddMessageOpts) {
		const internalId = crypto.randomUUID();

		messages.value.push({
			...opts,
			isItalic: opts.isItalic ?? false,
			createdAt: new Date(),
			internalId,
		});

		const timeout = opts.messageTimeout ?? messageTimeout.value;

		if (timeout) {
			setTimeout(() => {
				removeMessageByInternalId(internalId);
			}, timeout * 1000);
		}
	}

	function removeMessageByInternalId(id: string) {
		messages.value = messages.value.filter((m) => m.internalId !== id);
	}

	function removeMessageById(id: string) {
		messages.value = messages.value.filter((m) => m.id !== id);
	}

	function removeMessageBySenderName(name: string) {
		messages.value = messages.value.filter((m) => m.sender !== name);
	}

	async function destroy() {
		if (client) {
			await client.disconnect();
			client.removeAllListeners();
			client = null;
		}
	}

	async function create(channel: string) {
		destroy();

		client = new Client({
			connection: {
				secure: true,
				reconnect: true,
			},
			channels: [],
		});

		client.on('message', (_, tags, message) => {
			addMessage({
				id: tags.id,
				type: 'message',
				chunks: makeMessageChunks(message, tags.emotes),
				sender: tags.username!,
				senderColor: tags.color,
				senderDisplayName: tags['display-name'],
				badges: tags.badges as Record<string, string> | undefined,
				isItalic: tags['message-type'] === 'action',
			});
		});

		// @ts-ignore
		client.on('usernotice', (msgId, channel, tags, msg) => {
			if(msgId === 'announcement') {
				addMessage({
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
				});
			}
		});

		client.on('messagedeleted', (_channel, _username, _msgText, userState) => {
			const msgId = userState['target-msg-id'];
			if (msgId) {
				removeMessageById(msgId);
			}
		});

		client.on('timeout', (_channel, username) => {
			removeMessageBySenderName(username);
		});

		client.on('ban', (_channel, username) => {
			removeMessageBySenderName(username);
		});

		client.on('connecting', () => {
			addMessage({
				type: 'info',
				chunks: [{
					type: 'text',
					value: 'Connecting to servers...',
				}],
				messageTimeout: 5,
			});
		});

		client.on('connected', async () => {
			addMessage({
				type: 'info',
				chunks: [{ type: 'text', value: 'Connected' }],
				messageTimeout: 6,
			});

			await client!.join(channel);
			addMessage({
				type: 'info',
				chunks: [{ type: 'text', value: `Joined channel ${channel}` }],
				messageTimeout: 7,
			});
		});

		await client.connect();
	}

	watch(channelName, () => {
		const name = unref(channelName);
		if (!name) return;

		create(name);
	});

	return {
		messages,
		messageTimeout,
	};
};
