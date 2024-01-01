<script setup lang="ts">
import {
	ChatBox,
} from '@twir/frontend-chat';
import type { Message } from '@twir/frontend-chat';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { useThirdPartyEmotes, type Opts as EmotesOpts } from '../../components/chat_tmi_emotes.js';
import { useChatOverlaySocket } from '../../sockets/chat_overlay.js';
import { type ChatSettings, useTmiChat, knownBots, ChatMessage } from '../../sockets/chat_tmi.ts';

const route = useRoute();
const apiKey = route.params.apiKey as string;

const messages = ref<Message[]>([]);
const maxMessages = ref(30);

const { settings } = useChatOverlaySocket(apiKey);

const emotesOpts = computed<EmotesOpts>(() => {
	return {
		channelName: settings.value.channelName,
		channelId: settings.value.channelId,
		ffz: true,
		bttv: true,
		sevenTv: true,
	};
});

useThirdPartyEmotes(emotesOpts);

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
	};
});

const { destroy } = useTmiChat(chatSettings);

onMounted(() => {
	document.body.style.overflow = 'hidden';
});

onUnmounted(async () => {
	destroy();
});
</script>

<template>
	<ChatBox :messages="messages" :settings="settings" />
</template>
