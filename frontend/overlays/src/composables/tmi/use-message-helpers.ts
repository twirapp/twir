import { type MessageChunk } from '@twir/frontend-chat/dist';
import emojiRegex from 'emoji-regex';
import { defineStore, storeToRefs } from 'pinia';

import { useEmotes } from './use-emotes.ts';

const emojiRegexp = emojiRegex();

export const useMessageHelpers = defineStore('message-helpers', () => {
	const emotesStore = useEmotes();
	const { emotes: thirdPartyEmotes } = storeToRefs(emotesStore);

	function makeMessageChunks(
		message: string,
		emotes?: {
			[emoteid: string]: string[];
		},
	): MessageChunk[] {
		const parsedTwitchEmotes = emotes
			? Object.entries(emotes).reduce((acc, [id, positions]) => {
				positions.forEach((position) => {
					const [from, to] = position.split('-').map(Number);
					acc.push({
						from,
						to,
						emoteId: id,
					});
				});
				return acc;
			}, [] as { from: number; to: number; emoteId: string }[])
			: [];

		const chunks: MessageChunk[] = [];

		let currentWordIndex = 0;
		for (const part of message.split(' ')) {
			const emote = parsedTwitchEmotes.find((e) => e.from === currentWordIndex);
			const thirdPartyEmote = thirdPartyEmotes.value[part];
			const emojiMatch = part.match(emojiRegexp);

			if (emojiMatch && emojiMatch.length) {
				chunks.push({
					type: 'emoji',
					value: part,
				});
			} else if (emote) {
				chunks.push({
					type: 'emote',
					value: emote.emoteId,
				});
			} else if (thirdPartyEmote) {
				const isZeroWidthModifier = thirdPartyEmote.isZeroWidth;
				const isModifier = typeof thirdPartyEmote.modifierFlag !== 'undefined';
				const url = thirdPartyEmote.urls.at(-1)!;
				const latestChunk = chunks.at(-1)!;

				if (isZeroWidthModifier) {
					chunks.at(-1)!.zeroWidthModifiers = [...(chunks.at(-1)!.zeroWidthModifiers ?? []), url];
				} else if (isModifier && latestChunk) {
					const flags = [
						...(chunks.at(-1)?.flags ?? []),
						thirdPartyEmote.modifierFlag as number,
					];
					latestChunk.flags = flags;
				} else {
					chunks.push({
						type: '3rd_party_emote',
						value: url,
						emoteHeight: thirdPartyEmote.height,
						emoteWidth: thirdPartyEmote.width,
						emoteName: thirdPartyEmote.name,
					});
				}
			} else {
				chunks.push({ type: 'text', value: part });
			}

			currentWordIndex = currentWordIndex + [...part].length + 1;
		}

		return chunks;
	}

	return {
		makeMessageChunks,
	};
});
