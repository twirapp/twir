<script setup lang="ts">
import { computed } from 'vue'

import type { Layer } from '@/composables/overlays/use-overlays.js'

const props = defineProps<{
	layer: Layer
	zIndex?: number
}>()

const emoteUrl = computed(() => props.layer.settings.emoteUrl.trim())
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
		<img
			v-if="emoteUrl"
			:src="emoteUrl"
			:alt="layer.settings.emoteName"
			style="width: 100%; height: 100%; object-fit: contain; display: block;"
		/>
	</div>
</template>
