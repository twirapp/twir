export enum EmoteFlag {
	Hidden = 2,
	FlipX = 3,
	FlipY = 5,
	GrowX = 9,
	Rainbow = 2049,
	HyperRed = 4097,
	HyperShake = 8193,
	Cursed = 16385,
	Jam = 32769,
	Bounce = 65537,
}

export const BttvZeroModifiers = [
	'cvMask',
	'cvHazmat',
	'SoSnowy',
	'IceCold',
	'TopHat',
	'SantaHat',
	'ReinDeer',
	'CandyCane',
]

export interface MessageChunk {
	type: 'text' | 'emote' | '3rd_party_emote' | 'emoji'
	value: string
	flags?: EmoteFlag[]
	modifier_flags?: number
	zeroWidthModifiers?: string[]
	emoteWidth?: number
	emoteHeight?: number
	emoteName?: string
}

export interface BadgeVersion {
	id: string
	image_url_1x: string
	image_url_2x: string
	image_url_4x: string
}

export interface ChatBadge {
	set_id: string
	versions: Array<BadgeVersion>
}

export interface Message {
	internalId: string
	id?: string
	type: 'message' | 'system'
	chunks: MessageChunk[]
	sender?: string
	senderColor?: string
	senderDisplayName?: string
	badges?: Record<string, string>
	isItalic: boolean
	createdAt: Date
	announceColor?: string
	isAnnounce: boolean
}

export interface Settings {
	channelId: string
	channelName: string
	channelDisplayName: string
	animation: string
	chatBackgroundColor: string
	direction: string
	fontFamily: string
	fontSize: number
	fontStyle: string
	fontWeight: number
	hideBots: boolean
	hideCommands: boolean
	messageHideTimeout: number
	messageShowDelay: number
	paddingContainer: number
	preset: string
	showAnnounceBadge: boolean
	showBadges: boolean
	textShadowColor: string
	textShadowSize: number
	globalBadges: ChatBadge[]
	channelBadges: ChatBadge[]
}

export interface MessageComponentProps {
	msg: Message
	settings: Settings
	userColor: string
}
