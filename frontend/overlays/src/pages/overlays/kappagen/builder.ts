import type { MessageChunk } from '@twir/frontend-chat';
import type { Emote } from 'kappagen';
import { computed } from 'vue';

import { emotes } from '../../../components/chat_tmi_emotes';

export const useKappagenBuilder = () => {
	const kappagenEmotes = computed(() => {
		const emotesArray = Object.values(emotes.value);

		return emotesArray.filter(e => !e.isZeroWidth && !e.isModifier);
	});

	// ПРОСТО ЧАТ
	const buildSpawnEmotes = (chunks: MessageChunk[]) => {
		const emotesChunks = chunks.filter(c => c.type !== 'text');

		const result: Emote[] = [];

		for (const chunk of emotesChunks) {
			if (chunk.type === 'text') continue;

			const zwe = chunk.zeroWidthModifiers?.map(z => ({ url: z })) ?? [];

			if (chunk.type === 'emote') {
				result.push({
					url: `https://static-cdn.jtvnw.net/emoticons/v2/${chunk.value}/default/dark/3.0`,
					zwe: chunk.zeroWidthModifiers?.map(z => ({ url: z })) ?? [],
				});
				continue;
			}

			if (chunk.type === '3rd_party_emote') {
				result.push({
					url: chunk.value,
					zwe,
					width: chunk.emoteWidth,
					height: chunk.emoteHeight,
				});
				continue;
			}
		}

		return result;
	};

	// КОМАНДА И ИВЕНТЫ
	const buildKappagenEmotes = (chunks: MessageChunk[]) => {
		const result: Emote[] = [];

		const emotesChunks = chunks.filter(c => c.type !== 'text');

		if (!emotesChunks.length) {
			const mappedEmotes = kappagenEmotes.value.map(v => ({
				url: v.urls.at(-1)!,
				width: v.width,
				height: v.height,
			}));

			result.push(...mappedEmotes);
		} else {
			for (const chunk of emotesChunks) {
				const emote = buildSpawnEmotes([chunk]);
				if (emote.length) {
					result.push(...emote);
				}
			}
		}

		return result;
	};

	return {
		buildKappagenEmotes,
		buildSpawnEmotes,
	};
};

export const twirEmote: Emote = {
	url: 'https://cdn.7tv.app/emote/6548b7074789656a7be787e1/4x.webp',
	zwe: [
		{
			url: 'https://cdn.7tv.app/emote/6128ed55a50c52b1429e09dc/4x.webp',
		},
	],
};
