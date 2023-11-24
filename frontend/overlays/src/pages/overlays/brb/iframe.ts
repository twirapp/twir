import type { OnStart, OnStop, SetSettings } from './types.js';

type Opts = {
	onSettings: SetSettings,
	onStart: OnStart,
	onStop: OnStop,
}

export const useIframe = (opts: Opts) => {
	const onWindowMessage = (msg: MessageEvent<string>) => {
		const parsedData = JSON.parse(msg.data);
		console.log('iframe data: ', parsedData);
		if (parsedData.key === 'settings') {
			opts.onSettings(parsedData.data);
		}
	};

	const create = () => {
		window.postMessage('getSettings');
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
