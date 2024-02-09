<script setup lang="ts">
import type { Track } from '@/composables/now-playing/use-now-playing-socket.ts';

defineProps<{
	track: Track | undefined | null,
}>();
</script>

<template>
	<div v-if="track" class="spotify">
		<img class="image" :src="track.image_url ?? '/overlays/public/images/play.png'"/>
		<div class="info">
			<span class="name">{{ track.title }}</span>
			<div class="artist">{{ track.artist }}</div>
		</div>
	</div>
</template>

<style scoped>
body {
	display: flex;
}

.spotify {
	display: flex;
	align-items: center;
	column-gap: 16px;

	/* Font */
	font-family: Inter;
	font-size: 24px;
	color: #fff;
}

.image {
	width: 52px;
	height: 52px;
	border-radius: 8px;
	background-color: #e9e9e9;
	overflow: hidden;
}

.artist {
	opacity: 0.6;
	font-size: 16px;
}

.info {
	display: flex;
	flex-direction: column;
	gap: 4px;
}
</style>
