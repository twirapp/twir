import { useWebSocket } from '@vueuse/core';

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


	const connect = open;

	return {
		connect,
		data,
	};
};
