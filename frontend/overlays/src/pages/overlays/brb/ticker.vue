<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { computed, ref, watch, watchEffect } from 'vue';

import { getTimeDiffInMilliseconds, millisecondsToTime } from './timeUtils.js';

const props = defineProps<{
	settings: Settings
	startTicks: number
}>();

const countDownTicks = ref(0);
watchEffect(() => countDownTicks.value = props.startTicks);

const countUpTicks = ref(0);

const isRunned = computed(() => countDownTicks.value > 0 || countUpTicks.value > 0);

watch(countUpTicks, (v) => {
	if (!v) return;

	setTimeout(() => countUpTicks.value++, 1000);
});

watch(countDownTicks, (v) => {
	if (!v) {
		if (props.settings.late?.enabled) {
			countUpTicks.value++;
		}

		return;
	} else {
		countUpTicks.value = 0;
	}

	setTimeout(() => countDownTicks.value--, 1000);
});
</script>

<template>
	<div v-if="isRunned" class="overlay">
		<div
			class="count-up"
			:style="{ fontSize: `${settings.fontSize / (countUpTicks > 0 ? 2 : 1)}px` }"
		>
			{{ settings.text }}
			{{
				countDownTicks > 0
					? millisecondsToTime(countDownTicks * 1000)
					: millisecondsToTime(getTimeDiffInMilliseconds(5))
			}}
		</div>
		<div
			v-if="countUpTicks > 0"
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
