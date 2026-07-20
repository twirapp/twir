import { describe, expect, test } from 'bun:test'
import { effectScope, nextTick, ref } from 'vue'

import { useTrackProgress } from '../src/composables/use-track-progress.ts'
import {
	PLAYBACK_RESTART_THRESHOLD_MS,
	formatTrackTime,
	isPlaybackRestart,
	normalizeTrackTiming,
} from '../src/utils/progress.ts'
import { Preset, type Track } from '../src/types.ts'

test('exposes timing fields and progress-aware presets', () => {
	const track: Track = {
		artist: 'Glass Animals',
		title: 'Heat Waves',
		progressMs: null,
		durationMs: null,
	}

	expect(track.progressMs).toBeNull()
	expect(track.durationMs).toBeNull()
	expect([
		Preset.PULSE_STRIP,
		Preset.AURA_STACK,
		Preset.VINYL_HAZE,
		Preset.SIGNAL_DECK,
	]).toEqual([
		'PULSE_STRIP',
		'AURA_STACK',
		'VINYL_HAZE',
		'SIGNAL_DECK',
	])
})

describe('normalizeTrackTiming', () => {
	test('interpolates timed progress and calculates its exact percentage', () => {
		const timing = normalizeTrackTiming({
			progressMs: 70_000,
			durationMs: 238_000,
		}, 5_000)

		expect(timing).toEqual({
			mode: 'timed',
			progressMs: 75_000,
			durationMs: 238_000,
			percent: 75_000 / 238_000 * 100,
		})
	})

	test('clamps progress above the duration', () => {
		expect(normalizeTrackTiming({
			progressMs: 237_000,
			durationMs: 238_000,
		}, 5_000)).toEqual({
			mode: 'timed',
			progressMs: 238_000,
			durationMs: 238_000,
			percent: 100,
		})
	})

	test('clamps negative progress to zero', () => {
		expect(normalizeTrackTiming({
			progressMs: -1_000,
			durationMs: 238_000,
		}, 0).progressMs).toBe(0)
	})

	test('uses ambient mode for null, missing, zero, or negative durations', () => {
		const invalidTracks = [
			{ progressMs: null, durationMs: null },
			{},
			{ progressMs: 0, durationMs: 0 },
			{ progressMs: 0, durationMs: -1 },
		]

		for (const track of invalidTracks) {
			expect(normalizeTrackTiming(track, 0)).toEqual({
				mode: 'ambient',
				progressMs: 0,
				durationMs: 0,
				percent: 0,
			})
		}
	})

	test('uses ambient mode for non-finite timing values', () => {
		expect(normalizeTrackTiming({
			progressMs: Number.POSITIVE_INFINITY,
			durationMs: 238_000,
		}, 0).mode).toBe('ambient')
		expect(normalizeTrackTiming({
			progressMs: 70_000,
			durationMs: Number.NaN,
		}, 0).mode).toBe('ambient')
		expect(normalizeTrackTiming({
			progressMs: 70_000,
			durationMs: 238_000,
		}, Number.NEGATIVE_INFINITY).mode).toBe('ambient')
	})
})

describe('formatTrackTime', () => {
	test('formats minute and hour durations', () => {
		expect(formatTrackTime(70_000)).toBe('01:10')
		expect(formatTrackTime(3_661_000)).toBe('1:01:01')
	})

	test('does not expose non-finite values', () => {
		expect(formatTrackTime(Number.POSITIVE_INFINITY)).toBe('00:00')
		expect(formatTrackTime(Number.NaN)).toBe('00:00')
	})
})

