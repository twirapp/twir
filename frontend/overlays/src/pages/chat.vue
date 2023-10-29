<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import ChatMessage from '../components/chatMessage.vue';
import { useChatSocket } from '../sockets/chat.js';
import { useTmiChat } from '../sockets/chat_tmi.js';

const route = useRoute();
const apiKey = route.params.apiKey as string;

const chat = useChatSocket(apiKey);

const channelName = computed(() => chat.settings.channelName);
const channelId = computed(() => chat.settings.channelId);
const messageTimeout = computed(() => chat.settings.messageTimeout);
const { messages } = useTmiChat(
	channelName,
	channelId,
	messageTimeout,
);
</script>

<template>
	<div class="chat">
		<TransitionGroup name="list" tag="div" class="messages">
			<ChatMessage
				v-for="(msg, index) of messages"
				:key="index"
				:msg="msg"
				:settings="chat.settings"
			/>
		</TransitionGroup>
	</div>
</template>

<style scoped>
@import url(https://fonts.googleapis.com/css?family=Roboto:700);

.chat {
  height: 100%;
  width: 100%;
  background-color: #000;
  color: #fff;
  font-size: 20px;
	font-family: 'Roboto';
}

.chat .messages {
	display: flex;
	flex-direction: column;
	gap: 8px;
	overflow: hidden;
}

.list-move, /* apply transition to moving elements */
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

/* ensure leaving items are taken out of layout flow so that moving
   animations can be calculated correctly. */
.list-leave-active {
  position: absolute;
}
</style>
