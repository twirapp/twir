<script setup lang="ts">
import { ChatBox } from '@twir/frontend-chat';
import type { Message } from '@twir/frontend-chat';
import { storeToRefs } from 'pinia';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { useChatOverlaySocket } from '@/composables/chat/use-chat-overlay-socket.js';
import {
	type ChatMessage,
	type ChatSettings,
	useChatTmi,
	knownBots,
} from '@/composables/tmi/use-chat-tmi.js';
import {
	useThirdPartyEmotes,
	type ThirdPartyEmotesOptions,
} from '@/composables/tmi/use-third-party-emotes.js';

const route = useRoute();

const messages = ref<Message[]>([]);
const maxMessages = ref(30);

const chatSocketStore = useChatOverlaySocket();
const { settings } = storeToRefs(chatSocketStore);

const emotesOptions = computed<ThirdPartyEmotesOptions>(() => {
	return {
		channelName: settings.value.channelName,
		channelId: settings.value.channelId,
		ffz: true,
		bttv: true,
		sevenTv: true,
	};
});

useThirdPartyEmotes(emotesOptions);

const removeMessageByInternalId = (id: string) => {
	messages.value = messages.value.filter(m => m.internalId !== id);
};

const removeMessageById = (id: string) => {
	messages.value = messages.value.filter(m => m.id !== id);
};

const removeMessageByUserName = (userName: string) => {
	messages.value = messages.value.filter(m => m.sender !== userName);
};

const onMessage = (m: ChatMessage) => {
	if (m.sender && settings.value.hideBots && knownBots.has(m.sender)) {
		return;
	}

	if (settings.value.hideCommands && m.chunks.at(0)?.value.startsWith('!')) {
		return;
	}

	const internalId = crypto.randomUUID();

	const showDelay = settings.value.messageShowDelay ?? settings.value.messageShowDelay;

	if (messages.value.length >= maxMessages.value) {
		messages.value = messages.value.slice(1);
	}

	setTimeout(() => {
		messages.value.push({
			...m,
			isItalic: m.isItalic ?? false,
			createdAt: new Date(),
			internalId,
			isAnnounce: m.isAnnounce ?? false,
		});
	}, showDelay * 1000);

	const hideTimeout = m.messageHideTimeout ?? settings.value.messageHideTimeout;

	if (hideTimeout) {
		setTimeout(() => removeMessageByInternalId(internalId), hideTimeout * 1000);
	}
};

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: settings.value.channelId,
		channelName: settings.value.channelName,
		onMessage,
		onRemoveMessage: removeMessageById,
		onRemoveMessageByUser: removeMessageByUserName,
		onChatClear: () => {
			messages.value = [];
		},
	};
});

const chatTmiStore = useChatTmi(chatSettings);

onMounted(() => {
	document.body.style.overflow = 'hidden';

	const apiKey = route.params.apiKey as string;
	const overlayId = route.query.id as string;
	chatSocketStore.connect(apiKey, overlayId);
});

onUnmounted(async () => {
	chatSocketStore.destroy();
	chatTmiStore.destroy();
});
</script>

<template>
	<chat-box :messages="messages" :settings="settings" />
</template>
