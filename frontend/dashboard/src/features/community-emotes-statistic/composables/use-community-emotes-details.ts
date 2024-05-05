import { defineStore, storeToRefs } from 'pinia'
import { computed, ref, watch } from 'vue'

import type {
	EmotesStatisticEmoteDetailedOpts,
	EmotesStatisticsDetailsQuery,
} from '@/gql/graphql'

import { useEmotesStatisticDetailsQuery } from '@/api/emotes-statistic'
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
		const range = ref(EmoteStatisticRange.LastDay)

		const opts = computed<EmotesStatisticEmoteDetailedOpts>(() => {
			return {
				emoteName: emoteName.value!,
				range: range.value,
			}
		})

		const details = ref<EmotesStatisticsDetailsQuery>()

		const { data } = useEmotesStatisticDetailsQuery(opts)

		watch(data, (newData) => {
			details.value = newData
		})

		return {
			range,
			details,
		}
	},
)
