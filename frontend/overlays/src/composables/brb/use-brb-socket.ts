import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { useWebSocket } from '@vueuse/core';
import { watch } from 'vue';

import { TwirWebSocketEvent } from '@/api.js';
import { SetSettings } from '@/types.js';

type Options = {
	brbUrl: string,
	onSettings: SetSettings,
	onStart: (minutes: number, text?: string) => void,
	onStop: () => void,
}

export const useBeRightBackOverlaySocket = (optins: Options) => {
	const { data, send, open, close } = useWebSocket(
		optins.brbUrl,
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

	watch(data, (v) => {
		const parsedData = JSON.parse(v) as TwirWebSocketEvent;

		if (parsedData.eventName === 'settings') {
			optins.onSettings(parsedData.data as Settings);
		}

		if (parsedData.eventName === 'start') {
			optins.onStart(parsedData.data.minutes, parsedData.data.text);
		}

		if (parsedData.eventName === 'stop') {
			optins.onStop();
		}
	});

	function create() {
		open();
	}

	function destroy() {
		close();
	}

	return {
		create,
		destroy,
	};
};
