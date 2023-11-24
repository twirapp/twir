<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { useIframe } from './brb/iframe.js';
import { useBeRightBackOverlaySocket } from './brb/socket.js';
import BrbTicker from './brb/ticker.vue';
import { getTimeDiffInMilliseconds } from './brb/timeUtils.js';

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

setTimeout(() => {
	// this should be called on socket event
	countDownTicks.value = parseInt((getTimeDiffInMilliseconds(0.1) / 1000).toString());
}, 1000);

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
	<BrbTicker :settings="settings" :start-ticks="countDownTicks" />
</template>
