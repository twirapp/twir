import { BttvZeroModifiers } from '@twir/frontend-chat';
import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

import type {
	BttvChannelResponse,
	BttvEmote,
	BttvGlobalResponse,
	FfzChannelResponse,
	FfzGlobalResponse,
	SevenTvChannelResponse,
	SevenTvEmote,
	SevenTvGlobalResponse,
} from '@/types.js';

type Emote = {
	urls: string[];
	isZeroWidth?: boolean;
	name: string;
	modifierFlag?: number;
	isModifier?: boolean;
	service: '7tv' | 'bttv' | 'ffz';
	width?: number;
	height?: number;
};

function isZeroWidthEmote(flags: number): boolean {
	return flags === 1 << 0;
}

function bttvEmoteUrls(id: string): string[] {
	return Array.from({ length: 3 }).map(
		(_, index) => `https://cdn.betterttv.net/emote/${id}/${index + 1}x.webp`,
	);
}

export const useEmotes = defineStore('emotes', () => {
	const emotes = ref<Record<string, Emote>>({});

	function setSevenTvEmotes(data: SevenTvChannelResponse | SevenTvGlobalResponse): void {
		let emotesForParse: Array<SevenTvEmote>;
		if ('emote_set' in data) {
			emotesForParse = data.emote_set.emotes;
		} else {
			emotesForParse = data.emotes;
		}

		for (const emote of emotesForParse) {
			updateSevenTvEmote(emote);
		}
	}

	function setBttvEmotes(data: BttvChannelResponse | BttvGlobalResponse): void {
		let emotesForParse: Array<BttvEmote>;

		if ('channelEmotes' in data) {
			emotesForParse = [...data.channelEmotes, ...data.sharedEmotes];
		} else {
			emotesForParse = data;
		}

		for (const emote of emotesForParse) {
			emotes.value[emote.code] = {
				urls: bttvEmoteUrls(emote.id),
				name: emote.code,
				service: 'bttv',
				height: emote.height,
				width: emote.width,
				isModifier: emote.modifier ?? false,
				isZeroWidth: BttvZeroModifiers.some((e) => e === emote.code),
			};
		}
	}

	function setFrankerFaceZEmotes(data: FfzChannelResponse | FfzGlobalResponse): void {
		const sets = Object.values(data.sets);
		for (const set of sets) {
			for (const emote of set.emoticons) {
				emotes.value[emote.name] = {
					urls: Object.values(emote.urls),
					name: emote.name,
					service: 'ffz',
					width: emote.width,
					height: emote.height,
					isModifier: emote.modifier,
					modifierFlag: emote.modifier_flags,
				};
			}
		}
	}

	function removeEmoteByName(emoteName: string): void {
		if (emotes.value[emoteName]) {
			delete emotes.value[emoteName];
		}
	}

	function updateSevenTvEmote(emote: SevenTvEmote): void {
		const files = emote.data.host.files.filter((file) => file.format === 'WEBP');
		const { height, width } = files.at(0)!;

		emotes.value[emote.name] = {
			urls: files.map((file) => `https:${emote.data.host.url}/${file.name}`),
			isZeroWidth: isZeroWidthEmote(emote.flags),
			name: emote.name,
			service: '7tv',
			width,
			height,
		};
	}

	const loadedEmotes = ref<string[]>([]);
	watch(() => emotes, (emotes) => {
		for (const emote of Object.values(emotes)) {
			for (const url of emote.urls) {
				if (loadedEmotes.value.includes(url)) continue;
				const image = new Image();
				image.src = url;
				loadedEmotes.value.push(url);
			}
		}
	});

	return {
		emotes,
		setSevenTvEmotes,
		updateSevenTvEmote,
		removeEmoteByName,
		setBttvEmotes,
		setFrankerFaceZEmotes,
	};
});
