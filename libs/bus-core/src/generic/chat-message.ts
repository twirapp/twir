export interface ChatMessageBadge {
	readonly id?: string
	readonly set_id: string
	readonly info?: string
	readonly text: string
}

export interface ChatMessageFragment {
	readonly text: string
	readonly position?: {
		readonly start?: number
		readonly end?: number
	}
	readonly type: number
}

export interface ChatMessage {
	readonly message?: {
		readonly text: string
		readonly fragments?: readonly ChatMessageFragment[]
	}
	readonly id?: string
	readonly broadcaster_user_id?: string
	readonly broadcaster_user_name?: string
	readonly broadcaster_user_login?: string
	readonly chatter_user_id?: string
	readonly chatter_user_name?: string
	readonly chatter_user_login?: string
	readonly message_type?: string
	readonly platform: string
	readonly channel_id: string
	readonly channel_binding_id: string
	readonly user_id: string
	readonly platform_channel_id: string
	readonly sender_id: string
	readonly sender_login: string
	readonly sender_display_name: string
	readonly message_id: string
	readonly text: string
	readonly badges?: readonly ChatMessageBadge[]
	readonly is_broadcaster: boolean
	readonly is_moderator: boolean
	readonly is_vip: boolean
	readonly is_subscriber: boolean
	readonly color: string
}
