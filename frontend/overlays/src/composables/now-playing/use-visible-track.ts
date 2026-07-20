import { type Ref, computed, onScopeDispose, ref, watch } from 'vue'

import type { Track } from '@twir/frontend-now-playing'
import { isPlaybackRestart } from '@twir/frontend-now-playing/progress'

export function useVisibleTrack(
	currentTrack: Readonly<Ref<Track | null | undefined>>,
	hideTimeout: Readonly<Ref<number | null | undefined>>,
) {
	const isVisible = ref(false)
	let timerId: number | undefined
	let previousTrack: Track | null | undefined
	let previousTimeout: number | null | undefined
	let initialized = false

	function clearTimer() {
		if (timerId === undefined) return
		if (typeof window !== 'undefined') window.clearTimeout(timerId)
		timerId = undefined
	}

	watch([
		() => currentTrack.value?.artist,
		() => currentTrack.value?.title,
		() => currentTrack.value?.imageUrl,
		() => currentTrack.value?.progressMs,
		() => currentTrack.value?.durationMs,
		() => hideTimeout.value,
	], () => {
		const track = currentTrack.value
		const metadataChanged = !initialized
			|| previousTrack?.artist !== track?.artist
			|| previousTrack?.title !== track?.title
			|| previousTrack?.imageUrl !== track?.imageUrl
		const playbackRestarted = previousTrack != null
			&& track != null
			&& isPlaybackRestart(previousTrack, track)
		const timeoutChanged = initialized && !Object.is(previousTimeout, hideTimeout.value)

		previousTrack = track == null ? track : { ...track }
		previousTimeout = hideTimeout.value
		initialized = true

		if (!metadataChanged && !playbackRestarted && !timeoutChanged) return

		clearTimer()
		isVisible.value = track != null

		const timeout = hideTimeout.value
		if (
			!isVisible.value
			|| typeof window === 'undefined'
			|| typeof timeout !== 'number'
			|| !Number.isFinite(timeout)
			|| timeout <= 0
		) {
			return
		}

		timerId = window.setTimeout(() => {
			timerId = undefined
			isVisible.value = false
		}, timeout)
	}, { immediate: true })

	const visibleTrack = computed(() => (
		isVisible.value ? currentTrack.value ?? null : null
	))

	onScopeDispose(clearTimer)

	return { visibleTrack }
}
