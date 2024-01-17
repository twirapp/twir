import { nextTick, ref, watch, type Ref } from 'vue';

import { initialChatMessages, liveChatMessages, type Message } from './hero-chat-messages.js';

function sleep(ms = 1000) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

function getRandomNum(min: number, max: number) {
	return Math.floor(Math.random() * (max - min) + min);
}

export function useChat(isEnabled: Ref<boolean>) {
	const currentMessageIndex = ref(0);
	const timeoutId = ref<ReturnType<typeof setTimeout>>(null);
	const messages = ref<Message[]>(initialChatMessages);

	watch(() => messages.value, async (msg) => {
		if (msg.length > 20) {
			await nextTick(() => {
				msg = msg.slice(10, msg.length - 1);
			});
		}
	}, { deep: true });

	function startTimeout() {
		timeoutId.value = setTimeout(() => processLiveMessages(), getRandomNum(1500, 2000));
	}

	function stopTimeout() {
		if (!timeoutId.value) return;
		clearTimeout(timeoutId.value);
		timeoutId.value = null;
	}

	async function processLiveMessages() {
		if (isEnabled.value) return;

		currentMessageIndex.value = currentMessageIndex.value % liveChatMessages.length;
		const message = liveChatMessages[currentMessageIndex.value];
		currentMessageIndex.value += 1;

		if (message.type === 'sleep') {
			await sleep(message.ms);
		}

		if (message.type === 'message') {
			messages.value.push(message);

			for (const msg of message.replyMessages) {
				if (msg.type === 'sleep') {
					await sleep(msg.ms);
					continue;
				} else {
					await sleep(800);
				}

				messages.value.push(msg);
			}
		}

		startTimeout();
	}

	return {
		messages,
		startTimeout,
		stopTimeout,
	};
}