describe('isPlaybackRestart', () => {
	test('classifies backward jumps at the documented threshold', () => {
		expect(PLAYBACK_RESTART_THRESHOLD_MS).toBe(10_000)
		expect(isPlaybackRestart(
			{ progressMs: 75_000, durationMs: 238_000 },
			{ progressMs: 72_000, durationMs: 238_000 },
		)).toBe(false)
		expect(isPlaybackRestart(
			{ progressMs: 75_000, durationMs: 238_000 },
			{ progressMs: 65_001, durationMs: 238_000 },
		)).toBe(false)
		expect(isPlaybackRestart(
			{ progressMs: 75_000, durationMs: 238_000 },
			{ progressMs: 65_000, durationMs: 238_000 },
		)).toBe(true)
		expect(isPlaybackRestart(
			{ progressMs: 238_000, durationMs: 238_000 },
			{ progressMs: 0, durationMs: 238_000 },
		)).toBe(true)
		expect(isPlaybackRestart(
			{ progressMs: 120_000, durationMs: 238_000 },
			{ progressMs: 90_000, durationMs: 238_000 },
		)).toBe(true)
	})

	test('rejects ambient and invalid timing values', () => {
		expect(isPlaybackRestart(
			{ progressMs: null, durationMs: null },
			{ progressMs: 0, durationMs: 238_000 },
		)).toBe(false)
		expect(isPlaybackRestart(
			{ progressMs: 120_000, durationMs: 238_000 },
			{ progressMs: 0, durationMs: 0 },
		)).toBe(false)
		expect(isPlaybackRestart(
			{ progressMs: Number.POSITIVE_INFINITY, durationMs: 238_000 },
			{ progressMs: 0, durationMs: 238_000 },
		)).toBe(false)
	})
})

test('useTrackProgress initializes timed state and runs without a browser', async () => {
	const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		Reflect.deleteProperty(globalThis, 'window')
		expect('window' in globalThis).toBe(false)

		scope = effectScope()
		const track = ref<Track | null>({
			artist: 'Glass Animals',
			title: 'Heat Waves',
			progressMs: 70_000,
			durationMs: 238_000,
		})
		const progress = scope.run(() => useTrackProgress(track))

		expect(progress).toBeDefined()
		if (!progress) return

		expect(progress.timing.value).toEqual({
			mode: 'timed',
			progressMs: 70_000,
			durationMs: 238_000,
			percent: 70_000 / 238_000 * 100,
		})
		expect(progress.elapsedLabel.value).toBe('01:10')
		expect(progress.durationLabel.value).toBe('03:58')
		expect(progress.progressStyle.value).toEqual({
			'--track-progress': `${70_000 / 238_000 * 100}%`,
		})

		track.value = {
			artist: 'Ambient Artist',
			title: 'Ambient Track',
		}
		await nextTick()

		expect(progress.timing.value.mode).toBe('ambient')
		expect(progress.elapsedLabel.value).toBe('')
		expect(progress.durationLabel.value).toBe('')
	} finally {
		scope?.stop()
		if (originalWindow) {
			Object.defineProperty(globalThis, 'window', originalWindow)
		} else {
			Reflect.deleteProperty(globalThis, 'window')
		}
	}
})

test('useTrackProgress advances and clears its browser interval', () => {
	const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
	const originalPerformance = Object.getOwnPropertyDescriptor(globalThis, 'performance')
	const intervalId = 42
	let monotonicTime = 10_000
	let intervalCallback: (() => void) | undefined
	let intervalDelay: number | undefined
	let clearedIntervalId: number | undefined
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		Object.defineProperty(globalThis, 'window', {
			configurable: true,
			value: {
				setInterval(callback: () => void, delay: number) {
					intervalCallback = callback
					intervalDelay = delay
					return intervalId
				},
				clearInterval(id: number) {
					clearedIntervalId = id
				},
			},
		})
		Object.defineProperty(globalThis, 'performance', {
			configurable: true,
			value: { now: () => monotonicTime },
		})

		scope = effectScope()
		const track = ref<Track>({
			artist: 'Glass Animals',
			title: 'Heat Waves',
			progressMs: 70_000,
			durationMs: 238_000,
		})
		const progress = scope.run(() => useTrackProgress(track))

		expect(intervalDelay).toBe(1000)
		expect(progress?.timing.value.progressMs).toBe(70_000)
		if (!intervalCallback) {
			throw new Error('window.setInterval callback was not captured')
		}

		monotonicTime += 2_000
		intervalCallback()
		expect(progress?.timing.value.progressMs).toBe(72_000)

		scope.stop()
		expect(clearedIntervalId).toBe(intervalId)
	} finally {
		scope?.stop()
		if (originalPerformance) {
			Object.defineProperty(globalThis, 'performance', originalPerformance)
		} else {
			Reflect.deleteProperty(globalThis, 'performance')
		}
		if (originalWindow) {
			Object.defineProperty(globalThis, 'window', originalWindow)
		} else {
			Reflect.deleteProperty(globalThis, 'window')
		}
	}
})

