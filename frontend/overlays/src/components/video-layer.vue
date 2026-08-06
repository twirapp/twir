<script setup lang="ts">
import { computed } from 'vue'

import type { Layer } from '@/composables/overlays/use-overlays.js'

const props = defineProps<{
	layer: Layer
	zIndex?: number
}>()

const videoUrl = computed(() => props.layer.settings.videoUrl.trim())
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
		<video
			v-if="videoUrl"
			:src="videoUrl"
			autoplay
			:loop="layer.settings.videoLoop"
			:muted="layer.settings.videoMuted"
			playsinline
			style="width: 100%; height: 100%; object-fit: cover; display: block;"
		/>
	</div>
</template>
