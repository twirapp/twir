import type { Font as InternalFont } from '@twir/grpc/generated/api/api/google_fonts_unprotected';
import type { Settings as ChatSettings } from '@twir/grpc/generated/api/api/overlays_chat';

export const enum EmoteFlag {
	Hidden = 0,
	Cursed,
	GrowX,
	NoSpace,
	FlipY,
	FlipX,
	Rotate90,
	Rotate270,
}

export const BttvOverlayEmotes = [
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
	type: 'text' | 'emote' | '3rd_party_emote';
	value: string;
	flags?: EmoteFlag[]
	zeroWidthModifiers?: string[]
	emoteWidth?: number
	emoteHeight?: number
}

export type Font = InternalFont

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
