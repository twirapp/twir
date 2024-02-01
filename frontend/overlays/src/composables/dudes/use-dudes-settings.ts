import { useFontSource } from '@twir/fontsource';
import type { DudesSettings } from '@twirapp/dudes/types';
import { defineStore } from 'pinia';
import { ref } from 'vue';

import type { ChannelData } from '@/types.js';

export const useDudesSettings = defineStore('dudes-settings', () => {
	const fontSource = useFontSource();
	const dudesSettings = ref<DudesSettings | null>(null);
	const channelData = ref<ChannelData>();

	function updateSettings(settings: DudesSettings): void {
		dudesSettings.value = settings;
	}

	function updateChannelData(data: ChannelData): void {
		channelData.value = data;
	}

	async function loadFont(fontFamily: string, fontWeight: number, fontStyle: string): Promise<string> {
		try {
			await fontSource.loadFont(
				fontFamily,
				fontWeight,
				fontStyle,
			);

			const fontKey = `${fontFamily}-${fontWeight}-${fontStyle}`;
			return fontKey;
		} catch (err) {
			console.error(err);
			return 'Arial';
		}
	}

	return {
		channelInfo: channelData,
		updateChannelData,
		dudesSettings,
		updateSettings,
		loadFont,
	};
});
