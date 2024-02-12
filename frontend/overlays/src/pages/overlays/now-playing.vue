<script setup lang="ts">
import { ChannelOverlayNowPlayingPreset } from '@twir/types/api';
import { storeToRefs } from 'pinia';
import { computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';

import PresetAidenRedesign from '@/components/now-playing/aiden-redesign.vue';
import PresetTransparent from '@/components/now-playing/transparent.vue';
import { useNowPlayingData } from '@/composables/now-playing/use-now-playing-data.ts';
import { useNowPlayingIframe } from '@/composables/now-playing/use-now-playing-iframe.ts';
import { useNowPlayingSocket } from '@/composables/now-playing/use-now-playing-socket.ts';

const route = useRoute();

const socket = useNowPlayingSocket({
	apiKey: route.params.apiKey as string,
	overlayId: route.query.id as string,
});
const iframe = useNowPlayingIframe();
const { settings } = storeToRefs(useNowPlayingData());

onMounted(() => {
	iframe.connect();

	if (route.params.apiKey && route.query.id && !iframe.isIframe) {
		socket.connect();
	}
});

onUnmounted(() => {
	iframe.destroy();
	socket.destroy();
});

const presetComponent = computed(() => {
	switch (settings.value?.preset) {
		case ChannelOverlayNowPlayingPreset.TRANSPARENT:
			return PresetTransparent;
		case ChannelOverlayNowPlayingPreset.AIDEN_REDESIGN:
			return PresetAidenRedesign;
		default:
			return PresetTransparent;
	}
});
</script>

<template>
	<component :is="presetComponent" />
</template>

<style>
body {
	background-color: #000;
}
</style>
