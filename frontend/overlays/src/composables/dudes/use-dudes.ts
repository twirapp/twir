import type { MessageChunk } from '@twir/frontend-chat';
import type { DudesOverlayMethods, Dude } from '@twirapp/dudes/types';
import { defineStore, storeToRefs } from 'pinia';
import { computed, ref, watch } from 'vue';

import { dudesSprites } from './dudes-config.js';
import { useDudesSettings } from './use-dudes-settings.js';

import { randomRgbColor } from '@/helpers.js';
import type { UserData } from '@/types.js';

export const useDudes = defineStore('dudes', () => {
	const dudesSettigsStore = useDudesSettings();
	const { dudesSettings } = storeToRefs(dudesSettigsStore);

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
			createDude(userData.userDisplayName);
		}
	}

	function createDude(
		name: string,
		color?: string,
		messageChunks?: MessageChunk[],
	): void {
		if (!dudes.value) return;

		const randomDudeSprite = dudesSprites[Math.floor(Math.random() * dudesSprites.length)];
		const dude = dudes.value.createDude(name, randomDudeSprite);

		if (color) {
			dude.bodyTint(color);
		}

		if (randomDudeSprite === 'sith') {
			dude.cosmeticsTint(randomRgbColor());
		}

		if (messageChunks) {
			setTimeout(() => showMessageDude(dude, messageChunks), 1000);
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

	watch(() => dudes.value, async (dudes) => {
		if (!dudes) return;
		await dudes.initDudes();
		isDudeReady.value = true;
	});

	return {
		dudes,
		jumpDude,
		deleteDude,
		createDude,
		showMessageDude,
		getProxiedEmoteUrl,
		isDudeOverlayReady,
	};
});
