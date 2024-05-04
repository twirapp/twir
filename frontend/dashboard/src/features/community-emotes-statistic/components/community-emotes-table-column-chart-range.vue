<script setup lang="ts">
import { CheckIcon, GanttChartIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommunityEmotesStatisticFilters } from '../composables/use-community-emotes-statistic-filters.js'

import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { EmoteStatisticRange } from '@/gql/graphql'

const { t } = useI18n()

const ranges = computed(() => ({
	[EmoteStatisticRange.LastDay]: t('community.emotesStatistic.table.lastDay'),
	[EmoteStatisticRange.LastWeek]: t('community.emotesStatistic.table.lastWeek'),
	[EmoteStatisticRange.LastMonth]: t('community.emotesStatistic.table.lastMonth'),
	[EmoteStatisticRange.LastThreeMonth]: t('community.emotesStatistic.table.lastThreeMonth'),
	[EmoteStatisticRange.LastYear]: t('community.emotesStatistic.table.lastYear'),
}))

const emotesStatisticFilter = useCommunityEmotesStatisticFilters()
</script>

<template>
	<div class="flex items-center space-x-2">
		<DropdownMenu>
			<DropdownMenuTrigger as-child>
				<Button
					variant="ghost"
					size="sm"
					class="-ml-3 h-8 data-[state=open]:bg-accent"
				>
					<span>{{ t('community.emotesStatistic.table.chart') }}</span>
					<GanttChartIcon class="ml-2 h-4 w-4" />
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="start">
				<DropdownMenuItem
					v-for="[type, text] of Object.entries(ranges)"
					:key="type"
					@click="emotesStatisticFilter.changeTableRange(type as EmoteStatisticRange)"
				>
					<CheckIcon
						v-if="emotesStatisticFilter.tableRange === type"
						class="mr-2 h-3.5 w-3.5"
					/>
					{{ text }}
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	</div>
</template>
