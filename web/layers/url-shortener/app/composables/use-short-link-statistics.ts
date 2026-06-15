import type { StatisticsPointDto } from '@twir/api/openapi'

import { ShortUrlGetStatisticsParamsIntervalEnum } from '@twir/api/openapi'

import { useOapi } from '~/composables/use-oapi'

export const useShortLinkStatistics = (shortId: MaybeRef<string>) => {
	const api = useOapi()
	const range = ref<'week' | 'month' | '3months'>('week')
	const isLoading = ref(false)
	const error = ref<Error | null>(null)
	const statistics = ref<StatisticsPointDto[]>([])

	const timeRange = computed(() => {
		const now = new Date()
		const fromDate = new Date()

		switch (range.value) {
			case 'week':
				fromDate.setDate(now.getDate() - 7)
				break
			case 'month':
				fromDate.setMonth(now.getMonth() - 1)
				break
			case '3months':
				fromDate.setMonth(now.getMonth() - 3)
				break
		}

		return {
			from: fromDate.getTime(),
			to: now.getTime(),
		}
	})

	const isDayRange = computed(() => range.value === 'week')

	const interval = computed(() =>
		isDayRange.value
			? ShortUrlGetStatisticsParamsIntervalEnum.Hour
			: ShortUrlGetStatisticsParamsIntervalEnum.Day
	)

	const fetchData = async () => {
		isLoading.value = true
		error.value = null
		// Clear old data immediately to prevent flickering
		statistics.value = []

		try {
			const response = await api.v1.shortUrlGetStatistics(unref(shortId), {
				from: timeRange.value.from,
				to: timeRange.value.to,
				interval: interval.value,
			})
			statistics.value = response.data.data ?? []
		} catch (err) {
			error.value = err as Error
			statistics.value = []
		} finally {
			isLoading.value = false
		}
	}

	// Watch for changes and refetch
	watch(
		[range, () => unref(shortId)],
		() => {
			fetchData()
		},
		{ immediate: true }
	)

	return {
		statistics: computed(() => statistics.value),
		isLoading: computed(() => isLoading.value),
		error: computed(() => error.value),
		range,
		isDayRange,
		refetch: fetchData,
	}
}
