import { DISCORD_INVITE_URL } from '@twir/brand'

import type { ColorVariant } from './hero-chat-avatar.vue'

interface ChatMessage {
	type: 'message'
	sender: 'user' | 'bot'
	text: string
	variant?: ColorVariant
	sleep?: number
	replyMessages?: Message[]
}

interface ChatRedemption {
	type: 'redemption'
	text: string
	input: string
}

interface ChatSleep {
	type: 'sleep'
	ms: number
}

export type Message = ChatMessage | ChatRedemption | ChatSleep

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
	}
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
	}
}

function redemptionMessage(username: string, text: string, input: string): ChatRedemption {
	return {
		type: 'redemption',
		text: `<b>${username}</b> activated channel reward: ${text}`,
		input,
	}
}

function chatSleep(ms: number): ChatSleep {
	return {
		type: 'sleep',
		ms,
	}
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
]

export const liveChatMessages: Message[] = [
	userMessage('!game Minecraft', 'lime', botMessage('✅ Game changed to Minecraft')),

	userMessage('!song', 'blue', botMessage('Linkin Park — Numb')),

	userMessage('!watchtime', 'gray', botMessage(`You're watching stream for 210.2h`)),

	userMessage('!top time', 'lime', botMessage(`Jon × 210.3h · Mike × 167.2h · Alice × 125.9h · Ivan × 96.1h · Brian × 80.8h`)),

	userMessage('!top messages', 'blue', botMessage(`Jon × 3086 · Mike × 2647 · Alice × 2529 · Ivan × 2105 · Brian × 1500`)),

	userMessage('!me', 'red', botMessage(`210.4h watched · 1519 messages · 6507 used emotes · 0 used points · 0 songs requestes`)),

	userMessage('!sr sandstorm', 'pink', botMessage(`You're requested "Darude - Sandstorm"`)),

	userMessage('!followage', 'orange', botMessage(`You're following channel 2mo 1d 50m`)),

	userMessage('!age', 'purple', botMessage(`Your account age is 5mo 14d 21h 47m`)),

	userMessage('!title history', 'turquoise', botMessage(`Watching memes with my friends · Playing League of Legends with my friend`)),

	userMessage('!so @satont', 'lime', botMessage(`Check out amazing @Satont stream, was streaming Software and game development`)),

	userMessage('!permit @melkam', 'red', botMessage(`@melkam can post 1 link`)),

	userMessage('!permit @melkam 10', 'red', botMessage(`@melkam can post 10 links`)),

	userMessage(`!commands add discord ${DISCORD_INVITE_URL}`, 'blue', botMessage(`✅ Command !discord with response ${DISCORD_INVITE_URL} was added`)),

	userMessage(`!8ball Are you bot?`, 'blue', botMessage(`Yeah, probably.`)),

	userMessage(`!faceit`, 'blue', botMessage(`Level: 10 · ELO: 3000 · Matches: 1000 · Winrate: 50%`)),

	userMessage(`!weather Moscow`, 'blue', botMessage(`Moscow (RU): Snowy, 🌡️ -11.6°C, ☁️ 100%, 💦 94%, 💨 4.7 m/sec`)),
]
