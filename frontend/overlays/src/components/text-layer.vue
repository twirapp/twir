<script setup lang="ts">
import type { Layer } from '@/composables/overlays/use-overlays.js'

defineProps<{
	layer: Layer
	zIndex?: number
}>()
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
		<div
			:style="{
				width: '100%',
				height: '100%',
				fontFamily: layer.settings.textFontFamily,
				fontSize: `${layer.settings.textFontSize}px`,
				fontWeight: layer.settings.textFontWeight,
				color: layer.settings.textColor,
				textAlign: layer.settings.textAlign,
				lineHeight: 1.2,
				whiteSpace: 'pre-wrap',
			}"
		>
			{{ layer.settings.textContent }}
		</div>
	</div>
</template>
