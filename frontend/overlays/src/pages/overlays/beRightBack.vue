<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { useIframe } from './brb/iframe.js';
import { useBeRightBackOverlaySocket } from './brb/socket.js';
import BrbTicker, { type Ticker } from './brb/ticker.vue';
import type { SetSettings, OnStart, OnStop } from './brb/types.js';

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

const ticker = ref<Ticker | null>(null);

const setSettings: SetSettings = (s) => {
	settings.value = s;
};

const onStart: OnStart = (minutes, incomingText) => {
	ticker.value?.start(minutes, incomingText);
};

const onStop: OnStop = () => {
	ticker.value?.stop();
};

const iframe = useIframe({
	onSettings: setSettings,
	onStart,
	onStop,
});

const socket = useBeRightBackOverlaySocket({
	apiKey,
	onSettings: setSettings,
	onStart,
	onStop,
});

setTimeout(() => {
	// this should be called on socket event
	onStart(0.1, settings.value.text);
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
	<brb-ticker
		ref="ticker"
		:settings="settings"
	/>
</template>
