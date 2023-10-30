<script setup lang="ts">
import {
	ChatBox,
} from '@twir/frontend-chat';
import { useRoute } from 'vue-router';
import '@twir/frontend-chat/style.css';

import { useChatSocket } from '../sockets/chat.js';
import { useTmiChat } from '../sockets/chat_tmi.js';

const route = useRoute();
const apiKey = route.params.apiKey as string;

const { settings } = useChatSocket(apiKey);
const { messages } = useTmiChat(settings);
</script>

<template>
	<ChatBox
		:messages="messages"
		:settings="settings"
	/>
</template>
