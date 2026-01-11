import { computed, type ComputedRef } from 'vue'

import { useEmotes } from '@/composables/tmi/use-emotes.js'

export interface TextChunk {
	type: 'text' | 'emote'
	value: string
	emoteWidth?: number
	emoteHeight?: number
}

export function useBrbTextParser(text: ComputedRef<string | null>) {
	const { emotes } = useEmotes()

	const chunks = computed<TextChunk[]>(() => {
		if (!text.value) return []

		const words = text.value.split(' ')
		const result: TextChunk[] = []

		for (const word of words) {
			const emote = emotes.value[word]

			if (emote) {
				// Use the highest quality emote (last URL in the array)
				const emoteUrl = emote.urls[emote.urls.length - 1]
				result.push({
					type: 'emote',
					value: emoteUrl,
					emoteWidth: emote.width,
					emoteHeight: emote.height,
				})
			} else {
				// If the last chunk is text, append to it with a space
				const lastChunk = result[result.length - 1]
				if (lastChunk && lastChunk.type === 'text') {
					lastChunk.value += ' ' + word
				} else {
					result.push({
						type: 'text',
						value: word,
					})
				}
			}
		}

		return result
	})

	return {
		chunks,
	}
}
