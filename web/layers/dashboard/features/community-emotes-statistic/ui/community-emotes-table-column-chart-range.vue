<script setup lang="ts">
import { CheckIcon, GanttChartIcon } from 'lucide-vue-next'


import { useCommunityEmotesStatisticFilters } from '../composables/use-community-emotes-statistic-filters.js'

import type { EmoteStatisticRange } from '~/gql/graphql'



import {
	useTranslatedRanges,
} from '~/features/community-emotes-statistic/composables/use-translated-ranges'

const { t } = useI18n()

const { ranges } = useTranslatedRanges()

const emotesStatisticFilter = useCommunityEmotesStatisticFilters()
</script>

<template>
	<div class="flex items-center space-x-2">
		<UiDropdownMenu>
			<UiDropdownMenuTrigger as-child>
				<UiButton
					variant="ghost"
					size="sm"
					class="-ml-3 h-8 data-[state=open]:bg-accent"
				>
					<span>{{ t('community.emotesStatistic.table.chart') }}</span>
					<GanttChartIcon class="ml-2 h-4 w-4" />
				</UiButton>
			</UiDropdownMenuTrigger>
			<UiDropdownMenuContent align="start">
				<UiDropdownMenuItem
					v-for="[type, text] of Object.entries(ranges)"
					:key="type"
					@click="emotesStatisticFilter.changeTableRange(type as EmoteStatisticRange)"
				>
					<CheckIcon
						v-if="emotesStatisticFilter.tableRange.value === type"
						class="mr-2 h-3.5 w-3.5"
					/>
					{{ text }}
				</UiDropdownMenuItem>
			</UiDropdownMenuContent>
		</UiDropdownMenu>
	</div>
</template>
