import type { KappagenMethods } from '@twirapp/kappagen/types';

import { animations } from './kappagen-animations.js';
import { twirEmote } from './use-kappagen-builder.js';
import { useKappagenSettings } from './use-kappagen-settings.js';

import type { KappagenSettings } from '@/types.js';

type Options = KappagenMethods

export const useKappagenIframe = (options: Options) => {
	const { updateSettings } = useKappagenSettings();

	const onWindowMessage = (msg: MessageEvent<string>) => {
		const parsedData = JSON.parse(msg.data) as { key: string, data?: any };

		if (parsedData.key === 'settings' && parsedData.data) {
			const settings = parsedData.data as KappagenSettings;
			updateSettings(settings);
		}

		if (parsedData.key === 'kappa') {
			options.playAnimation([twirEmote], animations[Math.floor(Math.random() * animations.length)]);
		}

		if (parsedData.key === 'kappaWithAnimation') {
			options.playAnimation([twirEmote], parsedData.data.animation);
		}

		if (parsedData.key === 'spawn') {
			options.showEmotes([twirEmote]);
		}

		if (parsedData.key === 'clear') {
			options.clear();
		}
	};

	function create() {
		window.addEventListener('message', onWindowMessage);
	}

	function destroy() {
		window.removeEventListener('message', onWindowMessage);
	}

	return {
		create,
		destroy,
	};
};
