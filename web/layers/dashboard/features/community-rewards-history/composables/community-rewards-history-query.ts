import { createGlobalState, debouncedRef } from '@vueuse/core'
import { computed, ref } from 'vue'

import type { TwitchRedemptionsOpts } from '@/gql/graphql.ts'

import { useProfile } from '@/api/auth'
import { useTwitchRewardsNew } from '@/api/twitch'
import { usePagination } from '@/composables/use-pagination.ts'

export const useCommunityRewardsHistoryQuery = createGlobalState(() => {
	const { data: profile } = useProfile()
	const { pagination, setPagination } = usePagination()

	const { data: existedRewards } = useTwitchRewardsNew()

	const searchInput = ref<string>()
	const searchInputDebounced = debouncedRef(searchInput, 200)

	const selectedRewards = ref<string[]>([])

	const query = computed<TwitchRedemptionsOpts>(() => ({
		byChannelId: profile.value?.selectedDashboardId,
		userSearch: searchInputDebounced.value,
		page: pagination.value.pageIndex,
		perPage: pagination.value.pageSize,
		rewardsIds: selectedRewards.value, // можно взять айдиники для селекта с rewards
	}))

	const rewardsOptions = computed(() => {
		return (
			existedRewards.value?.twitchRewards.map((reward) => ({
				id: reward.id,
				title: reward.title,
				image: reward.imageUrls?.at(0),
			})) ?? []
		)
	})

	function handleRewardFilter(id: string) {
		const index = selectedRewards.value.indexOf(id)
		if (index === -1) {
			selectedRewards.value.push(id)
		} else {
			selectedRewards.value.splice(index, 1)
		}
	}

	return {
		pagination,
		query,
		searchInput,
		handleRewardFilter,

		rewardsOptions,

		setPagination,
	}
})
