<script setup lang="ts">
import { storeToRefs } from 'pinia';

import { useNowPlayingData } from '@/composables/now-playing/use-now-playing-data.ts';

const { currentTrack } = storeToRefs(useNowPlayingData());
</script>

<template>
	<div v-if="currentTrack" class="spotify">
		<img class="cover" :src="currentTrack.image_url ?? '/overlays/public/images/play.png'"/>
		<div class="info">
			<span
				class="artist"
			>
				{{ currentTrack.artist }}
			</span>
			<span
				class="name"
			>
				{{ currentTrack.title }}
			</span>
		</div>
	</div>
</template>

<style scoped>
.spotify {
	height: 70px;
	width: 350px;
	background: #1E1E1E;
	display: flex;
	align-items: center;

	font-family: Inter;
	font-size: 24px;
	color: #fff;
	border-radius: 10px;
	margin: 5px;
}

.artist {
	color: #AEAEAE;
	font-weight: 400;
	margin: 0;
	font-size: 1rem;
	text-transform: uppercase;
}

.cover {
	width: 70px;
	height: 70px;
	border-radius: 10px 0 0 10px;
}

.info {
	display: flex;
	flex-direction: column;
	padding-left: 5px;
	overflow: hidden;
	white-space: nowrap;
}

.info span {
	overflow: hidden;
	text-overflow: ellipsis;
}

.name {
	color: #06AB4F;
	margin: 4px 0 0;
	font-size: 1.5rem;
}
</style>
