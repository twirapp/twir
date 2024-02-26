import { EmojiStyle } from '@twir/api/messages/overlays_kappagen/overlays_kappagen';
import type { MessageChunk } from '@twir/frontend-chat';
import type { Emote } from '@twirapp/kappagen/types';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';

import { useKappagenSettings } from '@/composables/kappagen/use-kappagen-settings.js';
import { useEmotes } from '@/composables/tmi/use-emotes.js';

const getEmojiStyleName = (style: EmojiStyle) => {
	switch (style) {
		case EmojiStyle.Blobmoji:
			return 'blob';
		case EmojiStyle.Noto:
			return 'noto';
		case EmojiStyle.Openmoji:
			return 'openmoji';
		case EmojiStyle.Twemoji:
			return 'twemoji';
	}
};

export type Buidler = {
	buildKappagenEmotes: (chunks: MessageChunk[]) => Emote[];
	buildSpawnEmotes: (chunks: MessageChunk[]) => Emote[];
}

export const useKappagenEmotesBuilder = (): Buidler => {
	const emotesStore = useEmotes();
	const { emotes } = storeToRefs(emotesStore);

	const kappagenSettingsStore = useKappagenSettings();
	const { overlaySettings } = storeToRefs(kappagenSettingsStore);

	const kappagenEmotes = computed(() => {
		if (!emotes.value) return [];
		const emotesArray = Object.values(emotes.value);
		return emotesArray.filter(e => !e.isZeroWidth && !e.isModifier);
	});

	// chat events
	const buildSpawnEmotes = (chunks: MessageChunk[]) => {
		const emotes: Emote[] = [];

		for (const chunk of chunks) {
			if (chunk.type === 'text') continue;

			const zwe = chunk.zeroWidthModifiers?.map(z => ({ url: z })) ?? [];

			if (chunk.emoteName && overlaySettings.value?.excludedEmotes?.includes(chunk.emoteName)) continue;

			if (chunk.type === 'emote') {
				emotes.push({
					url: chunk.value,
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

			const emojiStyle = overlaySettings.value?.emotes?.emojiStyle;
			if (chunk.type === 'emoji' && emojiStyle) {
				const code = chunk.value.codePointAt(0)?.toString(16);
				if (!code) continue;
				emotes.push({
					url: `https://cdn.frankerfacez.com/static/emoji/images/${getEmojiStyleName(emojiStyle)}/${code}.png`,
				});
			}
		}

		return emotes;
	};

	// command, twitch events
	const buildKappagenEmotes = (chunks: MessageChunk[]) => {
		const result: Emote[] = [];

		const emotesChunks = chunks.filter(c => c.type !== 'text');
		if (!emotesChunks.length) {
			const mappedEmotes = kappagenEmotes.value
				.filter(v => !overlaySettings.value?.excludedEmotes?.includes(v.name))
				.map(v => ({
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
