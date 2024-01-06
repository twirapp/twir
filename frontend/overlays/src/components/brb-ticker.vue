<script setup lang="ts">
import { useFontSource } from '@twir/fontsource';
import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { useIntervalFn } from '@vueuse/core';
import { computed, ref, watch } from 'vue';

import { getTimeDiffInMilliseconds, millisecondsToTime } from '@/helpers.js';
import type { OnStart, OnStop } from '@/types.js';


const props = defineProps<{
	settings: Settings
}>();

const minutes = ref(0);
const text = ref<string>();
const countDownTicks = ref(0);
const countUpTicks = ref(0);

export type Ticker = {
	start: OnStart,
	stop: OnStop,
}

const countDownInterval = useIntervalFn(() => {
	countDownTicks.value--;

	if (!countDownTicks.value) {
		countDownInterval.pause();
	}

	if (!countDownTicks.value && props.settings.late?.enabled) {
		countUpInterval.resume();
	}
}, 1000, {
	immediate: false,
});

const countUpInterval = useIntervalFn(() => {
	countUpTicks.value++;
}, 1000, {
	immediate: false,
});

const start: OnStart = (incomingMinutes, incomingText) => {
	stop();
	const ticks = parseInt((getTimeDiffInMilliseconds(incomingMinutes) / 1000).toString());

	countDownTicks.value = ticks;
	minutes.value = getTimeDiffInMilliseconds(incomingMinutes);
	text.value = incomingText;

	countDownInterval.resume();
};

const stop: OnStop = () => {
	countDownTicks.value = 0;
	countUpTicks.value = 0;
	minutes.value = 0;
	text.value = undefined;

	countDownInterval.pause();
	countUpInterval.pause();
};

defineExpose({
	start,
	stop,
});

const showCountDown = computed(() => {
	const isActive = countDownInterval.isActive.value;

	if (isActive) return true;
	if (countUpInterval.isActive.value && !props.settings.late?.displayBrbTime) return false;

	return true;
});

const fontSource = useFontSource();
watch(() => props.settings.fontFamily, () => {
	fontSource.loadFont(props.settings.fontFamily, 400, 'normal');
}, { immediate: true });

const fontFamily = computed(() => {
	return `"${props.settings.fontFamily}-400-normal"`;
});
</script>

<template>
	<Transition name="overlay" appear>
		<div
			v-if="countDownInterval.isActive.value || countUpInterval.isActive.value"
			class="overlay"
			:style="{
				backgroundColor: settings.backgroundColor || 'rgba(9, 8, 8, 0.50)',
				color: settings.fontColor || '#fff',
			}"
		>
			<div
				v-if="showCountDown"
				class="count-up"
				:style="{ fontSize: `${settings.fontSize / (countUpInterval.isActive.value ? 2 : 1)}px` }"
			>
				{{ text || settings.text }}
				{{
					countDownTicks > 0
						? millisecondsToTime(countDownTicks * 1000)
						: millisecondsToTime(minutes)
				}}
			</div>
			<div
				v-if="countUpInterval.isActive.value && props.settings.late?.enabled"
				class="count-down"
				:style="{ fontSize: `${settings.fontSize}px`}"
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
