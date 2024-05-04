<script setup lang="ts">
import { ColorType, createChart } from 'lightweight-charts'
import { onMounted, onUnmounted, ref } from 'vue'

import type { IChartApi , UTCTimestamp } from 'lightweight-charts'

const props = defineProps<{
	usages: {
		timestamp: number
		count: number
	}[]
}>()

let chart: IChartApi | null = null
const chartContainer = ref<HTMLElement>()

onMounted(() => {
	if (!chartContainer.value) {
		return
	}
	chart = createChart(chartContainer.value, {
		layout: {
			textColor: 'black',
			background: { type: ColorType.Solid, color: 'transparent' },
		},
		grid: {
			horzLines: {
				visible: false,
			},
			vertLines: {
				visible: false,
			},
		},
		width: 300,
		height: 50,
		leftPriceScale: {
			visible: false,
		},
		rightPriceScale: { visible: false },
		timeScale: {
			fixLeftEdge: true,
			visible: false, // эта хрень отключает точки снизу, но при этом и тултипы отключает. Хз чё делать.
		},
		crosshair: {
			horzLine: {
				visible: false,
			},
		},
		handleScroll: {
			horzTouchDrag: false,
			mouseWheel: false,
			pressedMouseMove: false,
			vertTouchDrag: false,
		},
		handleScale: {
			mouseWheel: false,
			axisDoubleClickReset: false,
			axisPressedMouseMove: false,
			pinch: false,
		},
	})
	const areaSeries = chart.addLineSeries({
		crosshairMarkerVisible: false,
		priceLineVisible: false,
	})

	areaSeries.setData(props.usages.map(({ timestamp, count }) => ({
		time: timestamp / 1000 as UTCTimestamp,
		value: count,
	})))

	chart.timeScale().fitContent()
})

onUnmounted(() => {
	if (chart) {
		chart.remove()
		chart = null
	}
})
</script>

<template>
	<div ref="chartContainer" class="lw-chart h-full"></div>
</template>
