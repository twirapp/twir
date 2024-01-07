import { useBrbSettings } from './use-brb-settings';

import type { BrbOnStartFn, BrbOnStopFn } from '@/types.js';

type Options = {
	onStart: BrbOnStartFn,
	onStop: BrbOnStopFn,
}

export const useBrbIframe = (options: Options) => {
	const { setSettings } = useBrbSettings();

	const onWindowMessage = (msg: MessageEvent<string>) => {
		const parsedData = JSON.parse(msg.data);
		if (parsedData.key === 'settings') {
			setSettings(parsedData.data);
		}

		if (parsedData.key === 'start') {
			options.onStart(parsedData.data.minutes, parsedData.data.text ?? parsedData.data.incomingText);
		}

		if (parsedData.key === 'stop') {
			options.onStop();
		}
	};

	function create(): void {
		window.addEventListener('message', onWindowMessage);
		window.parent.postMessage(JSON.stringify({ key: 'getSettings' }));
	}

	function destroy(): void {
		window.removeEventListener('message', onWindowMessage);
	}

	return {
		create,
		destroy,
	};
};
