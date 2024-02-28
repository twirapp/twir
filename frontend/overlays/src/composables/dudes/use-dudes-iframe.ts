import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import { defineStore } from 'pinia';

import { dudesTwir, getSprite } from './dudes-config.js';
import { useDudesSettings } from './use-dudes-settings.js';
import { useDudesSocket } from './use-dudes-socket.js';
import { useDudes } from './use-dudes.js';

import { randomEmoji } from '@/helpers.js';

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
		if (!dudesStore.dudes) return;

		const parsedData = JSON.parse(msg.data) as DudesPostMessage;

		const dude = dudesStore.dudes.getDude(dudesTwir);
		if (!dude) return;

		if (parsedData.action === 'settings' && parsedData.data) {
			const settings = parsedData.data as Required<Settings>;
			dudesSocketStore.updateSettingFromSocket(settings);
			dude.spriteName = getSprite(settings.dudeSettings.defaultSprite);
			dude.playAnimation('Run', true);
			return;
		}

		if (parsedData.action === 'reset') {
			dudesStore.dudes.removeAllDudes();
			spawnIframeDude();
		}

		if (parsedData.action === 'jump') {
			dude.jump();
		}

		if (parsedData.action === 'grow') {
			dude.grow();
		}

		if (parsedData.action === 'spawn-emote') {
			const emote = dudesStore.getProxiedEmoteUrl({
				type: '3rd_party_emote',
				value: 'https://cdn.7tv.app/emote/60b00d1f0d3a78a196f803e3/1x.gif',
			});
			dude.spitEmotes([emote]);
		}

		if (parsedData.action === 'show-message') {
			dude.addMessage(`Hello, ${dudesSettingsStore.channelData!.channelDisplayName}! ${randomEmoji('emoticons')}`);
		}
	}

	function spawnIframeDude() {
		if (
			!dudesStore.dudes ||
			!dudesSettingsStore.dudesSettings ||
			!dudesSettingsStore.channelData ||
			dudesStore.dudes.getDude(dudesTwir)
		) return;

		const emote = dudesStore.getProxiedEmoteUrl({
			type: '3rd_party_emote',
			value: 'https://cdn.7tv.app/emote/65413498dc0468e8c1fbcdc6/1x.gif',
		});

		const dudeSprite = getSprite(dudesSettingsStore.dudesSettings.overlay.defaultSprite);
		const dude = dudesStore.dudes.createDude(dudesTwir, dudeSprite);
		dude.bodyTint('#8a2be2');
		dude.addMessage(`Hello, ${dudesSettingsStore.channelData.channelDisplayName}! ${randomEmoji('emoticons')}`);
		dude.spitEmotes([emote]);
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
