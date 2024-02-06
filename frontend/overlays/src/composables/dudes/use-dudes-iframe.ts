import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import { defineStore } from 'pinia';

import { dudesTwir } from './dudes-config.js';
import { useDudesSettings } from './use-dudes-settings.js';
import { useDudes } from './use-dudes.js';

interface DudesPostMessage {
	action: string;
	data?: any
}

export const useDudesIframe = defineStore('dudes-iframe', () => {
	const isIframe = Boolean(window.frameElement);
	const dudesStore = useDudes();
	const dudesSettingsStore = useDudesSettings();

	function onPostMessage(msg: MessageEvent<string>) {
		const parsedData = JSON.parse(msg.data) as DudesPostMessage;
		console.log(parsedData);

		if (parsedData.action === 'settings' && parsedData.data) {
			const settings = parsedData.data as Settings;
			dudesSettingsStore.updateSettings({
				ignore: settings.ignoreSettings!,
				dudes: {
					dude: {
						...settings.dudeSettings!,
						sounds: {
							enabled: settings.dudeSettings!.soundsEnabled,
							volume: settings.dudeSettings!.soundsVolume,
						},
					},
					name: settings.nameBoxSettings!,
					message: settings.messageBoxSettings!,
					spitter: {
						enabled: settings.spitterEmoteSettings!.enabled,
					},
				},
			});

			return;
		}

		const dude = dudesStore.dudes?.getDude(dudesTwir);
		if (parsedData.action === 'jump' && dude) {
			dude.jump();
		} else if (parsedData.action === 'spawn-emote' && dude) {
			dude.spitEmotes([
				'https://cdn.7tv.app/emote/613937fcf7977b64f644c0d2/1x.webp',
				'https://cdn.7tv.app/emote/60b00d1f0d3a78a196f803e3/1x.webp',
			]);
		} else if (parsedData.action === 'show-message' && dude) {
			dude.addMessage(`Hello, ${dudesSettingsStore.channelData!.channelDisplayName}!`);
		}
	}

	function spawnIframeDude() {
		if (dudesStore.dudes?.getDude(dudesTwir)) return;
		dudesStore.createDude(
			dudesTwir,
			'#8a2be2',
			[
				{
					type: 'text',
					value: `Hello, ${dudesSettingsStore.channelData!.channelDisplayName}!`,
				},
			],
		);
	}

	function create() {
		if (!isIframe) return;
		window.addEventListener('message', onPostMessage);
	}

	function destroy() {
		if (!isIframe) return;
		window.removeEventListener('message', onPostMessage);
	}

	return {
		isIframe,
		spawnIframeDude,
		create,
		destroy,
	};
});
