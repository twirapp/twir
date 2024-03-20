import type { MessageChunk } from '@twir/frontend-chat';
import { DudesSprite } from '@twir/types/overlays';
import { DudesLayers } from '@twirapp/dudes';
import type { Dude, DudesMethods } from '@twirapp/dudes/types';
import { defineStore, storeToRefs } from 'pinia';
import { computed, ref, watch } from 'vue';

import { getSprite } from './dudes-config.js';
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

	async function createDude(
		{ userId, userName, color }: { userId: string, userName?: string, color?: string },
	) {
		if (!dudes.value || !dudesSettings.value) return;

		const actualDude = dudes.value.dudes.get(userId) as Dude;
		if (actualDude) {
			return createDudeInstance(actualDude);
		}

		if (
			dudesSettings.value.overlay.maxOnScreen !== 0 &&
			dudes.value.dudes.size === dudesSettings.value.overlay.maxOnScreen
		) return;

		const userSettings = requestDudeUserSettings(userId);
		if (!userSettings) {
			dudesSettingsStore.dudesUserSettings.set(userId, {
				userId,
				userDisplayName: userName,
				dudeColor: color,
			});
			return;
		}

		const dudeColor = userSettings.dudeColor
			?? color
			?? dudesSettings.value.dudes.dude.bodyColor;

		const dudeSprite = getSprite(userSettings?.dudeSprite ?? dudesSettings.value.overlay.defaultSprite);
		const dude = await dudes.value.createDude({
			id: userSettings.userId,
			name: userSettings.userDisplayName!,
			sprite: dudeSprite,
		});

		updateDudeColors(dude, dudeColor);

		const dudeInstance = createDudeInstance(dude);
		dudeInstance.isCreated = true;

		return dudeInstance;
	}

	function updateDudeColors(dude: Dude, color?: string): void {
		if (color) {
			dude.updateColor(DudesLayers.Body, color);
		}

		if (dude.config.sprite.name.startsWith(DudesSprite.girl)) {
			dude.updateColor(DudesLayers.Hat, '#FF0000');
		}

		if (dude.config.sprite.name.startsWith(DudesSprite.santa)) {
			dude.updateColor(DudesLayers.Hat, '#FFF');
		}

		if (dude.config.sprite.name.startsWith(DudesSprite.agent)) {
			dude.updateColor(DudesLayers.Cosmetics, '#8a2be2');
		}

		if (dude.config.sprite.name.startsWith(DudesSprite.sith)) {
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
