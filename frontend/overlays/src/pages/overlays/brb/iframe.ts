import type { SetSettings } from './types.js';

export const useIframe = (setSettings: SetSettings) => {
	const onWindowMessage = (msg: MessageEvent<string>) => {
		const parsedData = JSON.parse(msg.data);
		console.log('iframe data: ', parsedData);
		if (parsedData.key === 'settings') {
			setSettings(parsedData.data);
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
