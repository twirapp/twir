import { defineStore, storeToRefs } from 'pinia'
import { computed, ref } from 'vue'

import type { EmotesStatisticEmoteDetailedOpts } from '@/gql/graphql'

import { useEmotesStatisticDetailsQuery } from '@/api/emotes-statistic'
import { usePagination } from '@/composables/use-pagination'
import { EmoteStatisticRange } from '@/gql/graphql'

export const useCommunityEmotesDetailsName = defineStore(
	'features/community-emotes-statistic-table/details-name',
	() => {
		const emoteName = ref<string>()

		function setEmoteName(name: string) {
			emoteName.value = name
		}

		return {
			emoteName,
			setEmoteName,
		}
	},
)

export const useCommunityEmotesDetails = defineStore(
	'features/community-emotes-statistic-table/details',
	() => {
		const { emoteName } = storeToRefs(useCommunityEmotesDetailsName())
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
	},
)
