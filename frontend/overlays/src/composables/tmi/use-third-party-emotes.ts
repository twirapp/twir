import { useIntervalFn } from '@vueuse/core';
import { watch, type Ref, onUnmounted } from 'vue';

import { useBetterTv } from './use-bettertv.js';
import { useFrankerFaceZ } from './use-ffz.js';
import { useSevenTv } from './use-seven-tv.js';

export type ThirdPartyEmotesOptions = {
	channelName?: string;
	channelId?: string;
	sevenTv?: boolean;
	bttv?: boolean;
	ffz?: boolean;
};

const ONE_MINUTE = 60 * 1000;

export function useThirdPartyEmotes(options: Ref<ThirdPartyEmotesOptions>) {
	const { fetchSevenTvEmotes, destroy: destroySevenTv } = useSevenTv();
	const { fetchBttvEmotes } = useBetterTv();
	const { fetchFrankerFaceZEmotes } = useFrankerFaceZ();

	function fetchBetterTv() {
		if (!options.value.channelId) return;
		fetchBttvEmotes(options.value.channelId);
	}

	function fetchFrankerFaceZ() {
		if (!options.value.channelId) return;
		fetchFrankerFaceZEmotes(options.value.channelId);
	}

	const { pause: bttvPause, resume: bttvResume } = useIntervalFn(fetchBetterTv, ONE_MINUTE * 5);
	const { pause: ffzPause, resume: ffzResume } = useIntervalFn(fetchFrankerFaceZ, ONE_MINUTE * 10);

	watch(() => options.value, async (options) => {
		if (!options.channelId) return;

		if (options.sevenTv) {
			fetchSevenTvEmotes(options.channelId);
		}

		if (options.bttv) {
			fetchBetterTv();
			bttvResume();
		} else {
			bttvPause();
		}

		if (options.ffz) {
			fetchFrankerFaceZ();
			ffzResume();
		} else {
			ffzPause();
		}
	});

	onUnmounted(() => {
		bttvPause();
		ffzPause();
		destroySevenTv();
	});
}
