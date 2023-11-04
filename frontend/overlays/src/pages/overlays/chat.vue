<script setup lang="ts">
import {
  ChatBox,
} from '@twir/frontend-chat';
import { useRoute } from 'vue-router';
import '@twir/frontend-chat/style.css';

import { useChatOverlaySocket } from '../../sockets/chat_overlay.js';
import { useTmiChat } from '../../sockets/chat_tmi.ts';

const route = useRoute();
const apiKey = route.params.apiKey as string;

const { settings } = useChatOverlaySocket(apiKey);
const { messages } = useTmiChat(settings);
</script>

<template>
	<ChatBox
		:messages="messages"
		:settings="settings"
	/>
</template>
