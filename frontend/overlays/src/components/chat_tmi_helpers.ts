import { type MessageChunk, BttvOverlayEmotes } from '@twir/frontend-chat';

import { bttvEmotes, ffzEmotes, sevenTvEmotes } from './chat_tmi_emotes.js';

export function makeMessageChunks(message: string, emotes?: {
	[emoteid: string]: string[];
}): MessageChunk[] {
	const parsedTwitchEmotes = emotes ? Object.entries(emotes).reduce((acc, [id, positions]) => {
		positions.forEach((position) => {
			const [from, to] = position.split('-').map(Number);
			acc.push({ from, to, emoteId: id });
		});
		return acc;
	}, [] as { from: number, to: number, emoteId: string }[]) : [];

	const chunks: MessageChunk[] = [];

	let currentWordIndex = 0;
	for (const part of message.split(' ')) {
		const emote = parsedTwitchEmotes.find(e => e.from === currentWordIndex);
		const thirdPartyEmote = sevenTvEmotes.value[part] || ffzEmotes.value[part] || bttvEmotes.value[part];
		if (emote) {
			chunks.push({ type: 'emote', value: emote.emoteId });
		} else if (thirdPartyEmote) {
			const bttvOverlayEmote = BttvOverlayEmotes.find(e => e === thirdPartyEmote.name);

			if (bttvOverlayEmote) {
				chunks.at(-1)!.zeroWidthModifiers = [...chunks.at(-1)!.zeroWidthModifiers ?? [], thirdPartyEmote.url];
			} else {
				chunks.push({
					type: '3rd_party_emote',
					value: thirdPartyEmote.url,
				});
			}
		} else {
			chunks.push({ type: 'text', value: part });
		}

		currentWordIndex = currentWordIndex + part.length + 1;
	}

	console.log(chunks)

	return chunks;
}
