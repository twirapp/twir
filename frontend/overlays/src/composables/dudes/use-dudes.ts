import type { MessageChunk } from '@twir/frontend-chat';
import type { DudesOverlayMethods, Dude } from '@twirapp/dudes/types';
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
			createDude(userData.userDisplayName);
		}
	}

	function createDude(name: string, color?: string, messageChunks?: MessageChunk[]) {
		if (!dudes.value) return;

		const randomDudeSprite = dudesSprites[Math.floor(Math.random() * dudesSprites.length)];
		const dude = dudes.value.createDude(name, randomDudeSprite);

		if (color) {
			dude.tint(color);
		}

		if (messageChunks) {
			setTimeout(() => showMessageDude(dude, messageChunks), 1000);
		}

		return dude;
	}

	function showMessageDude(dude: Dude, messageChunks: MessageChunk[]): void {
		// TODO: fix types and fix ignore commands
		if (
			dudesSettings.value?.messageBox.ignoreCommands &&
			messageChunks?.at(0)?.value.startsWith('!')
		) {
			return;
		}

		const message = messageChunks
			.filter((chunk) => chunk.type === 'text' || chunk.type === 'emoji')
			.map((chunk) => chunk.value)
			.join(' ');
		const emotes = messageChunks
			.filter((chunk) => chunk.type === '3rd_party_emote' || chunk.type === 'emote')
			.map((chunk) => {
				// twitch emote
				if (chunk.type === 'emote') {
					return getProxiedEmoteUrl(`https://static-cdn.jtvnw.net/emoticons/v2/${chunk.value}/default/dark/1.0`);
				}
				return getProxiedEmoteUrl(chunk.value);
			});

		dude.addMessage(message);

		if (emotes.length) {
			dude.spitEmotes(emotes);
		}
	}

	function getProxiedEmoteUrl(url: string) {
		return `${window.location.origin}/api/proxy?url=${url}`;
	}

	function deleteDude(displayName: string): void {
		if (!dudes.value) return;
		dudes.value.removeDude(displayName);
	}

	function createNewDudeFromIframe() {
		if (window.frameElement) {
			createDude(
				'TWIR',
				'#8a2be2',
				[
					{
						type: 'text',
						value: `Hello, ${channelInfo.value!.channelDisplayName}!`,
					},
					{
						type: '3rd_party_emote',
						value: 'https://cdn.7tv.app/emote/63706216d49eb6644629aa52/1x.webp',
					},
				],
			);
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
		deleteDude,
		createDude,
		showMessageDude,
		createNewDudeFromIframe,
		isDudeOverlayReady,
	};
});
