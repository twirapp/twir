import type { Settings as ChatSettings } from '@twir/grpc/generated/api/api/overlays_chat';

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
];

export type MessageChunk = {
	type: 'text' | 'emote' | '3rd_party_emote' | 'emoji';
	value: string;
	flags?: EmoteFlag[]
	modifier_flags?: number
	zeroWidthModifiers?: string[]
	emoteWidth?: number
	emoteHeight?: number
	emoteName?: string
}

export type BadgeVersion = {
	id: string,
	image_url_1x: string,
	image_url_2x: string,
	image_url_4x: string,
}

export type ChatBadge = {
	set_id: string,
	versions: Array<BadgeVersion>
}

export type Message = {
	internalId: string,
	id?: string,
	type: 'message' | 'system',
	chunks: MessageChunk[],
	sender?: string,
	senderColor?: string,
	senderDisplayName?: string
	badges?: Record<string, string>,
	isItalic: boolean;
	createdAt: Date;
	announceColor?: string;
	isAnnounce: boolean;
};

export type Settings = {
	channelId: string
	channelName: string
	channelDisplayName: string
	globalBadges: Map<string, ChatBadge>
	channelBadges: Map<string, BadgeVersion>
} & ChatSettings;


export type MessageComponentProps = {
	msg: Message,
	settings: Settings,
	userColor: string,
}
