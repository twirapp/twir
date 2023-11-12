import { useWebSocket } from '@vueuse/core';
import { watch } from 'vue';

import type { TwirWebSocketEvent } from './types.js';

export const useKappagenOverlaySocket = (apiKey: string) => {
	const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
	const host = window.location.host;


	const { data, send, open } = useWebSocket(
		`${protocol}://${host}/socket/overlays/kappagen?apiKey=${apiKey}`,
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

	watch(data, (d: string) => {
		const event = JSON.parse(d) as TwirWebSocketEvent;

		console.log(event);
	});

	const connect = open;

	return {
		connect,
	};
};
