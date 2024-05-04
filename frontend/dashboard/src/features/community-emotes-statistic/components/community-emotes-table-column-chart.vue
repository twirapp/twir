<script setup lang="ts">
import { ColorType, createChart } from 'lightweight-charts'
import { storeToRefs } from 'pinia'
import { computed, onMounted, onUnmounted, ref, shallowRef, watch } from 'vue'

import { useCommunityChartSize } from '../composables/use-community-chart-size.js'
import { useCommunityChartStyles } from '../composables/use-community-chart-styles.js'

import type { DeepPartial, IChartApi, TimeChartOptions, UTCTimestamp } from 'lightweight-charts'

const props = defineProps<{
	isDayRange: boolean
	usages: {
		timestamp: number
		count: number
	}[]
}>()

const chart = shallowRef<IChartApi | null>(null)
const chartContainer = ref<HTMLElement | null>(null)

const { chartSizes, setChartSize } = useCommunityChartSize()
const { chartStyles } = storeToRefs(useCommunityChartStyles())

const chartOptions = computed<DeepPartial<TimeChartOptions>>(() => ({
	layout: {
		fontSize: 12,
		fontFamily: 'Inter, system-ui, Avenir, Helvetica, Arial, sans-serif',
		textColor: chartStyles.value.textColor,
		background: {
			type: ColorType.Solid,
			color: 'transparent',
		},
	},
	grid: {
		horzLines: {
			visible: false,
		},
		vertLines: {
			visible: false,
		},
	},
	timeScale: {
		fixLeftEdge: true,
		timeVisible: props.isDayRange,
		borderColor: chartStyles.value.borderColor,
	},
	rightPriceScale: {
		borderColor: chartStyles.value.borderColor,
	},
	handleScroll: {
		mouseWheel: false,
	},
	handleScale: {
		axisDoubleClickReset: false,
		axisPressedMouseMove: false,
	},
}))

watch(chartOptions, (options) => {
	if (!chart.value) return
	// styles are not updated :(
	chart.value.applyOptions(options)
})

function resizeHandler() {
	if (!chart.value || !chartContainer.value) return

	const dimensions = chartContainer.value.getBoundingClientRect()
	if (dimensions.width !== 0 || dimensions.height !== 0) {
		setChartSize(dimensions.width, dimensions.height)
	}

	chart.value.resize(chartSizes.value.width, chartSizes.value.height)
	chart.value.timeScale().fitContent()
}

onMounted(() => {
	if (!chartContainer.value) return

	chart.value = createChart(chartContainer.value, chartOptions.value)

	const areaSeries = chart.value.addLineSeries({
		crosshairMarkerVisible: false,
		priceLineVisible: false,
	})

	areaSeries.setData(props.usages.map(({ timestamp, count }) => ({
		time: timestamp / 1000 as UTCTimestamp,
		value: count,
	})))

	resizeHandler()
	window.addEventListener('resize', resizeHandler)
})

onUnmounted(() => {
	if (!chart.value) return

	chart.value.remove()
	chart.value = null

	window.removeEventListener('resize', resizeHandler)
})
</script>

<template>
	<div ref="chartContainer" class="w-full h-[100px]"></div>
</template>
