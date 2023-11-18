import type { MessageChunk } from '@twir/frontend-chat';
import type { Emote } from 'kappagen';
import { computed } from 'vue';

import { emotes } from './chat_tmi_emotes';
import { animations } from './kappagen_animations';

export const useKappagenBuilder = () => {
	const kappagenEmotes = computed(() => {
		const emotesArray = Object.values(emotes.value);

		return emotesArray.filter(e => !e.isZeroWidth && !e.isModifier);
	});

	// ПРОСТО ЧАТ
	const buildSpawnEmotes = (chunks: MessageChunk[]) => {
		const emotesChunks = chunks.filter(c => c.type !== 'text');

		const emotes: Emote[] = [];

		for (const chunk of emotesChunks) {
			if (chunk.type === 'text') continue;

			const zwe = chunk.zeroWidthModifiers?.map(z => ({ url: z })) ?? [];

			if (chunk.type === 'emote') {
				emotes.push({
					url: `https://static-cdn.jtvnw.net/emoticons/v2/${chunk.value}/default/dark/3.0`,
					zwe: chunk.zeroWidthModifiers?.map(z => ({ url: z })) ?? [],
				});
				continue;
			}

			if (chunk.type === '3rd_party_emote') {
				emotes.push({
					url: chunk.value,
					zwe,
					width: chunk.emoteWidth,
					height: chunk.emoteHeight,
				});
				continue;
			}
		}

		return emotes;
	};

	// КОМАНДА И ИВЕНТЫ
	const buildKappagenEmotes = (chunks: MessageChunk[]) => {
		const enabledAnimations = animations.filter(a => a.style !== 'Text');
		const animation = enabledAnimations[Math.floor(Math.random() * enabledAnimations.length)];

		let count = 0;

		if ('count' in animation) {
			count = animation.count;
		} else {
			count = 1;
		}

		const emotes: Emote[] = [];

		const emotesChunks = chunks.filter(c => c.type !== 'text');

		if (!chunks.length) {
			const randomEmotes: Emote[] = Array(count)
			.fill(null)
			.map(() => {
				const randomEmote = kappagenEmotes.value[Math.floor(Math.random() * kappagenEmotes.value.length)];

				return {
					url: randomEmote.urls.at(-1)!,
					width: randomEmote.width,
					height: randomEmote.height,
				};
			});
			emotes.push(...randomEmotes);
		} else {
			for (const chunk of emotesChunks) {
				const emote = buildSpawnEmotes([chunk]);
				if (emote.length) {
					emotes.push(...emote);
				}
			}
		}

		return emotes;
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
