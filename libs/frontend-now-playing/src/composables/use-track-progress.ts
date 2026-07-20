import { type Ref, computed, onScopeDispose, ref, watch } from 'vue'

import type { Track } from '../types.js'
import {
	type TrackTimingInput,
	formatTrackTime,
	isPlaybackRestart,
	normalizeTrackTiming,
} from '../utils/progress.js'

interface TrackIdentity {
	artist: string
	title: string
	imageUrl: string | null
}

function getMonotonicTime(): number {
	return typeof performance === 'undefined' ? Date.now() : performance.now()
}

function getTrackIdentity(track: Track | null | undefined): TrackIdentity | undefined {
	if (!track) return undefined

	return {
		artist: track.artist,
		title: track.title,
		imageUrl: track.imageUrl ?? null,
	}
}

function hasSameIdentity(
	current: TrackIdentity | undefined,
	next: TrackIdentity | undefined,
): boolean {
	return current !== undefined
		&& next !== undefined
		&& current.artist === next.artist
		&& current.title === next.title
		&& current.imageUrl === next.imageUrl
}

export function useTrackProgress(track: Ref<Track | null | undefined>) {
	const receivedAt = ref(getMonotonicTime())
	const clock = ref(receivedAt.value)
	const baselineTiming = ref<TrackTimingInput>({})
	const currentIdentity = ref<TrackIdentity>()
	const timerWindow = typeof window === 'undefined' ? undefined : window
	let interval: number | undefined

	const stopInterval = () => {
		if (interval === undefined || timerWindow === undefined) return

		timerWindow.clearInterval(interval)
		interval = undefined
	}

	const syncInterval = () => {
		const mode = normalizeTrackTiming(baselineTiming.value, 0).mode
		if (mode !== 'timed' || timerWindow === undefined) {
			stopInterval()
			return
		}
		if (interval !== undefined) return

		interval = timerWindow.setInterval(() => {
			clock.value = getMonotonicTime()
		}, 1000)
	}

	watch(
		() => [
			track.value?.artist,
			track.value?.title,
			track.value?.imageUrl,
			track.value?.progressMs,
			track.value?.durationMs,
		] as const,
		() => {
			const now = getMonotonicTime()
			const currentTiming = normalizeTrackTiming(
				baselineTiming.value,
				now - receivedAt.value,
			)
			const incomingTiming: TrackTimingInput = {
				progressMs: track.value?.progressMs,
				durationMs: track.value?.durationMs,
			}
			const normalizedIncoming = normalizeTrackTiming(incomingTiming, 0)
			const nextIdentity = getTrackIdentity(track.value)

			if (
				hasSameIdentity(currentIdentity.value, nextIdentity)
				&& currentTiming.mode === 'timed'
				&& normalizedIncoming.mode === 'timed'
				&& !isPlaybackRestart(currentTiming, normalizedIncoming)
			) {
				baselineTiming.value = {
					progressMs: Math.max(currentTiming.progressMs, normalizedIncoming.progressMs),
					durationMs: normalizedIncoming.durationMs,
				}
			} else {
				baselineTiming.value = incomingTiming
			}

			currentIdentity.value = nextIdentity
			receivedAt.value = now
			clock.value = now
			syncInterval()
		},
		{ immediate: true },
	)

	onScopeDispose(stopInterval)

	const timing = computed(() => normalizeTrackTiming(
		baselineTiming.value,
		clock.value - receivedAt.value,
	))
	const elapsedLabel = computed(() => (
		timing.value.mode === 'timed' ? formatTrackTime(timing.value.progressMs) : ''
	))
	const durationLabel = computed(() => (
		timing.value.mode === 'timed' ? formatTrackTime(timing.value.durationMs) : ''
	))
	const progressStyle = computed(() => ({
		'--track-progress': `${timing.value.percent}%`,
	}))

	return {
		timing,
		elapsedLabel,
		durationLabel,
		progressStyle,
	}
}
