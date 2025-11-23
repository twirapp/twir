import type { BrbSettings as BrbOverlaySettings } from '@/composables/brb/use-brb-settings.js'

// emotes start
export type SevenTvEmote = {
	id: string
	name: string
	flags: number
	data: {
		animated: boolean
		host: {
			url: string
			files: Array<{ name: string; format: string; height: number; width: number }>
		}
	}
}

export type SevenTvChannelResponse = {
	user: {
		id: string
	}
	emote_set: {
		id: string
		emotes: Array<SevenTvEmote>
	}
}

export type SevenTvGlobalResponse = {
	emotes: Array<SevenTvEmote>
}

export type BttvEmote = {
	code: string
	imageType: string
	id: string
	animated: boolean
	height?: number
	width?: number
	modifier?: boolean
}

export type BttvChannelResponse = {
	channelEmotes: Array<BttvEmote>
	sharedEmotes: Array<BttvEmote>
}

export type BttvGlobalResponse = Array<BttvEmote>

export type FfzEmote = {
	name: string
	urls: Record<string, string>
	height: number
	width: number
	modifier: boolean
	modifier_flags?: number
}

export type FfzChannelResponse = {
	sets: {
		[x: string]: {
			emoticons: FfzEmote[]
		}
	}
}

export type FfzGlobalResponse = {
	sets: {
		[x: string]: {
			emoticons: FfzEmote[]
		}
	}
}
// emotes end

// brb start
export type BrbSetSettingsFn = (settings: BrbOverlaySettings) => void
export type BrbOnStartFn = (minutes: number, text: string) => void
export type BrbOnStopFn = () => void
// brb end

// tts start
export type TTSSayMessage = {
	text: string
	voice: string
	rate: string
	pitch: string
	volume: string
}
export type TTSOnSayFn = (message: TTSSayMessage) => void
export type TTSOnSkipFn = () => void
// tts end

export type KappagenTriggerRequestEmote = {
	id: string
	positions: string[]
}

// kappagen end

// dudes start
export type ChannelData = {
	channelDisplayName: string
	channelId: string
	channelName: string
}

export type UserData = {
	channelId: string
	userDisplayName: string
	userId: string
	userName: string
}
// dudes end
