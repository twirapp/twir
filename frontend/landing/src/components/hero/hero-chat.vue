<script setup lang="ts">
import { ref } from 'vue';

import ChatAvatar from './hero-chat-avatar.vue';
import { type ChatMessage } from './hero-chat-messages.js';

import ChatMessageTail from '@/assets/chat-message-tail.svg?component';

const props = defineProps<{
	messages: ChatMessage[]
}>();

const chatMessages = ref<ChatMessage[]>(props.messages);
</script>

<template>
	<div
		class="flex flex-col justify-end gap-[12px] w-full max-h-[540px] xl:max-w-lg relative -top-5"
		style="-webkit-mask-image: linear-gradient(0deg, #D9D9D9 75%, rgba(217, 217, 217, 0) 100%)"
	>
		<TransitionGroup name="list">
			<div
				v-for="(message, index) of chatMessages"
				:key="index"
				class="flex items-start gap-[16px] w-full"
			>
				<chat-avatar
					v-if="message.type === 'message'"
					:is-bot="message.sender === 'bot'"
				/>

				<div
					v-if="message.type === 'message'" :class="[
						'flex flex-col px-[16px] py-[10px] rounded-lg rounded-tl-none text-white relative',
						{
							'bg-[#534FDB]': message.sender === 'bot',
							'bg-[#232427]': message.sender === 'user',
						}
					]"
				>
					<ChatMessageTail
						:class="{
							'absolute h-[21px] top-0 -left-[10px]': true,
							'fill-[#534FDB]': message.sender === 'bot',
							'fill-[#232427]': message.sender === 'user',
						}"
					/>
					<span v-html="message.text"></span>
				</div>

				<div
					v-if="message.type === 'redemption'"
					class="font-normal flex flex-col py-3 px-5 bg-[#4C47F5]/[.15] gap-2 rounded-md relative w-full"
				>
					<span class="text-sm leading-normal text-white/90" v-html="message.text"></span>
					<span class="font-semibold">{{ message.user }}</span>
					<span class="absolute bg-[#4C47F5] w-[2px] rounded-sm h-[calc(100%-24px)] left-0"></span>
				</div>
			</div>
		</TransitionGroup>
	</div>
</template>

<style>
.list-move,
.list-enter-active,
.list-leave-active {
	transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
	opacity: 0;
	transform: translateY(50px);
}

.list-leave-active {
	position: absolute;
}
</style>
