import type { MessageChunk } from '@twir/frontend-chat';
import { DudesSprite } from '@twir/types/overlays';
import { DudesLayers } from '@twirapp/dudes';
import type { DudesMethods, Dude } from '@twirapp/dudes/types';
import { defineStore, storeToRefs } from 'pinia';
import { computed, ref, watch } from 'vue';

import { dudesTwir, getSprite } from './dudes-config.js';
import { useDudesSettings } from './use-dudes-settings.js';

import { randomRgbColor } from '@/helpers.js';

export const useDudes = defineStore('dudes', () => {
	const dudesSettingsStore = useDudesSettings();
	const { dudesSettings } = storeToRefs(dudesSettingsStore);

	const dudes = ref<DudesMethods | null>(null);
	const isDudeReady = ref(false);
	const isDudeOverlayReady = computed(() => {
		return dudes.value && dudesSettings.value && isDudeReady.value;
	});

	function createDudeInstance(dude: Dude) {
		return {
			dude,
			isCreated: false,
			showMessage: function (messageChunks: MessageChunk[]) {
				if (this.isCreated) {
					setTimeout(() => showMessageDude(this.dude, messageChunks), 1000);
				} else {
					showMessageDude(this.dude, messageChunks);
				}
			},
		};
	}

	async function createDude(name: string, userId: string, color?: string) {
		if (!dudes.value || !dudesSettings.value) return;

		const dudeFromCanvas = dudes.value.dudes.get(name) as Dude;
		if (dudeFromCanvas) {
			return createDudeInstance(dudeFromCanvas);
		}

		if (
			dudesSettings.value.overlay.maxOnScreen !== 0 &&
			dudes.value.dudes.size === dudesSettings.value.overlay.maxOnScreen
		) return;

		const userSettings = requestDudeUserSettings(userId);
		const dudeColor = userSettings?.dudeColor
			?? color
			?? dudesSettings.value.dudes.dude.bodyColor;

		const dudeSprite = getSprite(name, userSettings?.dudeSprite ?? dudesSettings.value.overlay.defaultSprite);
		const dude = await dudes.value.createDude(name, dudeSprite);
		dude.updateColor(DudesLayers.Body, dudeColor);
		updateDudeColors(dude);

		const dudeInstance = createDudeInstance(dude);
		dudeInstance.isCreated = true;

		return dudeInstance;
	}

	function updateDudeColors(dude: Dude): void {
		const isTwir = dude.spriteData.name.startsWith(dudesTwir);

		if (dude.spriteData.name.startsWith(DudesSprite.girl) || isTwir) {
			dude.updateColor(DudesLayers.Hat, '#FF0000');
		}

		if (dude.spriteData.name.startsWith(DudesSprite.santa) || isTwir) {
			dude.updateColor(DudesLayers.Hat, '#FFF');
		}

		if (dude.spriteData.name.startsWith(DudesSprite.agent) || isTwir) {
			dude.updateColor(DudesLayers.Cosmetics, '#8a2be2');
		}

		if (dude.spriteData.name.startsWith(DudesSprite.sith) || isTwir) {
			dude.updateColor(DudesLayers.Cosmetics, randomRgbColor());
		}
	}

	function showMessageDude(dude: Dude, messageChunks: MessageChunk[]): void {
		if (
			dudesSettings.value?.ignore.ignoreCommands &&
			messageChunks?.at(0)?.value.startsWith('!')
		) {
			return;
		}

		const message = messageChunks
			.filter((chunk) => chunk.type === 'text')
			.map((chunk) => chunk.value)
			.join(' ');

		dude.addMessage(message);

		const emotes = messageChunks
			.filter((chunk) => chunk.type !== 'text')
			.map(getProxiedEmoteUrl);

		if (emotes.length) {
			dude.addEmotes([...new Set(emotes)]);
		}
	}

	function getProxiedEmoteUrl(messageChunk: MessageChunk): string {
		if (messageChunk.type === 'emoji') {
			const code = messageChunk.value.codePointAt(0)?.toString(16);
			return `https://cdn.frankerfacez.com/static/emoji/images/twemoji/${code}.png`;
		}

		return `${window.location.origin}/api/proxy?url=${messageChunk.value}`;
	}

	function requestDudeUserSettings(userId: string) {
		const userSettings = dudesSettingsStore.dudesUserSettings.get(userId);
		if (userSettings) return userSettings;
		document.dispatchEvent(new CustomEvent<string>('get-user-settings', { detail: userId }));
	}

	watch(() => dudes.value, async (dudes) => {
		if (!dudes) return;
		await dudes.initDudes();
		isDudeReady.value = true;
	});

	return {
		dudes,
		createDude,
		showMessageDude,
		updateDudeColors,
		getProxiedEmoteUrl,
		isDudeOverlayReady,
	};
});
