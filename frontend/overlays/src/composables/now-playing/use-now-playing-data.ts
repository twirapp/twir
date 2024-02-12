import type { Settings } from '@twir/api/messages/overlays_now_playing/overlays_now_playing';
import type { ChannelOverlayNowPlayingPreset } from '@twir/types/api';
import { defineStore } from 'pinia';
import { ref } from 'vue';

export type Track = {
	artist: string,
	title: string,
	image_url?: string,
}

type SettingsWithTypedPreset = Settings & { preset: ChannelOverlayNowPlayingPreset }

export const useNowPlayingData = defineStore('now-playing-data', () => {
	const currentTrack = ref<Track | null | undefined>();
	const settings = ref<SettingsWithTypedPreset>();

	return {
		currentTrack,
		settings,
	};
});
