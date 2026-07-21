export const SONG_REQUEST_OVERLAY_STYLES = [
	'CINEMA',
	'COMPACT',
	'TICKER',
	'STUDIO',
	'PORTRAIT',
	'PILL',
] as const

export type SongRequestOverlayStyle = (typeof SONG_REQUEST_OVERLAY_STYLES)[number]

export interface SongRequestOverlayProps {
	style?: SongRequestOverlayStyle | string | null
	title?: string | null
	requester?: string | null
	videoId?: string | null
	position?: number | null
	duration?: number | null
	isPlaying?: boolean | null
	accentColor?: string | null
	tickerBackgroundColor?: string | null
	tickerTextColor?: string | null
	tickerSpeed?: number | null
}

export interface ResolvedSongRequestOverlayProps {
	style: SongRequestOverlayStyle
	title: string
	requester: string
	videoId: string
	position: number
	duration: number
	isPlaying: boolean
	accentColor: string
	tickerBackgroundColor: string
	tickerTextColor: string
	tickerSpeed: number
}

export type SongRequestOverlayVisualProps = Omit<ResolvedSongRequestOverlayProps, 'style'>

export interface SongRequestPlaybackMetrics {
	position: number
	duration: number
	progress: number
	formattedPosition: string
	formattedDuration: string
}

export const SONG_REQUEST_OVERLAY_DEFAULTS = {
	style: 'CINEMA',
	title: '',
	requester: '',
	videoId: '',
	position: 0,
	duration: 0,
	isPlaying: false,
	accentColor: '#8B5CF6',
	tickerBackgroundColor: '#111827E6',
	tickerTextColor: '#FFFFFF',
	tickerSpeed: 35,
} satisfies ResolvedSongRequestOverlayProps

function normalizeSeconds(value: number | null | undefined): number {
	return Number.isFinite(value) ? Math.max(0, value ?? 0) : 0
}

export function normalizeSongRequestOverlayStyle(
	style: SongRequestOverlayProps['style']
): SongRequestOverlayStyle {
	return SONG_REQUEST_OVERLAY_STYLES.includes(style as SongRequestOverlayStyle)
		? (style as SongRequestOverlayStyle)
		: SONG_REQUEST_OVERLAY_DEFAULTS.style
}

export function resolveSongRequestOverlayProps(
	props: SongRequestOverlayProps
): ResolvedSongRequestOverlayProps {
	return {
		style: normalizeSongRequestOverlayStyle(props.style),
		title: props.title ?? SONG_REQUEST_OVERLAY_DEFAULTS.title,
		requester: props.requester ?? SONG_REQUEST_OVERLAY_DEFAULTS.requester,
		videoId: props.videoId ?? SONG_REQUEST_OVERLAY_DEFAULTS.videoId,
		position: normalizeSeconds(props.position),
		duration: normalizeSeconds(props.duration),
		isPlaying: props.isPlaying ?? SONG_REQUEST_OVERLAY_DEFAULTS.isPlaying,
		accentColor: props.accentColor ?? SONG_REQUEST_OVERLAY_DEFAULTS.accentColor,
		tickerBackgroundColor:
			props.tickerBackgroundColor ?? SONG_REQUEST_OVERLAY_DEFAULTS.tickerBackgroundColor,
		tickerTextColor: props.tickerTextColor ?? SONG_REQUEST_OVERLAY_DEFAULTS.tickerTextColor,
		tickerSpeed: Number.isFinite(props.tickerSpeed)
			? Math.min(100, Math.max(10, props.tickerSpeed ?? SONG_REQUEST_OVERLAY_DEFAULTS.tickerSpeed))
			: SONG_REQUEST_OVERLAY_DEFAULTS.tickerSpeed,
	}
}

export function getYouTubeThumbnailUrl(videoId: string | null | undefined): string {
	const normalizedVideoId = videoId?.trim()
	return normalizedVideoId
		? `https://i.ytimg.com/vi/${encodeURIComponent(normalizedVideoId)}/hqdefault.jpg`
		: ''
}

export function formatSongRequestTime(seconds: number | null | undefined): string {
	const totalSeconds = Math.floor(normalizeSeconds(seconds))
	const hours = Math.floor(totalSeconds / 3600)
	const minutes = Math.floor((totalSeconds % 3600) / 60)
	const remainingSeconds = totalSeconds % 60

	if (hours > 0) {
		return `${hours}:${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`
	}

	return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

export function getSongRequestPlaybackMetrics(
	position: number | null | undefined,
	duration: number | null | undefined
): SongRequestPlaybackMetrics {
	const safeDuration = normalizeSeconds(duration)
	const safePosition = normalizeSeconds(position)
	const clampedPosition = safeDuration > 0 ? Math.min(safePosition, safeDuration) : safePosition
	const progress =
		safeDuration > 0 ? Math.min(100, Math.max(0, (clampedPosition / safeDuration) * 100)) : 0

	return {
		position: clampedPosition,
		duration: safeDuration,
		progress,
		formattedPosition: formatSongRequestTime(clampedPosition),
		formattedDuration: formatSongRequestTime(safeDuration),
	}
}
