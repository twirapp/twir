import type { DudesSettings } from '@twirapp/dudes/types';
import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useDudesSettings = defineStore('dudes-settings', () => {
	const dudesSettings = ref<DudesSettings>();

	const channelInfo = ref<{ channelId: string, channelName: string }>();

	function updateSettings(settings: DudesSettings) {
		dudesSettings.value = settings;
	}

	return {
		channelInfo,
		dudesSettings,
		updateSettings,
	};
});
