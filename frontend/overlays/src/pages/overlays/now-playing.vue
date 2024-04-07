<script setup lang="ts">
import { NowPlaying, type Track } from '@twir/frontend-now-playing';
import { useDebounceFn } from '@vueuse/core';
import { onMounted, onUnmounted, ref, watch } from 'vue';
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

const showedTrack = ref<Track | null | undefined>(null);

const debouncedShowTrack = useDebounceFn((track: Track | null | undefined) => {
	showedTrack.value = track;
	if (settings.value?.hideTimeout) {
		setTimeout(() => {
			showedTrack.value = null;
		}, settings.value.hideTimeout * 1000);
	}
}, 5000);

watch(currentTrack, (track) => {
	debouncedShowTrack(track);
});
</script>

<template>
	<NowPlaying v-if="settings" :settings="settings" :track="showedTrack" />
</template>
