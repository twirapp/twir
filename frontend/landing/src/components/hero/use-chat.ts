import { nextTick, ref, watch } from 'vue';

import { initialChatMessages, liveChatMessages, type Message } from './hero-chat-messages.js';

function sleep(ms = 1000) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

export function useChat() {
	const messages = ref<Message[]>(initialChatMessages);

	watch(() => messages.value, async (msg) => {
		if (msg.length > 20) {
			await nextTick(() => {
				msg = msg.slice(10, msg.length - 1);
			});
		}
	}, { deep: true });

	async function processLiveMessages() {
		for (const message of liveChatMessages) {
			if (message.type === 'sleep') {
				await sleep(message.ms);
				continue;
			} else {
				await sleep(1500);
			}

			messages.value.push(message);

			if (message.type === 'message') {
				if (!message.replyMessages) continue;

				for (const msg of message.replyMessages) {
					if (msg.type === 'sleep') {
						await sleep(msg.ms);
						continue;
					} else {
						await sleep(650);
					}

					messages.value.push(msg);
				}
			}
		}
	}

	return {
		messages,
		processLiveMessages,
	};
}
