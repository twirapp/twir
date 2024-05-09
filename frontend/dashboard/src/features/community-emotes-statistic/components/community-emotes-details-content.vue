<script setup lang="ts">
import { ColorType, createChart } from 'lightweight-charts'
import { RadioGroupItem, RadioGroupRoot } from 'radix-vue'
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import type { AutoscaleInfo, IChartApi , ISeriesApi, UTCTimestamp } from 'lightweight-charts'

import CommunityEmotesDetailsContentUsersHistory
	from '@/features/community-emotes-statistic/components/community-emotes-details-content-users-history.vue'
import CommunityEmotesDetailsContentUsersTop
	from '@/features/community-emotes-statistic/components/community-emotes-details-content-users-top.vue'
import {
	useCommunityEmotesDetails,
	useCommunityEmotesDetailsName,
} from '@/features/community-emotes-statistic/composables/use-community-emotes-details'
import {
	useTranslatedRanges,
} from '@/features/community-emotes-statistic/composables/use-translated-ranges'

const { t } = useI18n()
const { ranges } = useTranslatedRanges()
const { details, range } = useCommunityEmotesDetails()
const { emoteName } = useCommunityEmotesDetailsName()

const chartContainer = ref<HTMLElement>()

let chart = null as IChartApi | null
let areaSeries = null as ISeriesApi<'Area'> | null

onMounted(() => {
	if (!chartContainer.value) return
	chart = createChart(chartContainer.value, {
		height: 240,
		autoSize: true,
		layout: {
			textColor: 'white',
			fontFamily: 'Inter',
			background: { type: ColorType.Solid, color: 'transparent' },
		},
		rightPriceScale: {
			borderColor: '#454545',
		},
		crosshair: {
			horzLine: {
				visible: false,
				color: '#999999',
				labelBackgroundColor: '#555555',
			},
			vertLine: {
				color: '#999999',
				labelBackgroundColor: '#555555',
			},
		},
		timeScale: {
			timeVisible: true,
			borderColor: '#454545',
		},
		grid: {
			horzLines: {
				color: 'rgb(255,255,255,0.15)',
			},
			vertLines: {
				color: 'rgb(255,255,255,0.15)',
			},
		},
	})

	areaSeries = chart.addAreaSeries({
		lineColor: '#01D154',
		topColor: 'rgba(1, 209, 84, 0.5)',
		bottomColor: 'rgb(0,135,54,0.1)',
		lineWidth: 2,
		priceLineColor: 'rgba(2,209,84,0.6)',
		lastValueVisible: false,
		autoscaleInfoProvider: (original: () => AutoscaleInfo | null) => ({
			priceRange: {
				minValue: 0,
				maxValue: original()?.priceRange.maxValue || 1,
			},
		}),
		priceFormat: {
			precision: 0,
			minMove: 1,
		},
	})

	setData()
})

function setData() {
	if (!chart || !areaSeries || !details.value?.emotesStatisticEmoteDetailedInformation?.graphicUsages) return
	areaSeries.setData(details.value.emotesStatisticEmoteDetailedInformation.graphicUsages.map(({ timestamp, count }) => ({
		time: timestamp / 1000 as UTCTimestamp,
		value: count,
	})))
	chart.timeScale().fitContent()
}

watch(details, () => {
	setData()
})

const tableTabs = [
	{ key: 'top', text: t('community.emotesStatistic.details.usersTabs.top') },
	{ key: 'history', text: t('community.emotesStatistic.details.usersTabs.history') },
]

const tableTab = ref<'top' | 'history'>('top')
</script>

<template>
	<div class="flex flex-col divide-y divide-white/10">
		<h1 class="text-4xl font-medium px-6 py-6">
			{{ emoteName }}
		</h1>
		<div class="flex flex-col gap-6 px-6 py-7">
			<div class="flex justify-between flex-wrap">
				<h1 class="text-2xl font-medium">
					{{ t('community.emotesStatistic.details.stats') }}
				</h1>
				<RadioGroupRoot v-model="range" class="inline-flex w-full rounded-[7px] bg-zinc-800 p-px md:w-auto">
					<RadioGroupItem
						v-for="([key, text]) of Object.entries(ranges)"
						:key="key"
						:value="key"
						class="h-8 flex-1 rounded-md px-3 text-[13px] text-white/75 transition-colors hover:bg-white/5 data-[active=true]:bg-white/20 data-[active=true]:text-white data-[active=true]:shadow-md md:flex-auto whitespace-nowrap"
					>
						{{ text }}
					</RadioGroupItem>
				</RadioGroupRoot>
			</div>

			<div
				ref="chartContainer"
				class="relative h-[240px]"
			></div>
		</div>
		<div class="flex flex-col gap-6 px-6 py-7">
			<div class="flex justify-between flex-wrap">
				<h1 class="text-2xl font-medium">
					{{ t('community.emotesStatistic.details.users') }}
				</h1>
				<RadioGroupRoot v-model="tableTab" class="inline-flex w-full rounded-[7px] bg-zinc-800 p-px md:w-auto">
					<RadioGroupItem
						v-for="tab of tableTabs"
						:key="tab.key"
						:value="tab.key"
						class="h-8 flex-1 rounded-md px-3 text-[13px] text-white/75 transition-colors hover:bg-white/5 data-[active=true]:bg-white/20 data-[active=true]:text-white data-[active=true]:shadow-md md:flex-auto whitespace-nowrap"
					>
						{{ tab.text }}
					</RadioGroupItem>
				</RadioGroupRoot>
			</div>
			<CommunityEmotesDetailsContentUsersTop v-if="tableTab === 'top'" />
			<CommunityEmotesDetailsContentUsersHistory v-else />
		</div>
	</div>
</template>
