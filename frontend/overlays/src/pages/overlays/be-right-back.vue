<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import BrbTicker, { type Ticker } from '@/components/brb-ticker.vue';
import { useBrbIframe } from '@/composables/brb/use-brb-iframe.js';
import { useBeRightBackOverlaySocket } from '@/composables/brb/use-brb-socket.js';
import { generateSocketUrlWithParams } from '@/helpers';
import type { SetSettings, OnStart, OnStop } from '@/types.js';

const route = useRoute();
const apiKey = route.params.apiKey as string;

const settings = ref<Settings>();

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

const iframe = useBrbIframe({
	onSettings: setSettings,
	onStart,
	onStop,
});

const brbUrl = generateSocketUrlWithParams('/overlays/brb', {
	apiKey,
});

const socket = useBeRightBackOverlaySocket({
	brbUrl,
	onSettings: setSettings,
	onStart,
	onStop,
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
	<div class="container">
		<brb-ticker
			v-if="settings"
			ref="ticker"
			:settings="settings"
		/>
	</div>
</template>

<style scoped>
.container {
	overflow: hidden;
}
</style>
