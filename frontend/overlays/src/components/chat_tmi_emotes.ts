/* eslint-disable no-empty */
import { useFetch, useIntervalFn } from '@vueuse/core';
import { ref, Ref, computed, onMounted, onUnmounted, watch } from 'vue';


type Emote = {
	url: string,
	isZeroWidth?: boolean,
	name: string,
}

export const sevenTvEmotes = ref<Record<string, Emote>>({});
export const bttvEmotes = ref<Record<string, Emote>>({});
export const ffzEmotes = ref<Record<string, Emote>>({});

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
				ffzEmotes.value[emote.name] = {
					url: Object.values(emote.urls).at(-1)!,
					name: emote.name,
				};
			}
		}
	});
	watch(ffzGlobalEmotes.data, (v) => {
		if (!v) return;

		const sets = Object.values(v!.sets);
		for (const set of sets) {
			for (const emote of set.emoticons) {
				ffzEmotes.value[emote.name] = {
					url: Object.values(emote.urls).at(0)!,
					name: emote.name,
				};
			}
		}
	});

	watch(sevenTvChannelEmotes.data, (v) => {
		if (!v) return;

		for (const emote of v.emote_set.emotes) {
			const file = emote.data.host.files.filter(f => f.format === 'WEBP');
			sevenTvEmotes.value[emote.name] = {
				url: `https:${emote.data.host.url}/${file.at(0)!.name}`,
				isZeroWidth: isZeroWidthEmote(emote.flags),
				name: emote.name,
			};
		}
	});
	watch(sevenTvGlobalEmotes.data, (v) => {
		if (!v) return;

		for (const emote of v.emotes) {
			const file = emote.data.host.files.filter(f => f.format === 'WEBP');
			sevenTvEmotes.value[emote.name] = {
				url: `https:${emote.data.host.url}/${file.at(0)!.name}`,
				isZeroWidth: isZeroWidthEmote(emote.flags),
				name: emote.name,
			};
		}
	});

	watch(bttvChannelEmotes.data, (v) => {
		if (!v) return;

		for (const emote of [...v.sharedEmotes, ...v.channelEmotes]) {
			bttvEmotes.value[emote.code] = {
				url: `https://cdn.betterttv.net/emote/${emote.id}/1x.webp`,
				name: emote.code,
			};
		}
	});

	watch(bttvGlobalEmotes.data, (v) => {
		if (!v) return;

		for (const emote of v) {
			bttvEmotes.value[emote.code] = {
				url: `https://cdn.betterttv.net/emote/${emote.id}/1x.webp`,
				name: emote.code,
			};
		}
	});

	return {
		sevenTvEmotes,
		bttvEmotes,
		ffzEmotes,
	};
};

type SevenTvEmote = {
	name: string,
	data: {
		host: { url: string, files: Array<{ name: string, format: string }> }
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

type BttvEmote = { code: string, imageType: string, id: string }
type BttvChannelResponse = {
	channelEmotes: Array<BttvEmote>
	sharedEmotes: Array<BttvEmote>
}
type BttvGlobalResponse = Array<BttvEmote>

type FfzChannelResponse = {
	sets: {
		[x: string]: {
			emoticons: Array<{ name: string, urls: Record<string, string>}>
		}
	}
}

type FfzGlobalResponse = {
	sets: {
		[x: string]: {
			emoticons: Array<{ name: string, urls: Record<string, string>}>
		}
	}
}