test('useTrackProgress does not rewind when the wall clock rolls back', () => {
	const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
	const originalPerformance = Object.getOwnPropertyDescriptor(globalThis, 'performance')
	const originalDateNow = Date.now
	let wallClock = 50_000
	let monotonicTime = 10_000
	let intervalCallback: (() => void) | undefined
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		Object.defineProperty(globalThis, 'window', {
			configurable: true,
			value: {
				setInterval(callback: () => void) {
					intervalCallback = callback
					return 42
				},
				clearInterval() {},
			},
		})
		Object.defineProperty(globalThis, 'performance', {
			configurable: true,
			value: { now: () => monotonicTime },
		})
		Date.now = () => wallClock

		scope = effectScope()
		const track = ref<Track>({
			artist: 'Glass Animals',
			title: 'Heat Waves',
			progressMs: 70_000,
			durationMs: 238_000,
		})
		const progress = scope.run(() => useTrackProgress(track))

		if (!intervalCallback) {
			throw new Error('window.setInterval callback was not captured')
		}

		monotonicTime += 2_000
		wallClock -= 10_000
		intervalCallback()

		expect(progress?.timing.value.progressMs).toBe(72_000)
	} finally {
		scope?.stop()
		Date.now = originalDateNow
		if (originalPerformance) {
			Object.defineProperty(globalThis, 'performance', originalPerformance)
		} else {
			Reflect.deleteProperty(globalThis, 'performance')
		}
		if (originalWindow) {
			Object.defineProperty(globalThis, 'window', originalWindow)
		} else {
			Reflect.deleteProperty(globalThis, 'window')
		}
	}
})

test('useTrackProgress preserves its monotonic baseline for in-place timing updates', async () => {
	const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
	const originalPerformance = Object.getOwnPropertyDescriptor(globalThis, 'performance')
	let monotonicTime = 10_000
	let intervalCallback: (() => void) | undefined
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		Object.defineProperty(globalThis, 'window', {
			configurable: true,
			value: {
				setInterval(callback: () => void) {
					intervalCallback = callback
					return 42
				},
				clearInterval() {},
			},
		})
		Object.defineProperty(globalThis, 'performance', {
			configurable: true,
			value: { now: () => monotonicTime },
		})

		scope = effectScope()
		const track = ref<Track>({
			artist: 'Glass Animals',
			title: 'Heat Waves',
			progressMs: 70_000,
			durationMs: 238_000,
		})
		const progress = scope.run(() => useTrackProgress(track))

		if (!intervalCallback) {
			throw new Error('window.setInterval callback was not captured')
		}
		const tick = intervalCallback

		monotonicTime += 5_000
		tick()
		expect(progress?.timing.value.progressMs).toBe(75_000)

		track.value.progressMs = 72_000
		await nextTick()
		expect(progress?.timing.value.progressMs).toBe(75_000)

		monotonicTime += 2_000
		tick()
		expect(progress?.timing.value.progressMs).toBe(77_000)

		track.value.durationMs = 300_000
		await nextTick()
		expect(progress?.timing.value.progressMs).toBe(77_000)
	} finally {
		scope?.stop()
		if (originalPerformance) {
			Object.defineProperty(globalThis, 'performance', originalPerformance)
		} else {
			Reflect.deleteProperty(globalThis, 'performance')
		}
		if (originalWindow) {
			Object.defineProperty(globalThis, 'window', originalWindow)
		} else {
			Reflect.deleteProperty(globalThis, 'window')
		}
	}
})

