import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { useWebSocket } from '@vueuse/core';
import { watch } from 'vue';

import type { SetSettings } from './types.js';
import { TwirWebSocketEvent } from '../../../sockets/types.js';

type Opts = {
	apiKey: string,
	onSettings: SetSettings,
	onStart: (minutes: number, text?: string) => void,
	onStop: () => void,
}

export const useBeRightBackOverlaySocket = (opts: Opts) => {
	const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
	const host = window.location.host;

	const { data, send, open, close } = useWebSocket(
		`${protocol}://${host}/socket/overlays/brb?apiKey=${opts.apiKey}`,
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
			opts.onSettings(parsedData.data as Settings);
		}

		if (parsedData.eventName === 'start') {
			opts.onStart(parsedData.data.minutes, parsedData.data.text);
		}

		if (parsedData.eventName === 'stop') {
			opts.onStop();
		}
	});

	const create = () => {
		open();
	};

	const destroy = () => {
		close();
	};

	return {
		create,
		destroy,
	};
};
