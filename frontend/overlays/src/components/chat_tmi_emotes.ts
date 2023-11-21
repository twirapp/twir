/* eslint-disable no-empty */
import { BttvZeroModifiers } from '@twir/frontend-chat';
import { useIntervalFn } from '@vueuse/core';
import { ref, Ref, computed, onUnmounted, watch } from 'vue';

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

export type Opts = {
	channelName?: string,
	channelId?: string,
	sevenTv?: boolean,
	bttv?: boolean,
	ffz?: boolean,
}

const setFfzEmotes = (data: FfzChannelResponse | FfzGlobalResponse) => {
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
};

const setSevenTvEmotes = (data: SevenTvChannelResponse | SevenTvGlobalResponse) => {
	let emotesForParse: Array<SevenTvEmote>;
	if ('emote_set' in data) {
		emotesForParse = data.emote_set.emotes;
	} else {
		emotesForParse = data.emotes;
	}

	for (const emote of emotesForParse) {
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
};

const genBttvUrls = (id: string) => {
	return Array.from({ length: 3 }).map((_, index) => `https://cdn.betterttv.net/emote/${id}/${index+1}x.webp`);
};

const setBttvEmotes = (data: BttvChannelResponse | BttvGlobalResponse) => {
	let emotesForParse: Array<BttvEmote>;

	if ('channelEmotes' in data) {
		emotesForParse = [...data.channelEmotes, ...data.sharedEmotes];
	} else {
		emotesForParse = data;
	}

	for (const emote of emotesForParse) {
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
};

export const useThirdPartyEmotes = (opts: Ref<Opts>) => {
	const seventvUrl = computed(() => `https://7tv.io/v3/users/twitch/${opts.value.channelId}`);
	const ffzUrl = computed(() => `https://api.frankerfacez.com/v1/room/id/${opts.value.channelId}`);
	const bttvUrl = computed(() => `https://api.betterttv.net/3/cached/users/twitch/${opts.value.channelId}`);

	const fetchFfz = async () => {
		if (!opts.value.ffz) return;

		try {
			const [global, channel] = await Promise.all([
				fetch('https://api.frankerfacez.com/v1/set/global'),
				fetch(ffzUrl.value),
			]);

			setFfzEmotes(await global.json() as FfzGlobalResponse);
			setFfzEmotes(await channel.json() as FfzChannelResponse);
		} catch {}
	};

	const ffzEmotesInterval = useIntervalFn(fetchFfz, 10 * 1000, {
		immediate: false,
	});

	const fetchSevenTv = async () => {
		if (!opts.value.sevenTv) return;

		try {
			const [global, channel] = await Promise.all([
				fetch('https://7tv.io/v3/emote-sets/global'),
				fetch(seventvUrl.value),
			]);

			setSevenTvEmotes(await global.json() as SevenTvChannelResponse);
			setSevenTvEmotes(await channel.json() as SevenTvGlobalResponse);
		} catch {}
	};

	const sevenTvEmotesInterval = useIntervalFn(fetchSevenTv, 10 * 1000, {
		immediate: false,
	});


	const fetchBttv = async () => {
		if (!opts.value.bttv) return;

		try {
			const [global, channel] = await Promise.all([
				fetch('https://api.betterttv.net/3/cached/emotes/global'),
				fetch(bttvUrl.value),
			]);

			setBttvEmotes(await global.json() as BttvGlobalResponse);
			setBttvEmotes(await channel.json() as BttvChannelResponse);
		} catch {}
	};
	const bttvEmotesInterval = useIntervalFn(fetchBttv, 10 * 1000, {
		immediate: false,
	});

	const stop = () => {
		sevenTvEmotesInterval.pause();
		ffzEmotesInterval.pause();
		bttvEmotesInterval.pause();
	};

	const start = () => {
		fetchBttv();
		fetchFfz();
		fetchSevenTv();

		sevenTvEmotesInterval.resume();
		ffzEmotesInterval.resume();
		bttvEmotesInterval.resume();
	};

	watch(opts, (v) => {
		if (!v.channelId) return;

		stop();
		start();
	});

	onUnmounted(() => {
		stop();
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
