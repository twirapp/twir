import { useEmotesStatisticDetailsQuery } from '#layers/dashboard/api/emotes-statistic'
import { createGlobalState } from '@vueuse/core'
import { computed, ref } from 'vue'

import type { EmotesStatisticEmoteDetailedOpts } from '~/gql/graphql'

import { usePagination } from '~/composables/use-pagination'
import { EmoteStatisticRange } from '~/gql/graphql'

export const useCommunityEmotesDetailsName = createGlobalState(() => {
	const emoteName = ref<string>()

	function setEmoteName(name: string) {
		emoteName.value = name
	}

	return {
		emoteName,
		setEmoteName,
	}
})

export const useCommunityEmotesDetails = createGlobalState(() => {
	const { emoteName } = useCommunityEmotesDetailsName()
	const { pagination: usagesPagination } = usePagination()
	const { pagination: topPagination } = usePagination()
	const range = ref(EmoteStatisticRange.LastDay)

	const opts = computed<EmotesStatisticEmoteDetailedOpts>(() => {
		return {
			emoteName: emoteName.value!,
			range: range.value,
			usagesByUsersPage: usagesPagination.value.pageIndex,
			usagesByUsersPerPage: usagesPagination.value.pageSize,
			topUsersPage: topPagination.value.pageIndex,
			topUsersPerPage: topPagination.value.pageSize,
		}
	})

	const { data: details, fetching: isLoading } = useEmotesStatisticDetailsQuery(opts)

	return {
		range,
		details,
		isLoading,
		usagesPagination,
		topPagination,
	}
})
