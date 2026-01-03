<script setup lang="ts">
import { computed } from 'vue'

import type { Layer } from '@/composables/overlays/use-overlays.js'

const props = defineProps<{
	layer: Layer
}>()

const imageUrl = computed(() => props.layer.settings.imageUrl || '')

// Check if imageUrl is valid (not empty/null/undefined)
const hasValidUrl = computed(() => {
	return imageUrl.value && imageUrl.value.trim().length > 0
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
		}"
	>
		<img
			v-if="hasValidUrl"
			:src="imageUrl"
			alt="Overlay image"
			style="width: 100%; height: 100%; object-fit: contain; display: block;"
		/>
	</div>
</template>
