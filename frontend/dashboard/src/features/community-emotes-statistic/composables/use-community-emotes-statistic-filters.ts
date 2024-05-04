import { refDebounced } from '@vueuse/core'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import type { SortingState } from '@tanstack/vue-table'

import { EmoteStatisticRange, EmotesStatisticsOptsOrder } from '@/gql/graphql.js'

export const useCommunityEmotesStatisticFilters = defineStore('features/community-emotes-statistic-filters', () => {
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

	const emotesRange = ref(EmoteStatisticRange.LastDay)
	function changeEmoteRange(range: EmoteStatisticRange) {
		emotesRange.value = range
	}

	return {
		searchInput,
		debouncedSearchInput,

		sortingState,
		tableOrder,

		emotesRange,
		changeEmoteRange,
	}
})
