<script setup lang="ts">
import { NowPlaying } from '@twir/frontend-now-playing';
import { onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';

import { useNowPlayingSocket } from '@/composables/now-playing/use-now-playing-socket.ts';

const route = useRoute();

const { connect, destroy, currentTrack, settings } = useNowPlayingSocket({
	apiKey: route.params.apiKey as string,
	overlayId: route.query.id as string,
});

onMounted(() => {
	if (route.params.apiKey && route.query.id) {
		connect();
	}
});

onUnmounted(() => {
	destroy();
});
</script>

<template>
	<NowPlaying v-if="settings" :settings="settings" :track="currentTrack" />
</template>
