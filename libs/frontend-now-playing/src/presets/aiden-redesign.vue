<script setup lang="ts">
import { Settings } from '@twir/api/messages/overlays_now_playing/overlays_now_playing';
import { computed } from 'vue';

import type { Track } from '../types.js';

const props = defineProps<{
	track?: Track | null
	settings: Settings
}>();

const bgColor = computed(() => {
	if (props.settings.backgroundColor === 'rgba(0, 0, 0, 0)') {
		return '#1E1E1E';
	}

	return props.settings.backgroundColor;
});
</script>

<template>
	<div v-if="track" class="spotify">
		<img class="cover" :src="track.image_url ?? '/overlays/public/images/play.png'" />
		<div class="info">
			<span
				class="artist"
			>
				{{ track.artist }}
			</span>
			<span
				class="name"
			>
				{{ track.title }}
			</span>
		</div>
	</div>
</template>

<style scoped>
.spotify {
	height: 70px;
	width: 350px;
	background-color: v-bind(bgColor);
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
	padding-left: 15px;
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
