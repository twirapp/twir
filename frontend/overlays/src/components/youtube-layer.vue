<script setup lang="ts">
import { computed } from 'vue'

import type { Layer } from '@/composables/overlays/use-overlays.js'

const props = defineProps<{
	layer: Layer
	zIndex?: number
}>()

const youtubeUrl = computed(() => {
	const videoId = props.layer.settings.youtubeVideoId.trim()
	if (!videoId) return ''

	const autoplay = props.layer.settings.youtubeAutoplay ? 1 : 0
	const mute = props.layer.settings.youtubeMuted ? 1 : 0
	const loop = props.layer.settings.youtubeLoop ? 1 : 0
	const playlist = props.layer.settings.youtubeLoop ? `&playlist=${encodeURIComponent(videoId)}` : ''

	return `https://www.youtube-nocookie.com/embed/${encodeURIComponent(videoId)}?autoplay=${autoplay}&mute=${mute}&loop=${loop}${playlist}&controls=0&playsinline=1`
})
</script>

<template>
	<div
		:id="'layer' + layer.id"
		style="position: absolute; overflow: hidden;"
		:style="{
			top: `${layer.posY}px`,
			left: `${layer.posX}px`,
			width: `${layer.width}px`,
			height: `${layer.height}px`,
			transform: `rotate(${layer.rotation || 0}deg)`,
			transformOrigin: 'center center',
			zIndex: zIndex ?? 0,
		}"
	>
		<iframe
			v-if="youtubeUrl"
			title="YouTube video"
			:src="youtubeUrl"
			allow="autoplay"
			style="width: 100%; height: 100%; border: none; display: block;"
		/>
	</div>
</template>
