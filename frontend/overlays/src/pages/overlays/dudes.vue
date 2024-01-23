<script setup lang="ts">
import DudesOverlay from '@twirapp/dudes';
import type { DudesOverlayMethods, DudesSettings } from '@twirapp/dudes/types';
import { storeToRefs } from 'pinia';
import { onMounted, reactive, ref, watch } from 'vue';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { useChatOverlaySocket } from '@/composables/chat/use-chat-overlay-socket.js';
import { dudesAssets, dudesSprites, type DudeSprite } from '@/composables/dudes/dudes-config.js';
import { useChatTmi, type ChatSettings, type ChatMessage, knownBots } from '@/composables/tmi/use-chat-tmi.js';

const dudesSettings = reactive<DudesSettings>({
  dude: {
    color: '#969696',
    maxLifeTime: 1000 * 60 * 5,
    gravity: 500,
    scale: 6,
  },
  messageBox: {
    borderRadius: 10,
    boxColor: '#000000',
    fontFamily: 'Courier New',
    fontSize: 20,
    padding: 10,
    showTime: 5 * 1000,
    fill: '#ffffff',
  },
  nameBox: {
    fontFamily: 'Arial',
    fontSize: 20,
    fill: '#ffffff',
    lineJoin: 'round',
    strokeThickness: 4,
    stroke: '#000000',
    fillGradientStops: [0],
    fillGradientType: 0,
    fontStyle: 'normal',
    fontVariant: 'normal',
    fontWeight: 'normal',
    dropShadow: false,
    dropShadowAlpha: 1,
    dropShadowAngle: 0,
    dropShadowBlur: 0.1,
    dropShadowDistance: 10,
    dropShadowColor: '#3ac7d9',
  },
});

const dudesRef = ref<DudesOverlayMethods<DudeSprite> | null>(null);

watch(dudesSettings, (settings) => {
  if (!dudesRef.value) return;
  dudesRef.value.updateSettings(settings);
});

const route = useRoute();
const chatSocketStore = useChatOverlaySocket();
const { settings } = storeToRefs(chatSocketStore);

function onMessage(chatMessage: ChatMessage) {
	if (!dudesRef.value || chatMessage.type === 'system') return;

	if (
		chatMessage.sender &&
		settings.value.hideBots &&
		knownBots.has(chatMessage.sender)
	) {
		return;
	}

	if (chatMessage.rawMessage?.startsWith('!')) {
		const [command, argument] = chatMessage.rawMessage.split(' ');
		if (command === '!jump') {
			dudesRef.value.getDude(chatMessage.senderDisplayName!)?.jump();
		}

		if (command === '!color') {
			dudesRef.value.getDude(chatMessage.senderDisplayName!)?.tint(argument);
		}

		return;
	}

	// if (
	// 	settings.value.hideCommands &&
	// 	chatMessage.chunks.at(0)?.value.startsWith('!')
	// ) {
	// 	return;
	// }

	const currentDude = dudesRef.value.getDude(chatMessage.senderDisplayName!);
	if (currentDude) {
		currentDude.addMessage(chatMessage.rawMessage!);
		currentDude.tint(chatMessage.senderColor!);
	} else {
		createNewDude(
			chatMessage.senderDisplayName!,
			chatMessage.rawMessage!,
			chatMessage.senderColor!,
		);
	}
}

function createNewDude(sender: string, message: string, color: string) {
	const randomDudeSprite = dudesSprites[Math.floor(Math.random() * dudesSprites.length - 1)];
	const dude = dudesRef.value!.createDude(sender, randomDudeSprite);
	dude.tint(color);
	setTimeout(() => dude.addMessage(message), 2000);
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: settings.value.channelId,
		channelName: settings.value.channelName,
		onMessage,
	};
});

useChatTmi(chatSettings);

onMounted(async () => {
  if (!dudesRef.value) return;
  await dudesRef.value.initDudes();

	const apiKey = route.params.apiKey as string;
	chatSocketStore.connect(apiKey);
});
</script>

<template>
	<dudes-overlay ref="dudesRef" :assets="dudesAssets" :settings="dudesSettings" />
</template>

<style>
body {
  overflow: hidden;
}
</style>
