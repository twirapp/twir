<script setup lang="ts">
import { ChatBox } from '@twir/frontend-chat'

import type { Message, Settings } from '@twir/frontend-chat'

const settings: Settings = {
	channelId: '1',
	channelName: 'test',
	channelDisplayName: 'test',
	animation: 'DEFAULT',
	chatBackgroundColor: 'rgba(0, 0, 0, 0)',
	direction: 'top',
	fontFamily: 'inter',
	fontSize: 20,
	fontStyle: 'normal',
	fontWeight: 400,
	hideBots: false,
	hideCommands: false,
	messageHideTimeout: 0,
	messageShowDelay: 0,
	paddingContainer: 0,
	preset: 'clean',
	showAnnounceBadge: true,
	showBadges: true,
	showPlatformIcon: true,
	textShadowColor: 'rgba(0,0,0,1)',
	textShadowSize: 0,
	globalBadges: [],
	channelBadges: [],
}

const messages: Message[] = [
	{
		internalId: '1',
		id: '1',
		type: 'message',
		platform: 'twitch',
		chunks: [{ type: 'text', value: 'Cool' }],
		sender: 'twitchuser',
		senderColor: '#ff0000',
		senderDisplayName: 'TwitchUser',
		isItalic: false,
		createdAt: new Date(),
		isAnnounce: false,
	},
	{
		internalId: '2',
		id: '2',
		type: 'message',
		platform: 'kick',
		chunks: [{ type: 'text', value: 'Perfect' }],
		sender: 'kickuser',
		senderColor: '#00ff00',
		senderDisplayName: 'KickUser',
		kickBadges: [{ type: 'subscriber', text: 'Subscriber' }],
		isItalic: false,
		createdAt: new Date(),
		isAnnounce: false,
	},
	{
		internalId: '3',
		id: '3',
		type: 'message',
		platform: 'kick',
		chunks: [{ type: 'text', value: 'Hello' }],
		sender: 'moduser',
		senderColor: '#00ff00',
		senderDisplayName: 'ModUser',
		kickBadges: [{ type: 'moderator', text: 'Moderator' }],
		isItalic: false,
		createdAt: new Date(),
		isAnnounce: false,
	},
]
</script>

<template>
	<ChatBox
		:messages
		:settings="settings"
	/>
</template>
