<script setup lang="ts">
import { ChannelOverlayNowPlayingPreset } from '@twir/types/api';
import { computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';

import PresetAidenRedesign from '@/components/now-playing/aiden-redesign.vue';
import PresetTransparent from '@/components/now-playing/transparent.vue';
import { useNowPlayingSocket } from '@/composables/now-playing/use-now-playing-socket.ts';

const route = useRoute();

const socket = useNowPlayingSocket({
	apiKey: route.params.apiKey as string,
	overlayId: route.query.id as string,
});

onMounted(() => {
	if (!route.params.apiKey || !route.query.id) {
		return;
	}

	socket.connect();
});

const presetComponent = computed(() => {
	switch (socket.settings.value?.preset) {
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
	<component :is="presetComponent" :track="socket.track.value"/>
</template>
