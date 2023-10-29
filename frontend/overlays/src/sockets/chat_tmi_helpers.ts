import { MessageChunk } from './chat_tmi.js';
import { bttvEmotes, ffzEmotes, sevenTvEmotes } from './chat_tmi_emotes.js';

export function makeMessageChunks(message: string, emotes?: {
	[emoteid: string]: string[];
}): MessageChunk[] {
	const parsedEmotes = emotes ? Object.entries(emotes).reduce((acc, [id, positions]) => {
		positions.forEach((position) => {
			const [from, to] = position.split('-').map(Number);
			acc.push({ from, to, emoteId: id });
		});
		return acc;
	}, [] as { from: number, to: number, emoteId: string }[]) : [];

	const chunks: MessageChunk[] = [];

	let currentWordIndex = 0;
	for (const part of message.split(' ')) {
		const emote = parsedEmotes.find(e => e.from === currentWordIndex);
		if (emote) {
			chunks.push({ type: 'emote', value: emote.emoteId });
		} else if (sevenTvEmotes.value[part]) {
			chunks.push({ type: '7tv_emote', value: sevenTvEmotes.value[part] });
		} else if (ffzEmotes.value[part]) {
			chunks.push({ type: 'ffz_emote', value: ffzEmotes.value[part] });
		} else if (bttvEmotes.value[part]) {
			chunks.push({ type: 'bttv_emote', value: bttvEmotes.value[part] });
		} else {
			chunks.push({ type: 'text', value: part });
		}

		currentWordIndex = currentWordIndex + part.length + 1;
	}

	return chunks;
}

export function normalizeDisplayName(userName: string, displayName: string) {
	if (userName === displayName.toLocaleLowerCase()) {
		return displayName;
	} else {
		return userName;
	}
}
