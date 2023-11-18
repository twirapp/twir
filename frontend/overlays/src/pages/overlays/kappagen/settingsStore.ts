import { ref } from 'vue';

import type { KappagenSettings } from './types.js';

export const kappagenSettings = ref<KappagenSettings>();

export const useChannelSettings = () => {
	const setSettings = (opts: KappagenSettings) => kappagenSettings.value = opts;

	return {
		kappagenSettings,
		setSettings,
	};
};
