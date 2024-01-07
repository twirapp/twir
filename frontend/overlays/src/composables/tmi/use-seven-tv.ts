import { useWebSocket } from '@vueuse/core';
import { defineStore, storeToRefs } from 'pinia';
import { ref, watch } from 'vue';

import { useEmotes } from './use-emotes.js';

import type { SevenTvChannelResponse, SevenTvEmote, SevenTvGlobalResponse } from '@/types.js';

// opcodes https://github.com/SevenTV/EventAPI#opcodes
const OPCODES = {
	// -- Server -> Client
	DISPATCH: 0,
	HELLO: 1,
	HEARTBEAT: 2,
	RECONNECT: 4,
	ACK: 5,
	ERROR: 6,
	END_OF_STREAM: 7,

	// -- Client -> Server
	IDENTIFY: 33,
	RESUME: 34,
	SUBSCRIBE: 35,
	UNSUBSCRIBE: 36,
	SIGNAL: 37,
	BRIDGE: 38,
} as const;

type Opcode = (typeof OPCODES)[keyof typeof OPCODES];

interface SevenTvWebsocketPayload {
	op: number;
	d: {
		type: string
		body: {
			id: string
			pulled?: SevenTvWebsocketBody[]
			pushed?: SevenTvWebsocketBody[]
			updated?: SevenTvWebsocketBody[]
		}
	};
}

interface SevenTvWebsocketBody {
	old_value?: SevenTvEmote;
	value: SevenTvEmote;
}

function sevenTvPayload(opcode: Opcode, data: any): string {
	return JSON.stringify({
		op: opcode,
		t: Date.now(),
		d: data,
	});
}

export const useSevenTv = defineStore('seven-tv', () => {
	const channelId = ref('');
	const sevenTvUserId = ref('');
	const currentEmoteSetId = ref('');

	const emotesStore = useEmotes();
	const { emotes } = storeToRefs(emotesStore);

	const { data, status, open, send, close } = useWebSocket('wss://events.7tv.io/v3', {
		immediate: false,
		autoReconnect: {
			delay: 500,
		},
	});

	function sendSevenTv(opcode: Opcode, data: any): void {
		send(sevenTvPayload(opcode, data));
	}

	watch(() => currentEmoteSetId.value, (emoteSetId) => {
		sendSevenTvEvent(OPCODES.SUBSCRIBE, emoteSetId);
	});

	function sendSevenTvEvent(opcode: Opcode, emoteSetId: string): void {
		sendSevenTv(opcode, {
			type: 'emote_set.update',
			condition: {
				object_id: emoteSetId,
			},
		});

		sendSevenTv(opcode, {
			type: 'user.update',
			condition: {
				object_id: sevenTvUserId.value,
			},
		});
	}

	watch(() => data.value, async (data) => {
		const parsedData = JSON.parse(data) as SevenTvWebsocketPayload;

		if (parsedData.op === OPCODES.DISPATCH) {
			const { type, body } = parsedData.d;

			if (type === 'emote_set.update') {
				// ADD
				if (body.pushed) {
					for (const emote of body.pushed) {
						emotesStore.updateSevenTvEmote(emote.value);
					}
					// UPDATE
				} else if (body.updated) {
					for (const emote of body.updated) {
						if (emote.old_value) {
							emotesStore.removeEmoteByName(emote.old_value.name);
						}

						emotesStore.updateSevenTvEmote(emote.value);
					}
					// REMOVE
				} else if (body.pulled) {
					for (const emote of body.pulled) {
						if (emote.old_value) {
							emotesStore.removeEmoteByName(emote.old_value.name);
						}
					}
				}
			}

			if (type === 'user.update') {
				// UPDATE
				if (body.updated) {
					// Unsubscribe from old emote set
					sendSevenTvEvent(OPCODES.UNSUBSCRIBE, currentEmoteSetId.value);

					// Delete all emotes
					for (const emote of Object.values(emotes)) {
						if (emote.service !== '7tv') continue;
						delete emotes.value[emote.name];
					}

					// Fetch actual emotes
					// eslint-disable-next-line @typescript-eslint/ban-ts-comment
					// @ts-ignore
					const newEmoteSetId = body.updated[0].value[0].value.id;
					const response = await fetch(`https://7tv.io/v3/emote-sets/${newEmoteSetId}?ts=${Date.now()}`);
					const newEmotes = (await response.json()) as SevenTvGlobalResponse;
					emotesStore.setSevenTvEmotes(newEmotes);
					currentEmoteSetId.value = newEmoteSetId;
				}
			}
		}
	});

	async function fetchSevenTvEmotes() {
		try {
			const [global, channel] = await Promise.all([
				fetch('https://7tv.io/v3/emote-sets/global'),
				fetch(`https://7tv.io/v3/users/twitch/${channelId.value}?ts=${Date.now()}`),
			]);

			const globalEmotes = (await global.json()) as SevenTvGlobalResponse;
			emotesStore.setSevenTvEmotes(globalEmotes);

			const channelEmotes = (await channel.json()) as SevenTvChannelResponse;
			emotesStore.setSevenTvEmotes(channelEmotes);

			currentEmoteSetId.value = channelEmotes.emote_set.id;
			sevenTvUserId.value = channelEmotes.user.id;
		} catch (err) {
			console.error(err);
		}
	}

	function connect(_channelId: string): void {
		if (status.value === 'OPEN') return;
		channelId.value = _channelId;
		open();
	}

	function destroy(): void {
		close();
	}

	return {
		connect,
		destroy,
		fetchSevenTvEmotes,
	};
});
