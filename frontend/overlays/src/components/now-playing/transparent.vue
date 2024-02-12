<script setup lang="ts">
import { storeToRefs } from 'pinia';

import { useNowPlayingData } from '@/composables/now-playing/use-now-playing-data.ts';

const { currentTrack } = storeToRefs(useNowPlayingData());
</script>

<template>
	<div v-if="currentTrack" class="spotify">
		<img class="image" :src="currentTrack.image_url ?? '/overlays/public/images/play.png'" />
		<div class="info">
			<span class="name">{{ currentTrack.title }}</span>
			<div class="artist">
				{{ currentTrack.artist }}
			</div>
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
