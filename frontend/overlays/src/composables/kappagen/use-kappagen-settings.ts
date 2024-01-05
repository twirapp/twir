import { defineStore } from 'pinia';
import { ref } from 'vue';

import type { KappagenSettings } from '@/types.js';

export const useKappagenSettings = defineStore('kappagen-settings', () => {
	const settings = ref<KappagenSettings>();

	function setSettings(newOptions: KappagenSettings): void {
		settings.value = newOptions;
	}

	return {
		settings,
		setSettings,
	};
});
