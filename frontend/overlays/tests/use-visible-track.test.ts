import { describe, expect, test } from 'bun:test'
import { effectScope, isReadonly, nextTick, ref } from 'vue'

import type { Track } from '@twir/frontend-now-playing'
import { isPlaybackRestart } from '@twir/frontend-now-playing/progress'

import { useVisibleTrack } from '../src/composables/now-playing/use-visible-track.ts'

interface ScheduledTimeout {
	id: number
	callback: () => void
	delay: number
}

function installFakeWindow() {
	const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
	const scheduled: ScheduledTimeout[] = []
	const cleared: number[] = []
	const active = new Map<number, () => void>()
	let nextId = 0

	Object.defineProperty(globalThis, 'window', {
		configurable: true,
		value: {
			setTimeout(callback: () => void, delay: number) {
				const id = ++nextId
				scheduled.push({ id, callback, delay })
				active.set(id, callback)
				return id
			},
			clearTimeout(id: number) {
				cleared.push(id)
				active.delete(id)
			},
		},
	})

	return {
		scheduled,
		cleared,
		run(id: number) {
			const callback = active.get(id)
			if (!callback) throw new Error(`Timeout ${id} is not active`)
			active.delete(id)
			callback()
		},
		restore() {
			if (originalWindow) {
				Object.defineProperty(globalThis, 'window', originalWindow)
			} else {
				Reflect.deleteProperty(globalThis, 'window')
			}
		},
	}
}

function createTrack(overrides: Partial<Track> = {}): Track {
	return {
		artist: 'Glass Animals',
		title: 'Heat Waves',
		imageUrl: 'cover-a',
		progressMs: 70_000,
		durationMs: 238_000,
		...overrides,
	}
}

const identityChanges: Array<{
	field: 'artist' | 'title' | 'imageUrl'
	firstValue: string
	secondValue: string
}> = [
	{ field: 'artist', firstValue: 'alt-J', secondValue: 'Foals' },
	{ field: 'title', firstValue: 'Gooey', secondValue: 'Youth' },
	{ field: 'imageUrl', firstValue: 'cover-b', secondValue: 'cover-c' },
]

