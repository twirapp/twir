<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { useIframe } from './brb/iframe.js';
import { useBeRightBackOverlaySocket } from './brb/socket.js';
import { getTimeDiffInMilliseconds, millisecondsToTime } from './brb/timeUtils.js';

const route = useRoute();
const apiKey = route.params.apiKey as string;

const settings = ref<Settings>({
	fontSize: 100,
	fontColor: '#fff',
	backgroundColor: 'rgb(231, 220, 220, 0.5)',
	text: 'AFK FOR',
	late: {
		text: 'LATE FOR',
		displayBrbTime: true,
		displayLateTime: true,
		enabled: true,
	},
});
const setSettings = (s: Settings) => {
	settings.value = s;
};

const iframe = useIframe(setSettings);
const socket = useBeRightBackOverlaySocket(apiKey, setSettings);

const countDownTicks = ref(0);
const isRunned = ref(false);
const isCountUpRunned = ref(false);
const countUpTicks = ref(0);

setTimeout(() => isRunned.value = true, 1000);

watch(isRunned, (v) => {
	if (v) {
		countDownTicks.value = parseInt((getTimeDiffInMilliseconds(0.1) / 1000).toString());
	} else {
		countDownTicks.value = 0;
	}
});

watch(isCountUpRunned, (v) => {
	if (!v) countUpTicks.value = 0;

	countUpTicks.value++;
});

watch(countUpTicks, (v) => {
	if (!v) return;

	setTimeout(() => countUpTicks.value++, 1000);
});

watch(countDownTicks, (v) => {
	if (!v) {
		if (settings.value.late?.enabled) {
			isCountUpRunned.value = true;
		}

		return;
	}

	setTimeout(() => countDownTicks.value--, 1000);
});

onMounted(() => {
	if (window.frameElement) {
		iframe.create();
	} else {
		socket.create();
	}

	return () => {
		iframe.destroy();
		socket.destroy();
	};
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
