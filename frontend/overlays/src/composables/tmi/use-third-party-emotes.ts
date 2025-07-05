import { useIntervalFn } from '@vueuse/core';
import { type Ref, watch } from 'vue';

import { useBetterTv } from './use-bettertv.js';
import type { ChatSettings } from './use-chat-tmi.js';
import { useFrankerFaceZ } from './use-ffz.js';
import { useSevenTv } from './use-seven-tv.js';

const ONE_MINUTE = 60 * 1000;

export function useThirdPartyEmotes(options: Ref<ChatSettings>) {
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

		if (options.emotes.sevenTv) {
			fetchSevenTvEmotes(options.channelId);
		}

		if (options.emotes.bttv) {
			fetchBetterTv();
			bttvResume();
		} else {
			bttvPause();
		}

		if (options.emotes.ffz) {
			fetchFrankerFaceZ();
			ffzResume();
		} else {
			ffzPause();
		}
	});

	function destroy() {
		bttvPause();
		ffzPause();
		destroySevenTv();
	}

	return {
		destroy,
	};
}
