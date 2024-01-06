import type { OnStart, OnStop, SetSettings } from '@/types.js';

type Opts = {
	onSettings: SetSettings,
	onStart: OnStart,
	onStop: OnStop,
}

export const useBrbIframe = (options: Opts) => {
	const onWindowMessage = (msg: MessageEvent<string>) => {
		const parsedData = JSON.parse(msg.data);
		if (parsedData.key === 'settings') {
			options.onSettings(parsedData.data);
		}

		if (parsedData.key === 'start') {
			options.onStart(parsedData.data.minutes, parsedData.data.text ?? parsedData.data.incomingText);
		}

		if (parsedData.key === 'stop') {
			options.onStop();
		}
	};

	const create = () => {
		window.addEventListener('message', onWindowMessage);
		window.parent.postMessage(JSON.stringify({ key: 'getSettings' }));
	};

	const destroy = () => {
		window.removeEventListener('message', onWindowMessage);
	};

	return {
		create,
		destroy,
	};
};
