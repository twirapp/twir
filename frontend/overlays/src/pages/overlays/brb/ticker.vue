<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { useIntervalFn } from '@vueuse/core';
import { ref } from 'vue';

import { getTimeDiffInMilliseconds, millisecondsToTime } from './timeUtils.js';
import { OnStart, OnStop } from './types';

defineProps<{
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
</script>

<template>
	<div v-if="countDownInterval.isActive.value || countUpInterval.isActive.value" class="overlay">
		<div
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
			v-if="countUpInterval.isActive.value"
			class="count-down"
			:style="{ fontSize: `${settings.fontSize}px`}"
		>
			{{ settings.late?.text }}	 {{ millisecondsToTime(countUpTicks * 1000) }}
		</div>
	</div>
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
	color: v-bind('settings.fontColor');
	flex-direction: column;
	background-color: v-bind('settings.backgroundColor');
}
</style>