describe('useVisibleTrack', () => {
	test('shows the initial track and schedules one valid positive timeout', () => {
		const browser = installFakeWindow()
		const scope = effectScope()

		try {
			const currentTrack = ref<Track | null>(createTrack())
			const result = scope.run(() => useVisibleTrack(currentTrack, ref(5_000)))

			expect(result).toBeDefined()
			expect(isReadonly(result?.visibleTrack)).toBe(true)
			expect(result?.visibleTrack.value).toEqual(currentTrack.value)
			expect(browser.scheduled).toHaveLength(1)
			expect(browser.scheduled[0]?.delay).toBe(5_000)
			expect(browser.cleared).toEqual([])
		} finally {
			scope.stop()
			browser.restore()
		}
	})

	test('updates new-object and in-place timing without restarting the timeout', async () => {
		const browser = installFakeWindow()
		const scope = effectScope()

		try {
			const currentTrack = ref<Track | null>(createTrack())
			const result = scope.run(() => useVisibleTrack(currentTrack, ref(5_000)))

			currentTrack.value = createTrack({ progressMs: 75_000 })
			await nextTick()
			expect(result?.visibleTrack.value?.progressMs).toBe(75_000)

			currentTrack.value.progressMs = 80_000
			currentTrack.value.durationMs = 300_000
			await nextTick()
			expect(result?.visibleTrack.value?.progressMs).toBe(80_000)
			expect(result?.visibleTrack.value?.durationMs).toBe(300_000)
			expect(browser.scheduled).toHaveLength(1)
			expect(browser.cleared).toEqual([])
		} finally {
			scope.stop()
			browser.restore()
		}
	})

	test('keeps a timed-out track hidden through later same-identity timing emissions', async () => {
		const browser = installFakeWindow()
		const scope = effectScope()

		try {
			const currentTrack = ref<Track | null>(createTrack())
			const result = scope.run(() => useVisibleTrack(currentTrack, ref(5_000)))
			const timeout = browser.scheduled[0]
			if (!timeout) throw new Error('Timeout was not scheduled')

			browser.run(timeout.id)
			expect(result?.visibleTrack.value).toBeNull()

			currentTrack.value = createTrack({ progressMs: 75_000 })
			await nextTick()
			currentTrack.value.progressMs = 80_000
			await nextTick()

			expect(result?.visibleTrack.value).toBeNull()
			expect(browser.scheduled).toHaveLength(1)
			expect(browser.cleared).toEqual([])
		} finally {
			scope.stop()
			browser.restore()
		}
	})

	test('reveals a hidden same-metadata playback restart and schedules a fresh timeout', async () => {
		const browser = installFakeWindow()
		const scope = effectScope()

		try {
			const previousTrack = createTrack({
				progressMs: 230_000,
				durationMs: 238_000,
			})
			const restartedTrack = createTrack({
				progressMs: 2_000,
				durationMs: 238_000,
			})
			const currentTrack = ref<Track | null>(previousTrack)
			const result = scope.run(() => useVisibleTrack(currentTrack, ref(5_000)))
			const firstTimeout = browser.scheduled[0]
			if (!firstTimeout) throw new Error('Timeout was not scheduled')
			browser.run(firstTimeout.id)
			expect(result?.visibleTrack.value).toBeNull()

			expect(isPlaybackRestart(previousTrack, restartedTrack)).toBe(true)
			currentTrack.value = restartedTrack
			await nextTick()

			expect(result?.visibleTrack.value).toEqual(currentTrack.value)
			expect(browser.scheduled).toHaveLength(2)
			expect(browser.scheduled[1]?.delay).toBe(5_000)
			expect(browser.cleared).toEqual([])
		} finally {
			scope.stop()
			browser.restore()
		}
	})

	test('keeps hidden same-metadata sub-threshold timing corrections in the same occurrence', async () => {
		const browser = installFakeWindow()
		const scope = effectScope()

		try {
			const currentTrack = ref<Track | null>(createTrack({ progressMs: 75_000 }))
			const result = scope.run(() => useVisibleTrack(currentTrack, ref(5_000)))
			const timeout = browser.scheduled[0]
			if (!timeout) throw new Error('Timeout was not scheduled')
			browser.run(timeout.id)

			currentTrack.value = createTrack({ progressMs: 72_000 })
			await nextTick()
			expect(result?.visibleTrack.value).toBeNull()

			currentTrack.value = createTrack({ progressMs: 64_000 })
			await nextTick()

			expect(result?.visibleTrack.value).toBeNull()
			expect(browser.scheduled).toHaveLength(1)
			expect(browser.cleared).toEqual([])
		} finally {
			scope.stop()
			browser.restore()
		}
	})

	for (const { field, firstValue, secondValue } of identityChanges) {
		test(`changing ${field} independently reveals the track and resets its timeout`, async () => {
			const browser = installFakeWindow()
			const scope = effectScope()

			try {
				const currentTrack = ref<Track | null>(createTrack())
				const result = scope.run(() => useVisibleTrack(currentTrack, ref(5_000)))
				const firstTimeout = browser.scheduled[0]
				if (!firstTimeout) throw new Error('Timeout was not scheduled')
				browser.run(firstTimeout.id)

				currentTrack.value = createTrack({ [field]: firstValue })
				await nextTick()
				const secondTimeout = browser.scheduled[1]
				if (!secondTimeout) throw new Error('Identity timeout was not scheduled')

				expect(result?.visibleTrack.value).toEqual(currentTrack.value)
				expect(browser.scheduled).toHaveLength(2)
				expect(secondTimeout.delay).toBe(5_000)

				currentTrack.value = createTrack({ [field]: secondValue })
				await nextTick()

				expect(result?.visibleTrack.value).toEqual(currentTrack.value)
				expect(browser.cleared).toEqual([secondTimeout.id])
				expect(browser.scheduled).toHaveLength(3)
				expect(browser.scheduled[2]?.delay).toBe(5_000)
			} finally {
				scope.stop()
				browser.restore()
			}
		})
	}

	test('never mutates track objects or the timeout ref across visibility transitions', async () => {
		const browser = installFakeWindow()
		const scope = effectScope()

		try {
			const currentTrack = ref<Track | null>(createTrack())
			const hideTimeout = ref<number | null>(1_000)
			const initialTrack = currentTrack.value
			if (!initialTrack) throw new Error('Initial track is missing')
			const initialTrackSnapshot = { ...initialTrack }
			const initialTimeoutSnapshot = hideTimeout.value
			const result = scope.run(() => useVisibleTrack(currentTrack, hideTimeout))
			const firstTimeout = browser.scheduled[0]
			if (!firstTimeout) throw new Error('Timeout was not scheduled')

			expect(currentTrack.value).toBe(initialTrack)
			expect(currentTrack.value).toEqual(initialTrackSnapshot)
			expect(hideTimeout.value).toBe(initialTimeoutSnapshot)

			const timingUpdate = createTrack({ progressMs: 75_000, durationMs: 300_000 })
			const timingUpdateSnapshot = { ...timingUpdate }
			currentTrack.value = timingUpdate
			const observedTimingUpdate = currentTrack.value
			await nextTick()

			expect(currentTrack.value).toBe(observedTimingUpdate)
			expect(currentTrack.value).toEqual(timingUpdateSnapshot)
			expect(hideTimeout.value).toBe(initialTimeoutSnapshot)

			currentTrack.value.progressMs = 80_000
			currentTrack.value.durationMs = 320_000
			const inPlaceTimingUpdate = currentTrack.value
			const inPlaceTimingSnapshot = { ...inPlaceTimingUpdate }
			await nextTick()

			expect(currentTrack.value).toBe(inPlaceTimingUpdate)
			expect(currentTrack.value).toEqual(inPlaceTimingSnapshot)
			expect(hideTimeout.value).toBe(initialTimeoutSnapshot)

			browser.run(firstTimeout.id)
			expect(result?.visibleTrack.value).toBeNull()
			expect(currentTrack.value).toBe(inPlaceTimingUpdate)
			expect(currentTrack.value).toEqual(inPlaceTimingSnapshot)
			expect(hideTimeout.value).toBe(initialTimeoutSnapshot)

			const identityUpdate = createTrack({ artist: 'alt-J', progressMs: 10_000 })
			const identityUpdateSnapshot = { ...identityUpdate }
			currentTrack.value = identityUpdate
			const observedIdentityUpdate = currentTrack.value
			await nextTick()

			expect(result?.visibleTrack.value).toEqual(identityUpdateSnapshot)
			expect(currentTrack.value).toBe(observedIdentityUpdate)
			expect(currentTrack.value).toEqual(identityUpdateSnapshot)
			expect(hideTimeout.value).toBe(initialTimeoutSnapshot)

			const trackBeforeTimeoutChange = currentTrack.value
			const trackBeforeTimeoutChangeSnapshot = { ...trackBeforeTimeoutChange }
			hideTimeout.value = 2_000
			const changedTimeoutSnapshot = hideTimeout.value
			await nextTick()

			expect(currentTrack.value).toBe(trackBeforeTimeoutChange)
			expect(currentTrack.value).toEqual(trackBeforeTimeoutChangeSnapshot)
			expect(hideTimeout.value).toBe(changedTimeoutSnapshot)
		} finally {
			scope.stop()
			browser.restore()
		}
	})

	test('timeout changes reveal and reschedule once while non-positive or null values do not schedule', async () => {
		const browser = installFakeWindow()
		const scope = effectScope()

		try {
			const currentTrack = ref<Track | null>(createTrack())
			const hideTimeout = ref<number | null>(1_000)
			const result = scope.run(() => useVisibleTrack(currentTrack, hideTimeout))
			const firstTimeout = browser.scheduled[0]
			if (!firstTimeout) throw new Error('Timeout was not scheduled')
			browser.run(firstTimeout.id)

			hideTimeout.value = 2_000
			await nextTick()
			const secondTimeout = browser.scheduled[1]
			if (!secondTimeout) throw new Error('Replacement timeout was not scheduled')
			expect(result?.visibleTrack.value).toEqual(currentTrack.value)
			expect(browser.scheduled).toHaveLength(2)
			expect(secondTimeout.delay).toBe(2_000)

			hideTimeout.value = 3_000
			await nextTick()
			const thirdTimeout = browser.scheduled[2]
			if (!thirdTimeout) throw new Error('Second replacement timeout was not scheduled')
			expect(browser.cleared).toEqual([secondTimeout.id])
			expect(browser.scheduled).toHaveLength(3)
			expect(thirdTimeout.delay).toBe(3_000)

			hideTimeout.value = 0
			await nextTick()
			expect(result?.visibleTrack.value).toEqual(currentTrack.value)
			expect(browser.cleared).toEqual([
				secondTimeout.id,
				thirdTimeout.id,
			])
			expect(browser.scheduled).toHaveLength(3)

			hideTimeout.value = -1
			await nextTick()
			hideTimeout.value = null
			await nextTick()
			expect(result?.visibleTrack.value).toEqual(currentTrack.value)
			expect(browser.scheduled).toHaveLength(3)
		} finally {
			scope.stop()
			browser.restore()
		}
	})

	test('clears the visible track and active timeout when the current track becomes null', async () => {
		const browser = installFakeWindow()
		const scope = effectScope()

		try {
			const currentTrack = ref<Track | null>(createTrack())
			const result = scope.run(() => useVisibleTrack(currentTrack, ref(5_000)))
			const timeout = browser.scheduled[0]
			if (!timeout) throw new Error('Timeout was not scheduled')

			currentTrack.value = null
			await nextTick()

			expect(result?.visibleTrack.value).toBeNull()
			expect(browser.cleared).toEqual([timeout.id])
			expect(browser.scheduled).toHaveLength(1)
		} finally {
			scope.stop()
			browser.restore()
		}
	})

	test('clears the active timeout when its effect scope is disposed', () => {
		const browser = installFakeWindow()
		const scope = effectScope()

		try {
			scope.run(() => useVisibleTrack(ref(createTrack()), ref(5_000)))
			const timeout = browser.scheduled[0]
			if (!timeout) throw new Error('Timeout was not scheduled')

			scope.stop()
			expect(browser.cleared).toEqual([timeout.id])
		} finally {
			scope.stop()
			browser.restore()
		}
	})

	test('runs without browser globals', () => {
		const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window')
		const scope = effectScope()

		try {
			Reflect.deleteProperty(globalThis, 'window')
			const currentTrack = ref<Track | null>(createTrack())
			const result = scope.run(() => useVisibleTrack(currentTrack, ref(5_000)))

			expect(result?.visibleTrack.value).toEqual(currentTrack.value)
		} finally {
			scope.stop()
			if (originalWindow) {
				Object.defineProperty(globalThis, 'window', originalWindow)
			} else {
				Reflect.deleteProperty(globalThis, 'window')
			}
		}
	})
})
