import { animations } from './animations.js';
import { twirEmote } from './builder.js';
import type { KappagenCallback, SetSettingsCallback, SpawnCallback, KappagenSettings } from './types.js';

type Opts = {
	kappagenCallback: KappagenCallback
	spawnCallback: SpawnCallback
	setSettingsCallback: SetSettingsCallback
	clearCallback?: () => void;
}

export const useIframe = (opts: Opts) => {
	const onWindowMessage = (msg: MessageEvent<string>) => {
		const parsedData = JSON.parse(msg.data) as { key: string, data?: any };

		if (parsedData.key === 'settings' && parsedData.data) {
			const settings = parsedData.data as KappagenSettings;
			opts.setSettingsCallback(settings);
		}

		if (parsedData.key === 'kappa') {
			opts.kappagenCallback([twirEmote], animations[Math.floor(Math.random() * animations.length)]);
		}

		if (parsedData.key === 'kappaWithAnimation') {
			opts.kappagenCallback([twirEmote], parsedData.data.animation);
		}

		if (parsedData.key === 'spawn') {
			opts.spawnCallback([twirEmote]);
		}

		if (parsedData.key === 'clear') {
			opts.clearCallback?.();
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
