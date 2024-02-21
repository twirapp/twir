<script setup lang="ts">
import { Settings } from '@twir/api/messages/overlays_now_playing/overlays_now_playing';
import { useElementSize } from '@vueuse/core';
import { computed, ref } from 'vue';

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

const infoRef = ref<HTMLElement>();
const nameRef = ref<HTMLElement>();
const artistRef = ref<HTMLElement>();

const { width: infoWidth } = useElementSize(infoRef);
const { width: nameWidth } = useElementSize(nameRef);
const { width: artistWidth } = useElementSize(artistRef);

const nameMarqueEnabled = computed(() => {
	return nameWidth.value > infoWidth.value;
});

const artistMarqueEnabled = computed(() => {
	return artistWidth.value > infoWidth.value;
});
</script>

<template>
	<div v-if="track" class="spotify">
		<img
			v-if="settings.showImage" class="cover"
			:src="track.image_url ?? '/overlays/images/play.png'"
		/>
		<div ref="infoRef" class="info">
			<div ref="artistRef" class="artist" :class="{ marque: artistMarqueEnabled }">
				{{ track.artist }}
			</div>
			<div ref="nameRef" class="name" :class="{ marque: nameMarqueEnabled }">
				{{ track.title }}
			</div>
		</div>
	</div>
</template>

<style scoped>
.spotify {
	white-space: nowrap;
	overflow: hidden;

	height: 70px;
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
	padding-left: 15px;
	overflow: hidden;
}

.info span {
	overflow: hidden;
	text-overflow: ellipsis;
}

.name, .artist {
	white-space: nowrap;
	width: max-content;
}

.name {
	color: #06AB4F;
	margin: 4px 0 0;
	font-size: 1.5rem;
}
</style>
