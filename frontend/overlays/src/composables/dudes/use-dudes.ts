import type { DudesOverlayMethods } from '@twirapp/dudes/types';
import { defineStore, storeToRefs } from 'pinia';
import { computed, ref, watch } from 'vue';

import { dudesSprites } from './dudes-config.js';
import { useDudesSettings } from './use-dudes-settings.js';

import type { UserData } from '@/types.js';

export const useDudes = defineStore('dudes', () => {
	const dudesSettigsStore = useDudesSettings();
	const { channelInfo, dudesSettings } = storeToRefs(dudesSettigsStore);

	const dudes = ref<DudesOverlayMethods | null>(null);
	const isDudeReady = ref(false);
	const isDudeOverlayReady = computed(() => {
		return dudes.value && dudesSettings.value && isDudeReady.value;
	});

	function jumpDude(userData: UserData): void {
		if (!dudes.value) return;

		const dude = dudes.value.getDude(userData.userDisplayName);
		if (dude) {
			dude.jump();
		} else {
			createNewDude(userData.userDisplayName);
		}
	}

	function createNewDude(name: string, color?: string, message?: string) {
		if (!dudes.value) return;

		const randomDudeSprite = dudesSprites[Math.floor(Math.random() * dudesSprites.length - 1)];
		const dude = dudes.value.createDude(name, randomDudeSprite);

		if (color) {
			dude.tint(color);
		}

		if (message) {
			setTimeout(() => dude.addMessage(message), 1000);
		}

		return dude;
	}

	function createNewDudeFromIframe() {
		if (window.frameElement) {
			createNewDude('TWIR', '#8a2be2', `Hello, ${channelInfo.value!.channelDisplayName}!`);
		}
	}

	watch(() => dudes.value, async (dudes) => {
		if (!dudes) return;
		await dudes.initDudes();
		isDudeReady.value = true;
	});

	watch([isDudeOverlayReady, dudesSettings], ([isReady, settings]) => {
		if (!isReady || !settings || !dudes.value) return;
		dudes.value.clearDudes();
		dudes.value.updateSettings(settings);
		createNewDudeFromIframe();
	});

	return {
		dudes,
		jumpDude,
		createNewDude,
		createNewDudeFromIframe,
		isDudeOverlayReady,
	};
});
