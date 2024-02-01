<script setup lang="ts">
import DudesOverlay from '@twirapp/dudes';
import { storeToRefs } from 'pinia';
import { onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';

import {
	dudesAssets,
	dudesSounds,
} from '@/composables/dudes/dudes-config.js';
import { useDudesSettings } from '@/composables/dudes/use-dudes-settings.js';
import { useDudesSocket } from '@/composables/dudes/use-dudes-socket.js';
import { useDudes } from '@/composables/dudes/use-dudes.js';
import { useChatTmi, type ChatSettings, type ChatMessage } from '@/composables/tmi/use-chat-tmi.js';

const route = useRoute();

const dudesStore = useDudes();
const { dudes } = storeToRefs(dudesStore);

const dudesSettingStore = useDudesSettings();
const { dudesSettings, channelInfo } = storeToRefs(dudesSettingStore);

const dudesSocketStore = useDudesSocket();

function onMessage(chatMessage: ChatMessage) {
	if (!dudes.value || chatMessage.type === 'system') return;

	const message = chatMessage.rawMessage!;
	const displayName = chatMessage.senderDisplayName!;
	const color = chatMessage.senderColor ?? dudesSettings.value?.dude.color;

	const dude = dudes.value.getDude(displayName);
	if (dude) {
		dudesStore.showDudeMessage(dude, message);
	} else {
		dudesStore.createNewDude(displayName, color, message);
	}
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: channelInfo.value?.channelId ?? '',
		channelName: channelInfo.value?.channelName ?? '',
		onMessage,
	};
});

useChatTmi(chatSettings);

onMounted(async () => {
	const apiKey = route.params.apiKey as string;
	const overlayId = route.params.id as string;
	dudesSocketStore.connect(apiKey, overlayId);
});
</script>

<template>
	<dudes-overlay
		ref="dudes"
		:assets="dudesAssets"
		:sounds="dudesSounds"
	/>
</template>

<style>
body {
	overflow: hidden;
}
</style>
