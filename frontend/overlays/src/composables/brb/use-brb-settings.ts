import type { Settings } from '@twir/api/messages/overlays_be_right_back/overlays_be_right_back';
import { defineStore } from 'pinia';
import { ref } from 'vue';

import type { BrbSetSettingsFn } from '@/types.js';

export const useBrbSettings = defineStore('brb-settings', () => {
	const settings = ref<Settings>();

	const setSettings: BrbSetSettingsFn = (newSettings) => {
		settings.value = newSettings;
	};

	return {
		settings,
		setSettings,
	};
});
