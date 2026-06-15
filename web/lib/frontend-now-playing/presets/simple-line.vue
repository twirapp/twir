<script setup lang="ts">
import { useElementSize } from '@vueuse/core'
import { computed, ref } from 'vue'

import type { Settings, Track } from '../types.js'

defineProps<{
	track?: Track | null
	settings: Settings
}>()

const spotifyRef = ref<HTMLElement>()
const infoRef = ref<HTMLElement>()

const { width: spotifyWidth } = useElementSize(spotifyRef)
const { width: infoWidth } = useElementSize(infoRef)

const nameMarqueEnabled = computed(() => {
	return infoWidth.value > spotifyWidth.value
})
</script>

<template>
	<div v-if="track" ref="spotifyRef" class="spotify">
		<img
			v-if="settings.showImage" class="cover"
			:src="track.imageUrl ?? '/overlays/images/play.png'"
		/>
		<div ref="infoRef" class="info" :class="{ marque: nameMarqueEnabled }">
			<div class="name">
				{{ track.title }}
			</div>
			<div class="separator">
				–
			</div>
			<div class="artist">
				{{ track.artist }}
			</div>
		</div>
	</div>
</template>

<style scoped>
.spotify {
	white-space: nowrap;
	overflow: hidden;
	display: flex;
	align-items: center;
	column-gap: 16px;
	background-color: v-bind('settings.backgroundColor');

	/* Font */
	font-family: Inter;
	font-size: 24px;
	color: #fff;
	padding-top: 4px;
	padding-bottom: 4px;
	padding-left: 8px;
	padding-right: 8px;
	border-radius: 8px;
}

.cover {
	width: 30px;
	height: 30px;
	border-radius: 5px;
}

.info {
	overflow: hidden;
}

.name, .separator, .artist {
	display: inline-block;
}
</style>
