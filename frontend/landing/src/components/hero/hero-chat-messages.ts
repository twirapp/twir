import type { ColorVariant } from './hero-chat-avatar.vue';

type ChatMessage = {
	type: 'message';
	sender: 'user' | 'bot';
	text: string;
	variant?: ColorVariant;
	sleep?: number;
	replyMessages?: Message[];
};

type ChatRedemption = {
	type: 'redemption';
	text: string;
	input: string;
};

type ChatSleep = {
	type: 'sleep';
	ms: number;
};

export type Message = ChatMessage | ChatRedemption | ChatSleep;

function userMessage(
	text: string,
	variant?: ColorVariant,
	...messages: Message[]
): ChatMessage {
	return {
		type: 'message',
		sender: 'user',
		text,
		variant,
		replyMessages: messages,
	};
}

function botMessage(
	text: string,
	variant?: ColorVariant,
	...messages: Message[]
): ChatMessage {
	return {
		type: 'message',
		sender: 'bot',
		text,
		variant,
		replyMessages: messages,
	};
}

function redemptionMessage(username: string, text: string, input: string): ChatRedemption {
	return {
		type: 'redemption',
		text: `<b>${username}</b> activated channel reward: ${text}`,
		input,
	};
}

function chatSleep(ms: number): ChatSleep {
	return {
		type: 'sleep',
		ms,
	};
}

export const initialChatMessages: Message[] = [
	userMessage('Hello, World'),
	botMessage('Message from timer: follow to my socials!'),
	userMessage('!title Playling League of Legends with my friend'),
	botMessage('✅ Title succesfully changed.'),
	redemptionMessage('melkam', 'timeout chatter (1000 🪙)', 'Satont'),
	botMessage('melkam disabled chat for <b>Satont</b> for 5 minutes'),
	userMessage('!song'),
	botMessage('Linkin Park — Numb'),
	userMessage('!category LOL'),
	botMessage('✅ Category changed to League of Legends.'),
];

export const liveChatMessages: Message[] = [
	chatSleep(1000),
	userMessage('!game Minecraft', 'lime',
		chatSleep(2000),
		botMessage('✅ Game changed to Minecraft'),
	),

	chatSleep(2000),
	userMessage('!song', 'gray',
		chatSleep(2000),
		botMessage('Linkin Park — Numb'),
	),
];
