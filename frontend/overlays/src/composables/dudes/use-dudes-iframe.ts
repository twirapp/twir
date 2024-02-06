import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import { defineStore } from 'pinia';

import { dudesTwir } from './dudes-config.js';
import { useDudesSettings } from './use-dudes-settings.js';
import { useDudesSocket } from './use-dudes-socket.js';
import { useDudes } from './use-dudes.js';

interface DudesPostMessage {
	action: string;
	data?: any
}

export const useDudesIframe = defineStore('dudes-iframe', () => {
	const isIframe = Boolean(window.frameElement);
	const dudesStore = useDudes();
	const dudesSocketStore = useDudesSocket();
	const dudesSettingsStore = useDudesSettings();

	async function onPostMessage(msg: MessageEvent<string>) {
		const parsedData = JSON.parse(msg.data) as DudesPostMessage;

		if (parsedData.action === 'settings' && parsedData.data) {
			const settings = parsedData.data as Required<Settings>;
			dudesSocketStore.updateSettingFromSocket(settings);
			return;
		}

		const dude = dudesStore.dudes?.getDude(dudesTwir);
		if (parsedData.action === 'jump' && dude) {
			dude.jump();
		} else if (parsedData.action === 'spawn-emote' && dude) {
			const emote = dudesStore.getProxiedEmoteUrl({
				type: '3rd_party_emote',
				value: 'https://cdn.7tv.app/emote/60b00d1f0d3a78a196f803e3/1x.gif',
			});
			dude.spitEmotes([emote]);
		} else if (parsedData.action === 'show-message' && dude) {
			dude.addMessage(`Hello, ${dudesSettingsStore.channelData!.channelDisplayName}!`);
		}
	}

	function spawnIframeDude() {
		if (dudesStore.dudes?.getDude(dudesTwir)) return;

		const emote = dudesStore.getProxiedEmoteUrl({
			type: '3rd_party_emote',
			value: 'https://cdn.7tv.app/emote/65413498dc0468e8c1fbcdc6/1x.gif',
		});

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
					value: emote,
				},
			],
		);
	}

	function connect() {
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
		connect,
		destroy,
	};
});
