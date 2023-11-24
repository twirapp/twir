import type { OnStart, OnStop, SetSettings } from './types.js';

type Opts = {
	onSettings: SetSettings,
	onStart: OnStart,
	onStop: OnStop,
}

export const useIframe = (opts: Opts) => {
	const onWindowMessage = (msg: MessageEvent<string>) => {
		const parsedData = JSON.parse(msg.data);
		if (parsedData.key === 'settings') {
			opts.onSettings(parsedData.data);
		}

		if (parsedData.key === 'start') {
			opts.onStart(parsedData.data.minutes, parsedData.data.text ?? parsedData.data.incomingText);
		}

		if (parsedData.key === 'stop') {
			opts.onStop();
		}
	};

	const create = () => {
		window.parent.postMessage(JSON.stringify({ key: 'getSettings' }));
		window.addEventListener('message', onWindowMessage);
	};

	const destroy = () => {
		window.removeEventListener('message', onWindowMessage);
	};

	return {
		create,
		destroy,
	};
};
