import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import { EmoteStatisticRange } from '~/gql/graphql'

export const useTranslatedRanges = createGlobalState(() => {
	const { t } = useI18n()

	const ranges = computed(() => ({
		[EmoteStatisticRange.LastDay]: t('community.emotesStatistic.table.lastDay'),
		[EmoteStatisticRange.LastWeek]: t('community.emotesStatistic.table.lastWeek'),
		[EmoteStatisticRange.LastMonth]: t('community.emotesStatistic.table.lastMonth'),
		[EmoteStatisticRange.LastThreeMonth]: t('community.emotesStatistic.table.lastThreeMonth'),
		[EmoteStatisticRange.LastYear]: t('community.emotesStatistic.table.lastYear'),
	}))

	return {
		ranges,
	}
})
