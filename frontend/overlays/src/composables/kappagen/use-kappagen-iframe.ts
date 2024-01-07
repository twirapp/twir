import { animations } from './kappagen-animations.js';
import { twirEmote } from './use-kappagen-builder.js';
import { useKappagenSettings } from './use-kappagen-settings.js';

import type { KappagenSpawnAnimatedEmotesFn, KappagenSpawnEmotesFn, KappagenSettings } from '@/types.js';

type Options = {
	kappagenCallback: KappagenSpawnAnimatedEmotesFn
	spawnCallback: KappagenSpawnEmotesFn
	clearCallback?: () => void;
}

export const useKappagenIframe = (options: Options) => {
	const { setSettings } = useKappagenSettings();

	const onWindowMessage = (msg: MessageEvent<string>) => {
		const parsedData = JSON.parse(msg.data) as { key: string, data?: any };

		if (parsedData.key === 'settings' && parsedData.data) {
			const settings = parsedData.data as KappagenSettings;
			setSettings(settings);
		}

		if (parsedData.key === 'kappa') {
			options.kappagenCallback([twirEmote], animations[Math.floor(Math.random() * animations.length)]);
		}

		if (parsedData.key === 'kappaWithAnimation') {
			options.kappagenCallback([twirEmote], parsedData.data.animation);
		}

		if (parsedData.key === 'spawn') {
			options.spawnCallback([twirEmote]);
		}

		if (parsedData.key === 'clear') {
			options.clearCallback?.();
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
