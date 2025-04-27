<script setup lang="ts">
import { useFontSource } from '@twir/fontsource'
import { useIntervalFn } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

import type { BrbOnStartFn, BrbOnStopFn } from '@/types.js'

import { useBrbSettings } from '@/composables/brb/use-brb-settings.js'
import { getTimeDiffInMilliseconds, millisecondsToTime } from '@/helpers.js'

const { settings } = useBrbSettings()

const minutes = ref(0)
const text = ref<string | null>(null)
const countDownTicks = ref(0)
const countUpTicks = ref(0)

const countUpInterval = useIntervalFn(() => {
	countUpTicks.value++
}, 1000, { immediate: false })

const countDownInterval = useIntervalFn(() => {
	countDownTicks.value--

	if (!countDownTicks.value) {
		countDownInterval.pause()
	}

	if (!countDownTicks.value && settings.value?.late?.enabled) {
		countUpInterval.resume()
	}
}, 1000, { immediate: false })

const stop: BrbOnStopFn = () => {
	countDownTicks.value = 0
	countUpTicks.value = 0
	minutes.value = 0
	text.value = null

	countDownInterval.pause()
	countUpInterval.pause()
}

const start: BrbOnStartFn = (incomingMinutes, incomingText) => {
	stop()
	const ticks = Number.parseInt((getTimeDiffInMilliseconds(incomingMinutes) / 1000).toString())

	countDownTicks.value = ticks
	minutes.value = getTimeDiffInMilliseconds(incomingMinutes)
	text.value = incomingText

	countDownInterval.resume()
}

export interface BrbTimerMethods {
	start: BrbOnStartFn
	stop: BrbOnStopFn
}

defineExpose<BrbTimerMethods>({
	start,
	stop,
})

const showCountDown = computed(() => {
	const isActive = countDownInterval.isActive.value

	if (isActive) return true
	if (countUpInterval.isActive.value && !settings.value?.late?.displayBrbTime) return false

	return true
})

const fontSource = useFontSource(false)
watch(() => settings.value?.fontFamily, (font) => {
	if (!font) return
	fontSource.loadFont(font, 400, 'normal')
}, { immediate: true })

const fontFamily = computed(() => {
	return `"${settings.value?.fontFamily}-400-normal"`
})

const backgroundColor = computed(() => {
	return settings.value?.backgroundColor || 'rgba(9, 8, 8, 0.50)'
})

const fontColor = computed(() => {
	return settings.value?.fontColor || '#fff'
})

const countUpFontSize = computed(() => {
	if (!settings.value?.fontSize) return '16px'
	return `${settings.value.fontSize / (countUpInterval.isActive.value ? 2 : 1)}px`
})

const countDownFontSize = computed(() => {
	return `${settings.value?.fontSize || 16}px`
})
</script>

<template>
	<Transition v-if="settings" name="overlay" appear>
		<div
			v-if="countDownInterval.isActive.value || countUpInterval.isActive.value"
			id="brb-overlay"
			class="overlay"
		>
			<div
				v-if="showCountDown"
				id="brb-count-up"
				class="count-up"
			>
				{{ text || settings.text }}
				{{
					countDownTicks > 0
						? millisecondsToTime(countDownTicks * 1000)
						: millisecondsToTime(minutes)
				}}
			</div>
			<div
				v-if="countUpInterval.isActive.value && settings.late?.enabled"
				id="brb-count-down"
				class="count-down"
			>
				{{ settings.late?.text }} {{ millisecondsToTime(countUpTicks * 1000) }}
			</div>
		</div>
	</Transition>
</template>

<style scoped>
.overlay {
	width: 100vw;
	height: 100vh;
	margin: 0;
	display: flex;
	justify-content: center;
	align-items: center;
	text-align: center;
	flex-direction: column;
	font-family: v-bind(fontFamily);
	background-color: v-bind(backgroundColor);
	color: v-bind(fontColor);
	font-variant-numeric: tabular-nums;
}

.count-up {
	font-size: v-bind(countUpFontSize);
}

.count-down {
	font-size: v-bind(countDownFontSize);
}

.overlay-enter-active,
.overlay-leave-active {
	transition: opacity 0.9s ease;
}

.overlay-enter-from,
.overlay-leave-to {
	opacity: 0;
}
</style>
