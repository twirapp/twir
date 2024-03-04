import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import { DudesLayers } from '@twirapp/dudes';
import { defineStore } from 'pinia';

import { dudesTwir, getSprite, type DudeSprite } from './dudes-config.js';
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

		if (parsedData.action === 'settings' && parsedData.data && dude) {
			const settings = parsedData.data as Required<Settings>;
			dudesSocketStore.updateSettingFromSocket(settings);
			const spriteData = getSprite(dudesTwir, settings.dudeSettings.defaultSprite as DudeSprite);
			await dude.updateSpriteData(spriteData);
			return;
		}

		if (parsedData.action === 'reset') {
			dudesStore.dudes.removeAllDudes();
			spawnIframeDude();
		}

		if (parsedData.action === 'jump' && dude) {
			dude.jump();
		}

		if (parsedData.action === 'grow' && dude) {
			dude.grow();
		}

		if (parsedData.action === 'leave' && dude) {
			dude.leave();
		}

		if (parsedData.action === 'spawn-emote' && dude) {
			const emote = dudesStore.getProxiedEmoteUrl({
				type: '3rd_party_emote',
				value: 'https://cdn.7tv.app/emote/60b00d1f0d3a78a196f803e3/1x.gif',
			});
			dude.addEmotes([emote]);
		}

		if (parsedData.action === 'show-message' && dude) {
			dude.addMessage(`Hello, ${dudesSettingsStore.channelData!.channelDisplayName}! ${randomEmoji('emoticons')}`);
		}
	}

	async function spawnIframeDude() {
		if (
			!dudesStore.dudes ||
			!dudesSettingsStore.dudesSettings ||
			!dudesSettingsStore.channelData
		) return;

		const emote = dudesStore.getProxiedEmoteUrl({
			type: '3rd_party_emote',
			value: 'https://cdn.7tv.app/emote/65413498dc0468e8c1fbcdc6/1x.gif',
		});

		const dudeSprite = getSprite(dudesTwir, dudesSettingsStore.dudesSettings.overlay.defaultSprite);
		const dude = await dudesStore.dudes.createDude(dudesTwir, dudeSprite);
		dude.updateColor(DudesLayers.Body, '#8a2be2');
		dudesStore.updateDudeColors(dude);
		dude.addMessage(`Hello, ${dudesSettingsStore.channelData.channelDisplayName}! ${randomEmoji('emoticons')}`);
		dude.addEmotes([emote]);
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
