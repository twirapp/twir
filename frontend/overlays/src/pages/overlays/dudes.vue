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
const { channelInfo, dudesSettings } = storeToRefs(dudesSettingStore);

const dudesSocketStore = useDudesSocket();

function onMessage(chatMessage: ChatMessage) {
	if (!dudes.value || chatMessage.type === 'system') return;

	const dudeName = chatMessage.senderDisplayName!;
	const dude = dudes.value.getDude(dudeName);
	if (dude) {
		dude.addMessage(chatMessage.rawMessage!);
	} else {
		dudesStore.createNewDude(
			dudeName,
			chatMessage.senderColor,
			chatMessage.rawMessage,
		);
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
	<template v-if="dudesSettings">
		<dudes-overlay
			ref="dudes"
			:assets="dudesAssets"
			:sounds="dudesSounds"
			:settings="dudesSettings"
		/>
	</template>
</template>

<style>
body {
	overflow: hidden;
}
</style>
