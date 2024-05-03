import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'

import type { EmotesStatisticsOpts } from '@/gql/graphql'

import { useEmotesStatisticQuery } from '@/api/emotes-statistic.js'

export const useCommunityEmotesStatisticTable = defineStore('features/community-emotes-statistic-table', () => {
	const opts = ref<EmotesStatisticsOpts>({})

	const { data } = useEmotesStatisticQuery(opts)

	const stats = computed(() => data?.value?.emotesStatistics)

	watch(stats, () => console.log(stats.value))

	return {
		stats
	}
})
