<script setup lang="ts">
import DudesOverlay from '@twirapp/dudes';
import { storeToRefs } from 'pinia';
import { onMounted, onUnmounted, computed, watch } from 'vue';
import { useRoute } from 'vue-router';

import {
	dudesAssets,
	dudesSounds,
	assetsLoadOptions,
} from '@/composables/dudes/dudes-config.js';
import { useDudesIframe } from '@/composables/dudes/use-dudes-iframe.js';
import { useDudesSettings } from '@/composables/dudes/use-dudes-settings.js';
import { useDudesSocket } from '@/composables/dudes/use-dudes-socket.js';
import { useDudes } from '@/composables/dudes/use-dudes.js';
import { useChatTmi, type ChatSettings, type ChatMessage } from '@/composables/tmi/use-chat-tmi.js';

const route = useRoute();

const dudesStore = useDudes();
const { dudes, isDudeOverlayReady } = storeToRefs(dudesStore);

const dudesSettingStore = useDudesSettings();
const { channelData, dudesSettings } = storeToRefs(dudesSettingStore);

const dudesSocketStore = useDudesSocket();

const iframe = useDudesIframe();

watch([isDudeOverlayReady, dudesSettings], ([isReady, settings]) => {
	if (!isReady || !settings || !dudes.value) return;
	dudes.value.updateSettings(settings.dudes);

	if (iframe.isIframe) {
		iframe.spawnIframeDude();
	}
});

function onMessage(chatMessage: ChatMessage): void {
	if (!dudes.value || chatMessage.type === 'system') return;

	if (
		dudesSettings.value?.ignore.ignoreUsers &&
		dudesSettings.value.ignore.users.includes(chatMessage.senderId!)
	) {
		return;
	}

	const displayName = chatMessage.senderDisplayName!;
	const color = chatMessage.senderColor ?? dudesSettings.value?.dudes.dude.color;

	const dude = dudes.value.getDude(displayName);
	if (dude) {
		dudesStore.showMessageDude(dude, chatMessage.chunks);
	} else {
		dudesStore.createDude(displayName, color, chatMessage.chunks);
	}
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: channelData.value?.channelId ?? '',
		channelName: channelData.value?.channelName ?? '',
		emotes: {
			isSmaller: true,
			bttv: true,
			ffz: true,
			sevenTv: true,
		},
		onMessage,
	};
});

const { destroy } = useChatTmi(chatSettings);

onMounted(async () => {
	const apiKey = route.params.apiKey as string;
	const overlayId = route.params.id as string;
	dudesSocketStore.connect(apiKey, overlayId);
	iframe.connect();
});

onUnmounted(() => {
	destroy();
	iframe.destroy();
});
</script>

<template>
	<dudes-overlay
		ref="dudes"
		:assets-load-options="assetsLoadOptions"
		:assets="dudesAssets"
		:sounds="dudesSounds"
	/>
</template>

<style>
body {
	overflow: hidden;
}
</style>
