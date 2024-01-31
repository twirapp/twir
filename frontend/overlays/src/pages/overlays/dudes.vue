<script setup lang="ts">
import { useFontSource } from '@twir/fontsource';
import DudesOverlay from '@twirapp/dudes';
import type { DudesOverlayMethods } from '@twirapp/dudes/types';
import { storeToRefs } from 'pinia';
import { onMounted, ref, watch, computed } from 'vue';
import { useRoute } from 'vue-router';

import {
	dudesAssets,
	dudesSprites,
	dudesSounds,
	type DudeSprite,
} from '@/composables/dudes/dudes-config.js';
import { useDudesSettings } from '@/composables/dudes/use-dudes-settings.js';
import { useDudesSocket } from '@/composables/dudes/use-dudes-socket.js';
import { useChatTmi, type ChatSettings, type ChatMessage } from '@/composables/tmi/use-chat-tmi.js';

const dudesRef = ref<DudesOverlayMethods | null>(null);

const fontSource = useFontSource(false);
const dudesSettingStore = useDudesSettings();
const { dudesSettings } = storeToRefs(dudesSettingStore);
const route = useRoute();

watch(() => dudesSettings.value, async (settings) => {
	if (!settings || !dudesRef.value) return;

	dudesRef.value.updateSettings(settings);
	if (window.frameElement) {
		dudesRef.value.createDude('Twir', 'dude');
	}
	// dudesRef.value.clearDudes();
}, { deep: true });

watch(() => dudesSettings.value?.nameBox.fontFamily, async (fontFamily) => {
	if (!dudesSettings.value || !fontFamily) return;

	try {
		const loadedFont = await fontSource.loadFont(
			fontFamily,
			Number(dudesSettings.value.nameBox.fontWeight),
			dudesSettings.value.nameBox.fontStyle,
		);

		if (loadedFont) {
			dudesSettings.value.nameBox.fontFamily = loadedFont.id;
			dudesSettings.value.messageBox.fontFamily = loadedFont.id;
		}
	} catch (err) {
		console.error(err);
	}
});

const dudesSocket = useDudesSocket();

function onMessage(chatMessage: ChatMessage) {
	if (!dudesRef.value || chatMessage.type === 'system') return;

	const dudeName = chatMessage.senderDisplayName!;
	let dude = dudesRef.value.getDude(dudeName);
	if (dude) {
		dude.addMessage(chatMessage.rawMessage!);
	} else {
		dude = createNewDude(
			dudeName,
			chatMessage.rawMessage!,
			chatMessage.senderColor!,
		);
	}

	if (chatMessage.rawMessage?.startsWith('!') && dude) {
		const [command, argument] = chatMessage.rawMessage.split(' ');
		if (command === '!jump') {
			dude.jump();
		} else if (command === '!color' && argument) {
			dude.tint(argument);
		} else if (command === '!sprite' && argument) {
			if (!dudesSprites.includes(argument as DudeSprite)) return;
			dudesRef.value.createDude(dudeName, argument, {
				dude: {
					color: dude.color,
				},
			});
		}
	}
}

function createNewDude(sender: string, message: string, color: string) {
	if (!dudesRef.value) return;
	const randomDudeSprite = dudesSprites[Math.floor(Math.random() * dudesSprites.length - 1)];
	const dude = dudesRef.value.createDude(sender, randomDudeSprite);
	dude.tint(color);
	setTimeout(() => dude.addMessage(message), 2000);
	return dude;
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: dudesSettingStore.channelInfo?.channelId ?? '',
		channelName: dudesSettingStore.channelInfo?.channelName ?? '',
		onMessage,
	};
});

useChatTmi(chatSettings);

onMounted(async () => {
	const apiKey = route.params.apiKey as string;
	const overlayId = route.params.id as string;
	dudesSocket.connect(apiKey, overlayId);

	if (!dudesRef.value) return;
	await dudesRef.value.initDudes();
});
</script>

<template>
	<dudes-overlay
		ref="dudesRef"
		:assets="dudesAssets"
		:sounds="dudesSounds"
		:settings="dudesSettingStore.dudesSettings"
	/>
</template>

<style>
body {
	overflow: hidden;
}
</style>
