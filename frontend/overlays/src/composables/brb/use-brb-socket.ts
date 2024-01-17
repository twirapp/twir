import type { Settings } from '@twir/api/messages/overlays_be_right_back/overlays_be_right_back';
import { useWebSocket } from '@vueuse/core';
import { watch } from 'vue';

import { useBrbSettings } from './use-brb-settings.js';

import type { TwirWebSocketEvent } from '@/api.js';
import { generateSocketUrlWithParams } from '@/helpers.js';

type Options = {
	apiKey: string,
	onStart: (minutes: number, text: string) => void,
	onStop: () => void,
}

export const useBeRightBackOverlaySocket = (options: Options) => {
	const brbUrl = generateSocketUrlWithParams('/overlays/brb', {
		apiKey: options.apiKey,
	});

	const { data, send, open, close } = useWebSocket(
		brbUrl,
		{
			immediate: false,
			autoReconnect: {
				delay: 500,
			},
			onConnected() {
				send(JSON.stringify({ eventName: 'getSettings' }));
			},
		},
	);

	const { setSettings } = useBrbSettings();

	watch(data, (v) => {
		const parsedData = JSON.parse(v) as TwirWebSocketEvent;

		if (parsedData.eventName === 'settings') {
			setSettings(parsedData.data as Settings);
		}

		if (parsedData.eventName === 'start') {
			options.onStart(parsedData.data.minutes, parsedData.data.text ?? '');
		}

		if (parsedData.eventName === 'stop') {
			options.onStop();
		}
	});

	function create(): void {
		open();
	}

	function destroy(): void {
		close();
	}

	return {
		create,
		destroy,
	};
};
