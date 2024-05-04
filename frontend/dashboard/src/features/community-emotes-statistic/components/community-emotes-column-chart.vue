<script setup lang="ts">
import { ColorType, createChart } from 'lightweight-charts'
import { onMounted, onUnmounted, ref } from 'vue'

import type { AreaData, IChartApi, Time, UTCTimestamp } from 'lightweight-charts'

const props = defineProps<{
	usages: {
		timestamp: number
		count: number
	}[]
}>()

let chart: IChartApi | null = null
const chartContainer = ref<HTMLElement>()

const tooltipVisible = ref(false)
const tooltipBarData = ref<AreaData<Time>>()
const tooltipPosition = ref({ x: 0, y: 0 })

const TOOLTIP_WIDTH = 108
const TOOLTIP_HEIGHT = 48
const TOOLTIP_MARGIN = 12

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
		height: 50,
		leftPriceScale: {
			visible: false,
		},
		rightPriceScale: { visible: false },
		timeScale: {
			fixLeftEdge: true,
			visible: false,
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

	chart.subscribeCrosshairMove((param) => {
		if (
			!chartContainer.value
			|| !chart
			|| !areaSeries
			|| !param.logical
			|| !param.point
			|| !param.time
		) {
			tooltipVisible.value = false
			return
		}

		// check if the cursor is outside the chart pane
		if (
			param.point.x < 0
			|| param.point.x > chartContainer.value.clientWidth
			|| param.point.y < 0
			|| param.point.y > chartContainer.value.clientHeight
		) {
			tooltipVisible.value = false
			return
		}

		const barData = areaSeries.dataByIndex(param.logical)
		if (barData === null || !('value' in barData)) {
			tooltipVisible.value = false
			return
		}

		const paneSize = chart.paneSize()

		let x = param.point.x + TOOLTIP_MARGIN
		if (x > paneSize.width - TOOLTIP_WIDTH - TOOLTIP_MARGIN) {
			x = param.point.x - TOOLTIP_WIDTH - TOOLTIP_MARGIN
		}

		let y = param.point.y + TOOLTIP_MARGIN
		if (y > paneSize.height - TOOLTIP_HEIGHT - TOOLTIP_MARGIN) {
			y = param.point.y - TOOLTIP_HEIGHT - TOOLTIP_MARGIN
		}

		tooltipPosition.value = { x, y }
		tooltipBarData.value = barData
		tooltipVisible.value = true
	})
})

onUnmounted(() => {
	if (chart) {
		chart.remove()
		chart = null
	}
})
</script>

<template>
	<div ref="chartContainer" class="lw-chart h-full relative">
		<div
			v-if="tooltipVisible && tooltipBarData"
			class="pointer-events-none absolute z-50 flex min-w-[108px] flex-col whitespace-nowrap rounded-md bg-stone-400/90 px-2.5 py-1.5 shadow-md backdrop-blur-sm"
			:style="{
				left: `${tooltipPosition.x}px`,
				top: `${tooltipPosition.y}px`,
			}"
		>
			<div class="inline-flex items-center">
				<span class="text-sm font-medium">{{ tooltipBarData.value }}</span>
			</div>

			<span class="text-xs text-white/75">
				{{ new Date(tooltipBarData.time * 1000 as number).toLocaleDateString() }}, {{ new Date(tooltipBarData.time * 1000 as number).toLocaleTimeString() }}
			</span>
		</div>
	</div>
</template>
