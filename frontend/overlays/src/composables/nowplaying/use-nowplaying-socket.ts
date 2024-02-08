import { useWebSocket } from '@vueuse/core';
import { ref, watch } from 'vue';

import type { TwirWebSocketEvent } from '@/api.js';
import { generateSocketUrlWithParams } from '@/helpers.js';

type Options = {
	apiKey: string,
}

type Track = {
	artist: string,
	title: string,
	image_url?: string,
}

export const useNowPlayingSocket = (options: Options) => {
	const brbUrl = generateSocketUrlWithParams('/overlays/nowplaying', {
		apiKey: options.apiKey,
	});
	const track = ref<Track | null | undefined>();

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
		if (parsedData.eventName === 'nowplaying') {
			track.value = parsedData.data as Track | null;
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
		track,
	};
};
