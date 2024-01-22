<script setup lang="ts">
import DudesOverlay from '@twirapp/dudes';
import type { DudesOverlayMethods } from '@twirapp/dudes/types';
import { storeToRefs } from 'pinia';
import { onMounted, ref } from 'vue';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { useChatOverlaySocket } from '@/composables/chat/use-chat-overlay-socket.js';
import { dudeAssets, dudeSprites } from '@/composables/dudes/dudes-config.js';
import { useChatTmi, type ChatSettings, type ChatMessage, knownBots } from '@/composables/tmi/use-chat-tmi.js';

const dudesRef = ref<DudesOverlayMethods<string> | null>(null);

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

	if (
		settings.value.hideCommands &&
		chatMessage.chunks.at(0)?.value.startsWith('!')
	) {
		return;
	}

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
	const randomDudeSprite = dudeSprites[Math.floor(Math.random() * dudeSprites.length - 1)];
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
  await dudesRef.value.initDudes(dudeAssets);

	const apiKey = route.params.apiKey as string;
	const overlayId = route.query.id as string;
	chatSocketStore.connect(apiKey, overlayId);
});
</script>

<template>
	<dudes-overlay ref="dudesRef" />
</template>

<style>
body {
  overflow: hidden;
}
</style>
