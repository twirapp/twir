<script setup lang="ts">
import { transform } from 'nested-css-to-flat';

import { Layer } from '../sockets/overlays';

defineProps<{
	layer: Layer
	parsedData?: string
}>();
</script>

<template>
	<component :is="'style'">
		{{
			transform(`#layer${layer.id} {
					${layer.settings.htmlOverlayCss}
				}`
			)
		}}
	</component>
	<div
		:id="'layer' + layer.id"
		style="position: absolute; overflow: hidden; 'text-wrap': 'nowrap'"
		:style="{
			top: `${layer.pos_y}px`,
			left: `${layer.pos_x}px`,
			width: `${layer.width}px`,
			height: `${layer.height}px`,
		}"
		v-html="parsedData"
	/>
</template>
