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
	if (props.settings.reverseMessages) return;
	await nextTick();
	scrollToBottom();
});

function scrollToBottom() {
	chatElement.value?.scrollIntoView(false);
}

function scrollToTop() {
	chatElement.value?.scrollIntoView(true);
}

watch(() => props.settings.reverseMessages, (newValue) => {
	if (newValue) {
		scrollToTop();
	} else {
		scrollToBottom();
	}
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
	} catch {
		return defaultFont;
	}
});
const fontUrl = computed(() => {
	return `https://fonts.googleapis.com/css?family=${fontFamily.value}`;
});

const messagesDirection = computed(() => {
	return props.settings.reverseMessages ? 'column-reverse' : 'column';
});
</script>

<template>
	<div ref="chatElement" class="chat">
		<component :is="'style'">
			@import url('{{ fontUrl }}')
		</component>
		<TransitionGroup name="list" tag="div" class="messages">
			<component
				:is="chatMessageComponent"
				v-for="msg of messages"
				:key="msg.internalId"
				:msg="msg"
				:settings="toValue(settings)"
			/>
		</TransitionGroup>
	</div>
</template>

<style>
.chat {
	height: 100dvh;
  width: 100%;
  color: #fff;
  font-size: v-bind(fontSize);
	font-family: v-bind(fontFamily);
	overflow: hidden;
	position: relative;
	background-color: v-bind('settings.chatBackgroundColor');
}

.chat .messages {
	display: flex;
	flex-direction: v-bind(messagesDirection);
	gap: 8px;
	overflow: hidden;
	padding-bottom: 5px;
}

.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}
</style>
