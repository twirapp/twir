import type { Settings } from '@twir/api/messages/overlays_now_playing/overlays_now_playing';
import type { ChannelOverlayNowPlayingPreset } from '@twir/types/api';
import { useWebSocket } from '@vueuse/core';
import { storeToRefs } from 'pinia';
import { watch } from 'vue';

import { useNowPlayingData } from './use-now-playing-data.js';

import type { TwirWebSocketEvent } from '@/api.js';
import { generateSocketUrlWithParams } from '@/helpers.js';

type Options = {
	apiKey: string,
	overlayId: string,
}

export type Track = {
	artist: string,
	title: string,
	image_url?: string,
}

type SettingsWithTypedPreset = Settings & { preset: ChannelOverlayNowPlayingPreset }

export const useNowPlayingSocket = (options: Options) => {
	const { settings, currentTrack } = storeToRefs(useNowPlayingData());
	const brbUrl = generateSocketUrlWithParams('/overlays/nowplaying', {
		apiKey: options.apiKey,
		id: options.overlayId,
	});

	const { data, open, close } = useWebSocket(
		brbUrl,
		{
			immediate: false,
			autoReconnect: {
				delay: 500,
			},
		},
	);

	watch(data, (v) => {
		const parsedData = JSON.parse(v) as TwirWebSocketEvent;
		if (parsedData.eventName === 'settings') {
			settings.value = parsedData.data as SettingsWithTypedPreset;
		}

		if (parsedData.eventName === 'nowplaying') {
			currentTrack.value = parsedData.data as Track | null;
		}
	});

	function connect(): void {
		open();
	}

	function destroy(): void {
		close();
	}

	return {
		connect,
		destroy,
	};
};
