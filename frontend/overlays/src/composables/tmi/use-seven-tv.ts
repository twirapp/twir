import { createGlobalState, useWebSocket } from '@vueuse/core'
import { ref, watch } from 'vue'

import { useEmotes } from './use-emotes.js'

import type { SevenTvChannelResponse, SevenTvEmote, SevenTvGlobalResponse } from '@/types.js'

import { requestWithOutCache } from '@/helpers.js'

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
} as const

type Opcode = (typeof OPCODES)[keyof typeof OPCODES]

interface SevenTvWebsocketPayload {
	op: number
	d: {
		type: string
		body: {
			id: string
			pulled?: SevenTvWebsocketBody[]
			pushed?: SevenTvWebsocketBody[]
			updated?: SevenTvWebsocketBody[]
		}
	}
}

interface SevenTvWebsocketBody {
	old_value?: SevenTvEmote
	value: SevenTvEmote
}

function sevenTvPayload(opcode: Opcode, data: any): string {
	return JSON.stringify({
		op: opcode,
		t: Date.now(),
		d: data,
	})
}

export const useSevenTv = createGlobalState(() => {
	const countRefetch = ref(0)
	const sevenTvUserId = ref('')
	const currentEmoteSetId = ref('')

	const {
		emotes,
		removeEmoteByName,
		setSevenTvEmotes,
		updateSevenTvEmote,
	} = useEmotes()

	const { data, status, open, send, close } = useWebSocket('wss://events.7tv.io/v3', {
		immediate: false,
		autoReconnect: {
			delay: 500,
		},
	})

	function sendSevenTv(opcode: Opcode, data: any): void {
		send(sevenTvPayload(opcode, data))
	}

	watch(() => currentEmoteSetId.value, (emoteSetId) => {
		sendSevenTvEvent(OPCODES.SUBSCRIBE, emoteSetId)
	})

	function sendSevenTvEvent(opcode: Opcode, emoteSetId: string): void {
		sendSevenTv(opcode, {
			type: 'emote_set.update',
			condition: {
				object_id: emoteSetId,
			},
		})

		sendSevenTv(opcode, {
			type: 'user.update',
			condition: {
				object_id: sevenTvUserId.value,
			},
		})
	}

	watch(() => data.value, async (data) => {
		const parsedData = JSON.parse(data) as SevenTvWebsocketPayload

		if (parsedData.op === OPCODES.DISPATCH) {
			const { type, body } = parsedData.d

			if (type === 'emote_set.update') {
				// ADD
				if (body.pushed) {
					for (const emote of body.pushed) {
						updateSevenTvEmote(emote.value)
					}
					// UPDATE
				} else if (body.updated) {
					for (const emote of body.updated) {
						if (emote.old_value) {
							removeEmoteByName(emote.old_value.name)
						}

						updateSevenTvEmote(emote.value)
					}
					// REMOVE
				} else if (body.pulled) {
					for (const emote of body.pulled) {
						if (emote.old_value) {
							removeEmoteByName(emote.old_value.name)
						}
					}
				}
			}

			if (type === 'user.update') {
				// UPDATE
				if (body.updated) {
					// Unsubscribe from old emote set
					sendSevenTvEvent(OPCODES.UNSUBSCRIBE, currentEmoteSetId.value)

					// Delete all emotes
					for (const emote of Object.values(emotes)) {
						if (emote.service !== '7tv') continue
						delete emotes.value[emote.name]
					}

					// Fetch actual emotes
					// eslint-disable-next-line ts/ban-ts-comment
					// @ts-expect-error
					const newEmoteSetId = body.updated[0].value[0].value.id
					const newEmotes = await requestWithOutCache<SevenTvGlobalResponse>(
						`https://7tv.io/v3/emote-sets/${newEmoteSetId}`,
					)
					setSevenTvEmotes(newEmotes)
					currentEmoteSetId.value = newEmoteSetId
				}
			}
		}
	})

	async function fetchSevenTvEmotes(channelId: string) {
		if (status.value === 'OPEN' || countRefetch.value > 3) return

		try {
			const [globalEmotes, channelEmotes] = await Promise.all([
				requestWithOutCache<SevenTvGlobalResponse>(
					'https://7tv.io/v3/emote-sets/global',
				),
				requestWithOutCache<SevenTvChannelResponse>(
					`https://7tv.io/v3/users/twitch/${channelId}`,
				),
			])

			setSevenTvEmotes(globalEmotes)
			setSevenTvEmotes(channelEmotes)

			currentEmoteSetId.value = channelEmotes.emote_set.id
			sevenTvUserId.value = channelEmotes.user.id

			open()
		} catch (err) {
			countRefetch.value++
			fetchSevenTvEmotes(channelId)
			console.error(err)
		}
	}

	function destroy(): void {
		if (status.value !== 'OPEN') return
		close()
	}

	return {
		destroy,
		fetchSevenTvEmotes,
	}
})
