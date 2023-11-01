<script setup lang="ts">
import { computed, nextTick, ref, toValue, watch } from 'vue';

import ChatMessageStyleBoxed from './styles/boxed.vue';
import ChatMessageStyleClean from './styles/clean.vue';
import type { Message, Settings } from './types.js';

const chatElement = ref<HTMLDivElement>();

const props = defineProps<{
	messages: Message[]
	settings: Settings
}>();

watch(() => props.messages.length, async () => {
	if (!chatElement.value || props.settings.reverseMessages) return;

	await nextTick();
	chatElement.value.scrollTo(0, chatElement.value.scrollHeight);
});

const chatMessageComponent = computed(() => {
	switch (props.settings.preset) {
		case 'boxed':
			return ChatMessageStyleBoxed;
		case 'clean':
		default:
			return ChatMessageStyleClean;
	}
});

const fontSize = computed(() => `${props.settings.fontSize}px`);

const defaultFont = 'Roboto';
const fontFamily = computed(() => {
	try {
		const [family] = props.settings.fontFamily.split(':');

		return family || defaultFont;
	} catch (e) {
		return defaultFont;
	}
});
const fontUrl = computed(() => {
	return `https://fonts.googleapis.com/css?family=${fontFamily.value}`;
});

const messagesDirection = computed(() => {
	return !props.settings.reverseMessages ? 'column' : 'column-reverse';
});
</script>

<template>
	<component :is="'style'">
		@import url('{{ fontUrl }}')
	</component>
	<div ref="chatElement" class="chat">
		<TransitionGroup name="list" tag="div" class="messages">
			<component
				:is="chatMessageComponent"
				v-for="(msg, index) of messages"
				:key="index"
				:msg="msg"
				:settings="toValue(settings)"
			/>
		</TransitionGroup>
	</div>
</template>

<style scoped>
.chat {
	max-height: 100vh;
  width: 100%;
  color: #fff;
  font-size: v-bind(fontSize);
	font-family: v-bind(fontFamily);
	overflow: hidden;
}

.chat .messages {
	display: flex;
	flex-direction: v-bind(messagesDirection);
	gap: 8px;
	overflow: hidden;
}

.chat .message .text .emote {
	max-height: 1em;
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
