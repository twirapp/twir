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
		console.log(msg);

		const parsedData = JSON.parse(msg.data) as DudesPostMessage;

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
		}

		if (parsedData.action === 'jump') {
			dudesStore.dudes?.getDude(dudesTwir)?.jump();
		}

		if (parsedData.action === 'spawn-emote') {
			dudesStore.dudes?.getDude(dudesTwir)?.spitEmotes([]);
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
				{
					type: '3rd_party_emote',
					value: 'https://cdn.7tv.app/emote/63706216d49eb6644629aa52/1x.gif',
				},
			],
		);
	}

	function create() {
		if (!isIframe) return;
		window.addEventListener('message', onPostMessage);
		window.parent.postMessage(JSON.stringify({ action: 'get-settings' }));
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
