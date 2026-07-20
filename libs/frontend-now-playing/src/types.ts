export interface Track {
	artist: string
	title: string
	imageUrl?: string | null
	progressMs?: number | null
	durationMs?: number | null
}

export const Preset = {
	TRANSPARENT: 'TRANSPARENT',
	AIDEN_REDESIGN: 'AIDEN_REDESIGN',
	SIMPLE_LINE: 'SIMPLE_LINE',
	PULSE_STRIP: 'PULSE_STRIP',
	AURA_STACK: 'AURA_STACK',
	VINYL_HAZE: 'VINYL_HAZE',
	SIGNAL_DECK: 'SIGNAL_DECK',
} as const

export interface Settings {
	id: string
	preset: keyof typeof Preset
	fontFamily: string
	fontWeight: number
	backgroundColor: string
	showImage: boolean
	hideTimeout?: number | null
}
