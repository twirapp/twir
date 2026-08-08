import { nextTick } from 'vue'
import { beforeEach, describe, expect, it } from 'vitest'

import {
	defaultHeaderStatsLayout,
	headerStatsLayoutStorageKey,
	useHeaderLayout,
} from './use-header-layout.js'

describe('useHeaderLayout', () => {
	beforeEach(() => {
		window.localStorage.clear()
	})

	it('defaults to aggregate when localStorage is empty', () => {
		const { layout } = useHeaderLayout()

		expect(layout.value).toBe('aggregate')
		expect(layout.value).toBe(defaultHeaderStatsLayout)
	})

	it('persists the selected layout to localStorage', async () => {
		const { layout } = useHeaderLayout()

		layout.value = 'grid'
		await nextTick()

		expect(window.localStorage.getItem(headerStatsLayoutStorageKey)).toBe('grid')

		const next = useHeaderLayout()
		expect(next.layout.value).toBe('grid')
	})

	it('reads a valid stored layout', () => {
		window.localStorage.setItem(headerStatsLayoutStorageKey, 'rows')

		const { layout } = useHeaderLayout()

		expect(layout.value).toBe('rows')
	})

	it('ignores an invalid stored value and falls back to aggregate', () => {
		window.localStorage.setItem(headerStatsLayoutStorageKey, 'banana')

		const { layout } = useHeaderLayout()

		expect(layout.value).toBe('aggregate')
	})
})
