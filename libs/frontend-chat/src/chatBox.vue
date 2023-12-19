<script setup lang="ts">
import { useFontSource } from '@twir/fontsource';
import { useWindowSize } from '@vueuse/core';
import { computed, nextTick, ref, toValue, watch } from 'vue';

import { getChatDirection, getUserColor } from './helpers.js';
import MessageHorizontalBoxed from './styles/boxed/message-horizontal.vue';
import MessageVerticalBoxed from './styles/boxed/message-vertical.vue';
import MessageHorizontalClean from './styles/clean/message-horizontal.vue';
import MessageVerticalClean from './styles/clean/message-vertical.vue';
import type { Message, Settings } from './types.js';

const props = defineProps<{
	messages: Message[]
	settings: Settings
}>();

const chatMessages = ref<HTMLDivElement>();

const fontSource = useFontSource();

watch(() => [
	props.settings.fontFamily,
	props.settings.fontWeight,
	props.settings.fontStyle,
], () => {
	fontSource.loadFont(
		props.settings.fontFamily,
		props.settings.fontWeight,
		props.settings.fontStyle,
	);
});

watch(() => props.messages.length, async () => {
	await nextTick();
	scrollByDirection(props.settings.direction);
});

watch(() => props.settings.direction, async (direction) => {
	await nextTick();
	scrollByDirection(direction);
});

function scrollToBottom() {
	if (!chatMessages.value) return;
	chatMessages.value.scrollIntoView(true);
}

function scrollToTop() {
	if (!chatMessages.value) return;
	chatMessages.value.scrollIntoView(false);
}

function scrollToLeft() {
	if (!chatMessages.value) return;
	chatMessages.value.scrollLeft += 999999;
}

function scrollToRight() {
	if (!chatMessages.value) return;
	chatMessages.value.scrollLeft -= 999999;
}

function scrollByDirection(direction: string) {
	if (direction === 'bottom') {
		scrollToBottom();
	}

	if (direction === 'top') {
		scrollToTop();
	}

	if (direction === 'left') {
		scrollToLeft();
	}

	if (direction === 'right') {
		scrollToRight();
	}
}

const chatDirection = computed(() => getChatDirection(props.settings.direction));

const messagesFlexDirection = computed(() => {
	switch (props.settings.direction) {
		case 'top':
			return 'column';
		case 'bottom':
			return 'column-reverse';
		case 'left':
			return 'row';
		case 'right':
			return 'row-reverse';
		default:
			return 'column';
	}
});

const chatMessageComponent = computed(() => {
	switch (props.settings.preset) {
		case 'boxed':
			if (chatDirection.value === 'horizontal') {
				return MessageHorizontalBoxed;
			}

			if (chatDirection.value === 'vertical') {
				return MessageVerticalBoxed;
			}

			return MessageVerticalBoxed;
		case 'clean':
			if (chatDirection.value === 'horizontal') {
				return MessageHorizontalClean;
			}

			if (chatDirection.value === 'vertical') {
				return MessageVerticalClean;
			}

			return MessageHorizontalClean;
		default:
			return MessageVerticalClean;
	}
});

const { height } = useWindowSize();
const windowHeight = computed(() => {
	if (getChatDirection(props.settings.direction) === 'horizontal') {
		return `${height.value}px`;
	}

	return 'auto';
});

const fontSize = computed(() => `${props.settings.fontSize}px`);
const fontFamily = computed(() => {
	return `"${props.settings.fontFamily}-${props.settings.fontWeight}-${props.settings.fontStyle}"`;
});
</script>

<template>
	<div class="chat">
		<div ref="chatMessages" class="messages">
			<TransitionGroup name="list">
				<component
					:is="chatMessageComponent"
					v-for="msg of messages"
					:key="msg.internalId"
					:msg="msg"
					:direction="chatDirection"
					:settings="toValue(settings)"
					:userColor="getUserColor(msg.senderColor)"
				/>
			</TransitionGroup>
		</div>
	</div>
</template>

<style>
.chat {
	height: 100vh;
	width: 100%;
	color: #fff;
	font-size: v-bind(fontSize);
	font-family: v-bind(fontFamily);
	font-weight: v-bind('settings.fontWeight');
	font-style: v-bind('settings.fontStyle');
	overflow: hidden;
	position: relative;
	background-color: v-bind('settings.chatBackgroundColor');
}

.messages {
	display: flex;
	flex-direction: v-bind(messagesFlexDirection);
	gap: 8px;
	overflow: hidden;
	height: v-bind(windowHeight);
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
