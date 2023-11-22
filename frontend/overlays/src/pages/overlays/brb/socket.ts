import { useWebSocket } from '@vueuse/core';
import { watch } from 'vue';

import type { SetSettings } from './types.js';

export const useBeRightBackOverlaySocket = (apiKey: string, setSettings: SetSettings) => {
	const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
	const host = window.location.host;

	const { data, send, open, close } = useWebSocket(
		`${protocol}://${host}/socket/overlays/brb?apiKey=${apiKey}`,
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
		console.log(v);
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
