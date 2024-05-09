import { createGlobalState, refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'

import type { SortingState } from '@tanstack/vue-table'

import { EmoteStatisticRange, EmotesStatisticsOptsOrder } from '@/gql/graphql.js'

export const useCommunityEmotesStatisticFilters = createGlobalState(() => {
	const searchInput = ref('')
	const debouncedSearchInput = refDebounced(searchInput, 500)

	const sortingState = ref<SortingState>([
		{
			desc: true,
			id: 'usages', // accessorKey
		},
	])

	const tableOrder = computed(() => {
		return sortingState.value[0].desc
			? EmotesStatisticsOptsOrder.Desc
			: EmotesStatisticsOptsOrder.Asc
	})

	const tableRange = ref(EmoteStatisticRange.LastDay)
	function changeTableRange(range: EmoteStatisticRange) {
		tableRange.value = range
	}

	return {
		searchInput,
		debouncedSearchInput,

		sortingState,
		tableOrder,

		tableRange,
		changeTableRange,
	}
})
