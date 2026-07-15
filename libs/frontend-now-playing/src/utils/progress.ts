export type ProgressMode = 'timed' | 'ambient'

/** Minimum backward jump treated as a new playback occurrence. */
export const PLAYBACK_RESTART_THRESHOLD_MS = 10_000

export interface TrackTimingInput {
	progressMs?: number | null
	durationMs?: number | null
}

export interface NormalizedTrackTiming {
	mode: ProgressMode
	progressMs: number
	durationMs: number
	percent: number
}

export function normalizeTrackTiming(
	track: TrackTimingInput,
	elapsedMs: number,
): NormalizedTrackTiming {
	const { progressMs, durationMs } = track
	if (
		typeof progressMs !== 'number'
		|| !Number.isFinite(progressMs)
		|| typeof durationMs !== 'number'
		|| !Number.isFinite(durationMs)
		|| durationMs <= 0
		|| !Number.isFinite(elapsedMs)
	) {
		return {
			mode: 'ambient',
			progressMs: 0,
			durationMs: 0,
			percent: 0,
		}
	}

	const normalizedProgress = Math.min(Math.max(progressMs + elapsedMs, 0), durationMs)

	return {
		mode: 'timed',
		progressMs: normalizedProgress,
		durationMs,
		percent: normalizedProgress / durationMs * 100,
	}
}

export function isPlaybackRestart(
	previous: TrackTimingInput,
	current: TrackTimingInput,
): boolean {
	const previousTiming = normalizeTrackTiming(previous, 0)
	const currentTiming = normalizeTrackTiming(current, 0)
	if (previousTiming.mode !== 'timed' || currentTiming.mode !== 'timed') return false

	return previousTiming.progressMs - currentTiming.progressMs >= PLAYBACK_RESTART_THRESHOLD_MS
}

export function formatTrackTime(ms: number): string {
	const safeMs = Number.isFinite(ms) ? Math.max(ms, 0) : 0
	const totalSeconds = Math.floor(safeMs / 1000)
	const seconds = totalSeconds % 60
	const totalMinutes = Math.floor(totalSeconds / 60)
	const minutes = totalMinutes % 60
	const hours = Math.floor(totalMinutes / 60)

	if (hours > 0) {
		return `${hours}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
	}

	return `${String(totalMinutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
}
