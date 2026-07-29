import { useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'

export const headerStatsLayouts = ['rows', 'aggregate', 'grid'] as const

export type HeaderStatsLayout = (typeof headerStatsLayouts)[number]

export const headerStatsLayoutStorageKey = 'twirHeaderStatsLayoutV1'

export const defaultHeaderStatsLayout: HeaderStatsLayout = 'aggregate'

export function isHeaderStatsLayout(value: unknown): value is HeaderStatsLayout {
	return typeof value === 'string' && (headerStatsLayouts as readonly string[]).includes(value)
}

export function useHeaderLayout() {
	const stored = useLocalStorage<string>(headerStatsLayoutStorageKey, defaultHeaderStatsLayout)

	const layout = computed<HeaderStatsLayout>({
		get: () => (isHeaderStatsLayout(stored.value) ? stored.value : defaultHeaderStatsLayout),
		set: (value) => {
			stored.value = value
		},
	})

	return {
		layout,
	}
}
