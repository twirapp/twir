import type { ChatMessage, ChatMessageBadge } from '@twir/bus-core'

export interface ChannelBinding {
	readonly id: string
	readonly channelId: string
	readonly platformChannelId: string
	readonly userId: string
	readonly enabled: boolean
}

export interface YoutubeTextChatMessage {
	readonly id: string
	readonly text: string
	readonly author: {
		readonly id: string
		readonly name: string
		readonly isModerator: boolean
		readonly badges: readonly ChatMessageBadge[]
	}
}

export function normalizeYoutubeTextMessage(
	binding: ChannelBinding,
	broadcasterName: string,
	message: YoutubeTextChatMessage
): ChatMessage {
	return {
		message: {
			text: message.text,
			fragments: [{
				text: message.text,
				position: { start: 0, end: Array.from(message.text).length },
				type: 0,
			}],
		},
		id: message.id,
		broadcaster_user_id: binding.platformChannelId,
		broadcaster_user_name: broadcasterName,
		broadcaster_user_login: broadcasterName,
		chatter_user_id: message.author.id,
		chatter_user_name: message.author.name,
		chatter_user_login: message.author.name,
		message_type: 'text',
		platform: 'youtube',
		channel_id: binding.channelId,
		channel_binding_id: binding.id,
		user_id: binding.userId,
		platform_channel_id: binding.platformChannelId,
		sender_id: message.author.id,
		sender_login: message.author.name,
		sender_display_name: message.author.name,
		message_id: message.id,
		text: message.text,
		badges: message.author.badges,
		is_broadcaster: message.author.id === binding.platformChannelId,
		is_moderator: message.author.isModerator,
		is_vip: false,
		is_subscriber: false,
		color: '',
	}
}
