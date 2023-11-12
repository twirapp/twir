/* eslint-disable no-empty */
import { BttvZeroModifiers } from '@twir/frontend-chat';
import { useFetch, useIntervalFn } from '@vueuse/core';
import { ref, Ref, computed, onMounted, onUnmounted, watch } from 'vue';

type Emote = {
	urls: string[],
	isZeroWidth?: boolean,
	name: string,
	modifierFlag?: number
	isModifier?: boolean
	service: '7tv' | 'bttv' | 'ffz'
	width?: number
	height?: number
}

export const emotes = ref<Record<string, Emote>>({});

const isZeroWidthEmote = (flags: number) => {
	return flags === (1 << 0);
};

export const useThirdPartyEmotes = (channelName: Ref<string>, channelId: Ref<string>) => {

	const seventvUrl = computed(() => `https://7tv.io/v3/users/twitch/${channelId.value}`);

	const sevenTvChannelEmotes = useFetch(seventvUrl, {
		refetch: true,
		beforeFetch({ cancel }) {
			if (!channelId.value) cancel();
		},
	}).get().json<SevenTvChannelResponse>();
	const sevenTvGlobalEmotes = useFetch('https://7tv.io/v3/emote-sets/global').get().json<SevenTvGlobalResponse>();
	const sevenTvEmotesInterval = useIntervalFn(() => {
		try {
			sevenTvChannelEmotes.execute(false);
			sevenTvGlobalEmotes.execute(false);
		} catch {}
	}, 10 * 1000, {
		immediate: false,
	});

	const ffzUrl = computed(() => `https://api.frankerfacez.com/v1/room/${channelName.value}`);
	const ffzChannelEmotes = useFetch(ffzUrl, {
		refetch: true,
		beforeFetch({ cancel }) {
			if (!channelName.value) cancel();
		},
	}).get().json<FfzChannelResponse>();
	const ffzGlobalEmotes = useFetch('https://api.frankerfacez.com/v1/set/global').get().json<FfzGlobalResponse>();
	const ffzEmotesInterval = useIntervalFn(() => {
		try {
			ffzChannelEmotes.execute(false);
			ffzGlobalEmotes.execute(false);
		} catch {}
	}, 10 * 1000, {
		immediate: false,
	});

	const bttvUrl = computed(() => `https://api.betterttv.net/3/cached/users/twitch/${channelId.value}`);
	const bttvChannelEmotes = useFetch(bttvUrl, {
		refetch: true,
		beforeFetch({ cancel }) {
			if (!channelId.value) cancel();
		},
	}).get().json<BttvChannelResponse>();
	const bttvGlobalEmotes = useFetch('https://api.betterttv.net/3/cached/emotes/global').get().json<BttvGlobalResponse>();
	const bttvEmotesInterval = useIntervalFn(() => {
		try {
			bttvChannelEmotes.execute(false);
			bttvGlobalEmotes.execute(false);
		} catch {}
	}, 10 * 1000, {
		immediate: false,
	});

	onMounted(() => {
		sevenTvEmotesInterval.resume();
		ffzEmotesInterval.resume();
		bttvEmotesInterval.resume();
	});

	onUnmounted(() => {
		sevenTvEmotesInterval.pause();
		ffzEmotesInterval.pause();
		bttvEmotesInterval.pause();
	});

	watch(ffzChannelEmotes.data, (v) => {
		if (!v) return;

		const sets = Object.values(v!.sets);
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
	});
	watch(ffzGlobalEmotes.data, (v) => {
		if (!v) return;

		const sets = Object.values(v!.sets);
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
	});

	watch(sevenTvChannelEmotes.data, (v) => {
		if (!v) return;

		for (const emote of v.emote_set.emotes) {
			const files = emote.data.host.files.filter(f => f.format === 'WEBP');

			emotes.value[emote.name] = {
				urls: files.map(f => `https:${emote.data.host.url}/${f.name}`),
				isZeroWidth: isZeroWidthEmote(emote.flags),
				name: emote.name,
				service: '7tv',
				width: files.at(0)!.width,
				height: files.at(0)!.height,
			};
		}
	});
	watch(sevenTvGlobalEmotes.data, (v) => {
		if (!v) return;

		for (const emote of v.emotes) {
			const files = emote.data.host.files.filter(f => f.format === 'WEBP');
			emotes.value[emote.name] = {
				urls: files.map(f => `https:${emote.data.host.url}/${f.name}`),
				isZeroWidth: isZeroWidthEmote(emote.flags),
				name: emote.name,
				service: '7tv',
				width: files.at(0)!.width,
				height: files.at(0)!.height,
			};
		}
	});

	const genBttvUrls = (id: string) => {
		return Array.from({ length: 3 }).map((_, index) => `https://cdn.betterttv.net/emote/${id}/${index+1}x.webp`);
	};

	watch(bttvChannelEmotes.data, (v) => {
		if (!v) return;

		for (const emote of [...v.sharedEmotes, ...v.channelEmotes]) {
			emotes.value[emote.code] = {
				urls: genBttvUrls(emote.id),
				name: emote.code,
				service: 'bttv',
				height: emote.height,
				width: emote.width,
				isModifier: emote.modifier ?? false,
				isZeroWidth: BttvZeroModifiers.some(e => e === emote.code),
			};
		}
	});

	watch(bttvGlobalEmotes.data, (v) => {
		if (!v) return;

		for (const emote of v) {
			emotes.value[emote.code] = {
				urls: genBttvUrls(emote.id),
				name: emote.code,
				service: 'bttv',
				height: emote.height,
				width: emote.width,
				isModifier: emote.modifier ?? false,
				isZeroWidth: BttvZeroModifiers.some(e => e === emote.code),
			};
		}
	});

	return {
		emotes,
	};
};

type SevenTvEmote = {
	name: string,
	data: {
		host: { url: string, files: Array<{ name: string, format: string, height: number, width: number }> }
	},
	flags: number,
}
type SevenTvChannelResponse = {
	emote_set: {
		emotes: Array<SevenTvEmote>
	}
}
type SevenTvGlobalResponse = {
	emotes: Array<SevenTvEmote>
}

type BttvEmote = {
	code: string,
	imageType: string,
	id: string,
	height?: number,
	width?: number,
	modifier?: boolean,
}
type BttvChannelResponse = {
	channelEmotes: Array<BttvEmote>
	sharedEmotes: Array<BttvEmote>
}
type BttvGlobalResponse = Array<BttvEmote>

type FfzEmote = {
	name: string,
	urls: Record<string, string>,
	height: number,
	width: number,
	modifier: boolean;
	modifier_flags?: number
}

type FfzChannelResponse = {
	sets: {
		[x: string]: {
			emoticons: FfzEmote[]
		}
	}
}

type FfzGlobalResponse = {
	sets: {
		[x: string]: {
			emoticons: FfzEmote[]
		}
	}
}
