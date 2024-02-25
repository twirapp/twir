import type { MessageChunk } from '@twir/frontend-chat';
import { DudesSprite } from '@twir/types/overlays';
import type { DudesOverlayMethods, Dude } from '@twirapp/dudes/types';
import { defineStore, storeToRefs } from 'pinia';
import { computed, ref, watch } from 'vue';

import { getRandomSprite } from './dudes-config.js';
import { useDudesSettings } from './use-dudes-settings.js';

import { randomRgbColor } from '@/helpers.js';

export const useDudes = defineStore('dudes', () => {
	const dudesSettingsStore = useDudesSettings();
	const { dudesSettings } = storeToRefs(dudesSettingsStore);

	const dudes = ref<DudesOverlayMethods | null>(null);
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

	function createDude(
		displayName: string,
		userId: string,
		color?: string,
	) {
		if (!dudes.value || !dudesSettings.value) return;

		const dudeFromCanvas = dudes.value.dudes.get(displayName) as Dude;
		if (dudeFromCanvas) {
			return createDudeInstance(dudeFromCanvas);
		}

		if (
			dudesSettings.value.overlay.maxOnScreen !== 0 &&
			dudes.value.dudes.size === dudesSettings.value.overlay.maxOnScreen
		) return;

		const userSettings = getDudeUserSettings(userId);
		const dudeColor = userSettings?.dudeColor
			?? color
			?? dudesSettings.value.dudes.dude.color;

		let dudeSprite = (userSettings?.dudeSprite ?? dudesSettings.value.overlay.defaultSprite) as keyof typeof DudesSprite;
		if (dudeSprite === DudesSprite.random) {
			dudeSprite = getRandomSprite();
		}

		const dude = dudes.value.createDude(displayName, dudeSprite);
		dude.bodyTint(dudeColor);

		if (dudeSprite === DudesSprite.sith) {
			dude.cosmeticsTint(randomRgbColor());
		}

		const dudeInstance = createDudeInstance(dude);
		dudeInstance.isCreated = true;

		return dudeInstance;
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
			dude.spitEmotes([...new Set(emotes)]);
		}
	}

	function getProxiedEmoteUrl(messageChunk: MessageChunk): string {
		if (messageChunk.type === 'emoji') {
			const code = messageChunk.value.codePointAt(0)?.toString(16);
			return `https://cdn.frankerfacez.com/static/emoji/images/twemoji/${code}.png`;
		}

		return `${window.location.origin}/api/proxy?url=${messageChunk.value}`;
	}

	function deleteDude(displayName: string): void {
		if (!dudes.value) return;
		dudes.value.removeDude(displayName);
	}

	function getDudeUserSettings(userId: string) {
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
		deleteDude,
		createDude,
		showMessageDude,
		getProxiedEmoteUrl,
		isDudeOverlayReady,
	};
});
