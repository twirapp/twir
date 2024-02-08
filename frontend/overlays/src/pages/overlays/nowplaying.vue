<script setup lang="ts">
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';

import { useNowPlayingSocket } from '@/composables/nowplaying/use-nowplaying-socket.ts';

const route = useRoute();

const socket = useNowPlayingSocket({
	apiKey: route.params.apiKey as string,
});

onMounted(() => {
	socket.connect();
});
</script>

<template>
	<div>
		<h1>Now Playing</h1>
		<h2>{{ socket.track }}</h2>
	</div>
</template>

<style scoped>

</style>
