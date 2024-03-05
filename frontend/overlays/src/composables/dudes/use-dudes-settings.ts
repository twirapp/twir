import type { IgnoreSettings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import { useFontSource } from '@twir/fontsource';
import type { DudesSprite, DudesUserSettings } from '@twir/types/overlays';
import type { DudesTypes } from '@twirapp/dudes/types';
import { defineStore } from 'pinia';
import { ref } from 'vue';

import type { ChannelData } from '@/types.js';

export type DudesOverlaySettings = {
	maxOnScreen: number
	defaultSprite: keyof typeof DudesSprite
}

export type DudesConfig = {
	ignore: IgnoreSettings
	dudes: {
		dude: DudesTypes.DudeStyles
		name: DudesTypes.NameBoxStyles
		message: DudesTypes.MessageBoxStyles
		emote: DudesTypes.EmotesStyles
	}
	overlay: DudesOverlaySettings
}

export const useDudesSettings = defineStore('dudes-settings', () => {
	const fontSource = useFontSource();
	const dudesSettings = ref<DudesConfig | null>(null);
	const dudesUserSettings = new Map<string, DudesUserSettings>();
	const channelData = ref<ChannelData>();

	function updateSettings(settings: DudesConfig): void {
		dudesSettings.value = settings;
	}

	function updateChannelData(data: ChannelData): void {
		channelData.value = data;
	}

	async function loadFont(
		fontFamily: string,
		fontWeight: number,
		fontStyle: string,
	): Promise<string> {
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
		channelData,
		updateChannelData,
		dudesSettings,
		dudesUserSettings,
		updateSettings,
		loadFont,
	};
});
