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
		style="position: absolute; overflow: hidden; 'text-wrap': 'nowrap'"
		:style="{
			top: layer.pos_y,
			left: layer.pos_x,
			width: layer.width,
			height: layer.height,
		}"
		:id="'layer' + layer.id"
		v-html="parsedData"
	/>
</template>
