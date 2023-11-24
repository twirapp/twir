<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { useIntervalFn } from '@vueuse/core';
import { computed, ref } from 'vue';

import { getTimeDiffInMilliseconds, millisecondsToTime } from './timeUtils.js';
import { OnStart, OnStop } from './types';

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

const fontUrl = computed(() => {
	return `https://fonts.googleapis.com/css?family=${props.settings.fontFamily}`;
});

const fontFamily = computed(() => {
	try {
		const [family] = props.settings.fontFamily.split(':');

		return family;
	} catch (e) {
		return '';
	}
});

const showCountDown = computed(() => {
	const isActive = countDownInterval.isActive.value;

	if (isActive) return true;
	if (countUpInterval.isActive.value && !props.settings.late?.displayBrbTime) return false;

	return true;
});

</script>

<template>
	<component :is="'style'">
		@import url('{{ fontUrl }}')
	</component>
	<div
		v-if="countDownInterval.isActive.value || countUpInterval.isActive.value"
		class="overlay"
		:style="{
			backgroundColor: settings.backgroundColor || 'rgba(9, 8, 8, 0.49)',
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
	flex-direction: column;
	font-family: v-bind(fontFamily);
}
</style>
