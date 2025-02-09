<script setup lang="ts">
import ChatMessageTail from 'assets/icons/chat-message-tail.svg'
import { onMounted, onUnmounted, ref } from 'vue'

import ChatAvatar from './hero-chat-avatar.vue'
import { useChat } from './use-chat.js'

const isEnabled = ref(false)
const { messages, startTimeout, stopTimeout } = useChat(isEnabled)

onMounted(() => startTimeout())
onUnmounted(() => stopTimeout())

function onMouseHover() {
	isEnabled.value = true
	stopTimeout()
}

function onMouseLeave() {
	isEnabled.value = false
	startTimeout()
}
</script>

<template>
	<div
		class="flex flex-col justify-end gap-[12px] w-full max-h-[540px] xl:max-w-lg relative -top-5"
		style="-webkit-mask-image: linear-gradient(0deg, #D9D9D9 75%, rgba(217, 217, 217, 0) 100%)"
		@mouseenter="onMouseHover"
		@mouseleave="onMouseLeave"
	>
		<TransitionGroup name="list">
			<div
				v-for="(message, index) of messages"
				:key="index"
				class="flex items-start gap-[16px] w-full"
			>
				<ChatAvatar
					v-if="message.type === 'message'"
					:fontControlled="false"
					:is-bot="message.sender === 'bot'"
					:variant="message.variant"
				/>

				<div
					v-if="message.type === 'message'" class="flex flex-col px-[16px] py-[10px] rounded-lg rounded-tl-none text-white relative" :class="[
						{
							'bg-[#534FDB]': message.sender === 'bot',
							'bg-[#232427]': message.sender === 'user',
						},
					]"
				>
					<ChatMessageTail
						:fontControlled="false"
						class="absolute h-[21px] w-[11px] top-0 -left-[10px]"
						:class="{
							'text-[#534FDB]': message.sender === 'bot',
							'text-[#232427]': message.sender === 'user',
						}"
					/>
					<span class="chat-message" v-html="message.text"></span>
				</div>

				<div
					v-if="message.type === 'redemption'"
					class="font-normal flex flex-col py-3 px-5 bg-[#4C47F5]/[.15] gap-2 rounded-md relative w-full"
				>
					<span class="text-sm leading-normal text-white/90" v-html="message.text"></span>
					<span class="font-semibold">{{ message.input }}</span>
					<span class="absolute bg-[#4C47F5] w-[2px] rounded-sm h-[calc(100%-24px)] left-0"></span>
				</div>
			</div>
		</TransitionGroup>
	</div>
</template>

<style scoped>
.chat-message {
	vertical-align: baseline;
}

.chat-emote {
	position: relative;
	display: inline-block;
	margin-left: 4px;
	margin-right: 4px;
	vertical-align: middle;
}

.list-move,
.list-enter-active {
	transition: all 0.3s ease;
}

.list-enter-from {
	opacity: 0;
	transform: translateY(50px);
}
</style>
