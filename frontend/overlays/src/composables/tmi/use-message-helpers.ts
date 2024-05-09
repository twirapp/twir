import { createGlobalState } from '@vueuse/core'
import emojiRegex from 'emoji-regex'

import { useEmotes } from './use-emotes.js'

import type { MessageChunk } from '@twir/frontend-chat'

const emojiRegexp = emojiRegex()

interface TwitchEmotes {
	[emoteId: string]: string
}

export const useMessageHelpers = createGlobalState(() => {
	const { emotes: thirdPartyEmotes } = useEmotes()

	function makeMessageChunks(
		message: string,
		emotes: {
			isSmaller: boolean
			emotesList: Record<string, string[]>
		},
	): MessageChunk[] {
		const chunks: MessageChunk[] = []
		const parsedTwitchEmotes: TwitchEmotes = {}

		for (const [emoteId, positions] of Object.entries(emotes.emotesList)) {
			for (const position of positions) {
				const [from] = position.split('-').map(Number)
				parsedTwitchEmotes[from] = emoteId
			}
		}

		let currentWordIndex = 0
		for (const part of message.split(' ')) {
			const emote = parsedTwitchEmotes[currentWordIndex]
			const thirdPartyEmote = thirdPartyEmotes.value[part]
			const emojiMatch = part.match(emojiRegexp)

			if (emojiMatch && emojiMatch.length) {
				chunks.push({
					type: 'emoji',
					value: part,
				})
			} else if (emote) {
				chunks.push({
					type: 'emote',
					value: `https://static-cdn.jtvnw.net/emoticons/v2/${emote}/default/dark/${emotes.isSmaller ? 1 : 3}.0`,
				})
			} else if (thirdPartyEmote) {
				const isZeroWidthModifier = thirdPartyEmote.isZeroWidth
				const isModifier = typeof thirdPartyEmote.modifierFlag !== 'undefined'
				const url = thirdPartyEmote.urls.at(emotes.isSmaller ? 0 : -1)!
				const latestChunk = chunks.at(-1)!

				if (isZeroWidthModifier && latestChunk) {
					latestChunk.zeroWidthModifiers = [...(latestChunk.zeroWidthModifiers ?? []), url]
				} else if (isModifier && latestChunk) {
					const flags = [
						...(latestChunk.flags ?? []),
						thirdPartyEmote.modifierFlag as number,
					]
					latestChunk.flags = flags
				} else {
					chunks.push({
						type: '3rd_party_emote',
						value: url,
						emoteHeight: thirdPartyEmote.height,
						emoteWidth: thirdPartyEmote.width,
						emoteName: thirdPartyEmote.name,
					})
				}
			} else {
				chunks.push({ type: 'text', value: part })
			}

			currentWordIndex = currentWordIndex + part.length + 1
		}

		return chunks
	}

	return {
		makeMessageChunks,
	}
})
