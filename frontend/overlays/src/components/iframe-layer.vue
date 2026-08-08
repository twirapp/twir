<script setup lang="ts">
import { computed } from 'vue'

import type { Layer } from '@/composables/overlays/use-overlays.js'

const props = defineProps<{
	layer: Layer
	zIndex?: number
}>()

const iframeUrl = computed(() => props.layer.settings.iframeUrl.trim())
const iframeScale = computed(() => props.layer.settings.iframeScale || 1)
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
			v-if="iframeUrl"
			:title="layer.settings.widgetKey || 'Overlay iframe'"
			:src="iframeUrl"
			:style="{
				width: `${100 / iframeScale}%`,
				height: `${100 / iframeScale}%`,
				border: 'none',
				display: 'block',
				transform: `scale(${iframeScale})`,
				transformOrigin: 'top left',
			}"
		/>
	</div>
</template>
