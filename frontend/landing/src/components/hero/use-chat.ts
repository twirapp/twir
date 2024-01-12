import { nextTick, onMounted, ref, watch } from 'vue';

import { initialChatMessages, liveChatMessages, type Message } from './hero-chat-messages.js';

function sleep(ms = 1000) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

export function useChat() {
	const messages = ref<Message[]>(initialChatMessages);

	function getNextMessage() {
		const message = liveChatMessages.shift();
		return message;
	}

	async function queueMessage() {
		const message = getNextMessage();
		if (!message) return;

		if (message.type === 'sleep') {
			await sleep(message.ms);
			await queueMessage();
			return;
		}

		messages.value.push(message);

		if (message.type === 'message') {
			if (message.replyMessages) {
				for (const msg of message.replyMessages) {
					if (msg.type === 'sleep') {
						await sleep(msg.ms);
						continue;
					}

					messages.value.push(msg);
				}
			}
		}

		await queueMessage();
	}

	watch(() => messages.value, async (msg) => {
		if (msg.length > 20) {
			await nextTick(() => {
				msg = msg.slice(10, msg.length - 1);
			});
		}
	}, { deep: true });

	onMounted(() => queueMessage());

	return {
		messages,
	};
}
