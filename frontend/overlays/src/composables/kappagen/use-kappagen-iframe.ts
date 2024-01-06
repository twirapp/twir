import type { KappagenAnimations } from 'kappagen';

import { twirEmote } from './use-kappagen-builder.js';

import type { KappagenCallback, SetSettingsCallback, SpawnCallback, KappagenSettings } from '@/types.js';

type Options = {
	kappagenCallback: KappagenCallback
	spawnCallback: SpawnCallback
	setSettingsCallback: SetSettingsCallback
	clearCallback?: () => void;
}

export const useKappagenIframe = (options: Options) => {
	const onWindowMessage = (msg: MessageEvent<string>) => {
		const parsedData = JSON.parse(msg.data) as { key: string, data?: any };

		if (parsedData.key === 'settings' && parsedData.data) {
			const settings = parsedData.data as KappagenSettings;
			options.setSettingsCallback(settings);
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
		window.postMessage('getSettings');
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

export const animations: KappagenAnimations[] = [
	{
		style: 'TheCube',
		prefs: {
			size: 0.2,
			center: false,
			faces: false,
			speed: 6,
		},
	},
	{
		style: 'Text',
		prefs: {
			message: ['Twir'],
			time: 3,
		},
	},
	{
		style: 'Confetti',
		count: 150,
	},
	{
		style: 'Spiral',
		count: 150,
	},
	{
		style: 'Stampede',
		count: 150,
	},
	{
		style: 'Burst',
		count: 150,
	},
	{
		style: 'Fountain',
		count: 150,
	},
	{
		style: 'SmallPyramid',
	},
	{
		style: 'Pyramid',
	},
	{
		style: 'Fireworks',
		count: 150,
	},
	{
		style: 'Conga',
		prefs: {
			avoidMiddle: false,
		},
	},
];