test('useTrackProgress does not rewind a same-track correction', async () => {
	const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
	const originalPerformance = Object.getOwnPropertyDescriptor(globalThis, 'performance')
	let monotonicTime = 10_000
	let intervalCallback: (() => void) | undefined
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		Object.defineProperty(globalThis, 'window', {
			configurable: true,
			value: {
				setInterval(callback: () => void) {
					intervalCallback = callback
					return 42
				},
				clearInterval() {},
			},
		})
		Object.defineProperty(globalThis, 'performance', {
			configurable: true,
			value: { now: () => monotonicTime },
		})

		scope = effectScope()
		const track = ref<Track>({
			artist: 'Glass Animals',
			title: 'Heat Waves',
			imageUrl: 'cover',
			progressMs: 70_000,
			durationMs: 238_000,
		})
		const progress = scope.run(() => useTrackProgress(track))

		if (!intervalCallback) {
			throw new Error('window.setInterval callback was not captured')
		}
		const tick = intervalCallback

		monotonicTime += 5_000
		tick()
		expect(progress?.timing.value.progressMs).toBe(75_000)

		track.value = {
			artist: 'Glass Animals',
			title: 'Heat Waves',
			imageUrl: 'cover',
			progressMs: 72_000,
			durationMs: 238_000,
		}
		await nextTick()
		expect(progress?.timing.value.progressMs).toBe(75_000)
		expect(track.value.progressMs).toBe(72_000)

		monotonicTime += 2_000
		tick()
		expect(progress?.timing.value.progressMs).toBe(77_000)
	} finally {
		scope?.stop()
		if (originalPerformance) {
			Object.defineProperty(globalThis, 'performance', originalPerformance)
		} else {
			Reflect.deleteProperty(globalThis, 'performance')
		}
		if (originalWindow) {
			Object.defineProperty(globalThis, 'window', originalWindow)
		} else {
			Reflect.deleteProperty(globalThis, 'window')
		}
	}
})

test('useTrackProgress resets a significant same-metadata rollback', async () => {
	const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
	const originalPerformance = Object.getOwnPropertyDescriptor(globalThis, 'performance')
	let monotonicTime = 10_000
	let intervalCallback: (() => void) | undefined
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		Object.defineProperty(globalThis, 'window', {
			configurable: true,
			value: {
				setInterval(callback: () => void) {
					intervalCallback = callback
					return 42
				},
				clearInterval() {},
			},
		})
		Object.defineProperty(globalThis, 'performance', {
			configurable: true,
			value: { now: () => monotonicTime },
		})

		scope = effectScope()
		const track = ref<Track>({
			artist: 'Glass Animals',
			title: 'Heat Waves',
			imageUrl: 'cover',
			progressMs: 110_000,
			durationMs: 238_000,
		})
		const progress = scope.run(() => useTrackProgress(track))

		if (!intervalCallback) {
			throw new Error('window.setInterval callback was not captured')
		}
		const tick = intervalCallback

		monotonicTime += 10_000
		tick()
		expect(progress?.timing.value.progressMs).toBe(120_000)

		track.value = {
			artist: 'Glass Animals',
			title: 'Heat Waves',
			imageUrl: 'cover',
			progressMs: 90_000,
			durationMs: 238_000,
		}
		await nextTick()
		expect(progress?.timing.value.progressMs).toBe(90_000)

		monotonicTime += 2_000
		tick()
		expect(progress?.timing.value.progressMs).toBe(92_000)
	} finally {
		scope?.stop()
		if (originalPerformance) {
			Object.defineProperty(globalThis, 'performance', originalPerformance)
		} else {
			Reflect.deleteProperty(globalThis, 'performance')
		}
		if (originalWindow) {
			Object.defineProperty(globalThis, 'window', originalWindow)
		} else {
			Reflect.deleteProperty(globalThis, 'window')
		}
	}
})

