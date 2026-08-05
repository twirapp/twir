import emojiRegex from 'emoji-regex'

import { useEmotes } from '../tmi/use-emotes.js'

import type { MessageChunk } from '@twir/frontend-chat'

const emojiRegexp = emojiRegex()

export interface ChatEventFragment {
	type: string
	text: string
	emoteId?: string | null
	emoteUrl?: string | null
}

const escapeRegExp = (value: string) => value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')

/**
 * Twitch prefixes reply messages with the parent mention ("@parentName ...").
 * Strip that leading mention so only the actual reply text is displayed.
 */
export function stripReplyMention(
	fragments: readonly ChatEventFragment[],
	parentUserName: string,
): ChatEventFragment[] {
	const [first, ...rest] = fragments
	if (!first) return [...fragments]

	if (first.type === 'mention') {
		const [second, ...others] = rest
		if (second && second.type === 'text') {
			const text = second.text.replace(/^\s+/, '')
			return text ? [{ ...second, text }, ...others] : others
		}
		return rest
	}

	if (first.type !== 'text' || !parentUserName) return [...fragments]

	const mentionRegexp = new RegExp(`^@${escapeRegExp(parentUserName)}\\s*`, 'i')
	if (!mentionRegexp.test(first.text)) return [...fragments]

	const text = first.text.replace(mentionRegexp, '')
	return text ? [{ ...first, text }, ...rest] : rest
}

export function useFragmentsToChunks() {
	const { emotes: thirdPartyEmotes } = useEmotes()

	function fragmentsToChunks(fragments: readonly ChatEventFragment[]): MessageChunk[] {
		const chunks: MessageChunk[] = []

		for (const fragment of fragments) {
			if (fragment.type === 'emote' && fragment.emoteUrl) {
				chunks.push({
					type: 'emote',
					value: fragment.emoteUrl,
					emoteName: fragment.text,
				})
				continue
			}

			for (const word of fragment.text.split(' ')) {
				if (!word) {
					chunks.push({ type: 'text', value: '' })
					continue
				}

				const thirdPartyEmote = thirdPartyEmotes.value[word]
				if (thirdPartyEmote) {
					const isZeroWidthModifier = thirdPartyEmote.isZeroWidth
					const isModifier = typeof thirdPartyEmote.modifierFlag !== 'undefined'
					const url = thirdPartyEmote.urls.at(-1)!
					const latestChunk = chunks.at(-1)

					if (isZeroWidthModifier && latestChunk) {
						latestChunk.zeroWidthModifiers = [...(latestChunk.zeroWidthModifiers ?? []), url]
					} else if (isModifier && latestChunk) {
						latestChunk.flags = [
							...(latestChunk.flags ?? []),
							thirdPartyEmote.modifierFlag as number,
						]
					} else {
						chunks.push({
							type: '3rd_party_emote',
							value: url,
							emoteHeight: thirdPartyEmote.height,
							emoteWidth: thirdPartyEmote.width,
							emoteName: thirdPartyEmote.name,
						})
					}

					continue
				}

				if (word.match(emojiRegexp)) {
					chunks.push({ type: 'emoji', value: word })
					continue
				}

				chunks.push({ type: 'text', value: word })
			}
		}

		return chunks
	}

	return {
		fragmentsToChunks,
	}
}