test('useTrackProgress resets lower progress for a different track identity', async () => {
	const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
	const originalPerformance = Object.getOwnPropertyDescriptor(globalThis, 'performance')
	let monotonicTime = 10_000
	let intervalCallback: (() => void) | undefined
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		Object.defineProperty(globalThis, 'window', {
			configurable: true,
			value: {
				setInterval(callback: () => void) {
					intervalCallback = callback
					return 42
				},
				clearInterval() {},
			},
		})
		Object.defineProperty(globalThis, 'performance', {
			configurable: true,
			value: { now: () => monotonicTime },
		})

		scope = effectScope()
		const track = ref<Track>({
			artist: 'Glass Animals',
			title: 'Heat Waves',
			imageUrl: 'cover-a',
			progressMs: 70_000,
			durationMs: 238_000,
		})
		const progress = scope.run(() => useTrackProgress(track))

		if (!intervalCallback) {
			throw new Error('window.setInterval callback was not captured')
		}
		const tick = intervalCallback

		monotonicTime += 5_000
		tick()
		expect(progress?.timing.value.progressMs).toBe(75_000)

		track.value = {
			artist: 'Glass Animals',
			title: 'Gooey',
			imageUrl: 'cover-a',
			progressMs: 10_000,
			durationMs: 200_000,
		}
		await nextTick()
		expect(progress?.timing.value.progressMs).toBe(10_000)

		monotonicTime += 1_000
		tick()
		track.value = {
			artist: 'Different Artist',
			title: 'Gooey',
			imageUrl: 'cover-a',
			progressMs: 5_000,
			durationMs: 200_000,
		}
		await nextTick()
		expect(progress?.timing.value.progressMs).toBe(5_000)

		monotonicTime += 1_000
		tick()
		track.value = {
			artist: 'Different Artist',
			title: 'Gooey',
			imageUrl: 'cover-b',
			progressMs: 2_000,
			durationMs: 200_000,
		}
		await nextTick()
		expect(progress?.timing.value.progressMs).toBe(2_000)
	} finally {
		scope?.stop()
		if (originalPerformance) {
			Object.defineProperty(globalThis, 'performance', originalPerformance)
		} else {
			Reflect.deleteProperty(globalThis, 'performance')
		}
		if (originalWindow) {
			Object.defineProperty(globalThis, 'window', originalWindow)
		} else {
			Reflect.deleteProperty(globalThis, 'window')
		}
	}
})

test('useTrackProgress only runs a browser interval for timed tracks', async () => {
	const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
	const originalPerformance = Object.getOwnPropertyDescriptor(globalThis, 'performance')
	const intervalId = 42
	const intervalDelays: number[] = []
	const clearedIntervalIds: number[] = []
	let monotonicTime = 10_000
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		Object.defineProperty(globalThis, 'window', {
			configurable: true,
			value: {
				setInterval(_callback: () => void, delay: number) {
					intervalDelays.push(delay)
					return intervalId
				},
				clearInterval(id: number) {
					clearedIntervalIds.push(id)
				},
			},
		})
		Object.defineProperty(globalThis, 'performance', {
			configurable: true,
			value: { now: () => monotonicTime },
		})

		scope = effectScope()
		const track = ref<Track>({
			artist: 'Ambient Artist',
			title: 'Ambient Track',
			progressMs: null,
			durationMs: null,
		})
		const progress = scope.run(() => useTrackProgress(track))

		expect(progress?.timing.value.mode).toBe('ambient')
		expect(intervalDelays).toEqual([])

		track.value.progressMs = 70_000
		track.value.durationMs = 238_000
		await nextTick()
		expect(progress?.timing.value.mode).toBe('timed')
		expect(intervalDelays).toEqual([1000])

		monotonicTime += 1_000
		track.value.progressMs = 71_000
		await nextTick()
		expect(intervalDelays).toEqual([1000])

		track.value.progressMs = null
		track.value.durationMs = null
		await nextTick()
		expect(progress?.timing.value.mode).toBe('ambient')
		expect(intervalDelays).toEqual([1000])
		expect(clearedIntervalIds).toEqual([intervalId])

		scope.stop()
		scope.stop()
		expect(clearedIntervalIds).toEqual([intervalId])
	} finally {
		scope?.stop()
		if (originalPerformance) {
			Object.defineProperty(globalThis, 'performance', originalPerformance)
		} else {
			Reflect.deleteProperty(globalThis, 'performance')
		}
		if (originalWindow) {
			Object.defineProperty(globalThis, 'window', originalWindow)
		} else {
			Reflect.deleteProperty(globalThis, 'window')
		}
	}
})
