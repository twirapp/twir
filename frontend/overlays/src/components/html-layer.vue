<script setup lang="ts">
import { transform } from 'nested-css-to-flat'
import { computed, nextTick, watch } from 'vue'

import type { Layer } from '@/composables/overlays/use-overlays.js'

const props = defineProps<{
	layer: Layer
	parsedData?: string
}>()

const executeFunc = computed(() => {
	// oxlint-disable-next-line no-new-func
	return new Function(`${props.layer.settings.htmlOverlayJs}; onDataUpdate();`)
})

watch(
	() => props.parsedData,
	async () => {
		await nextTick()
		executeFunc.value?.()
	}
)
</script>

<template>
	<component :is="'style'">
		{{
			transform(`#layer${layer.id} {
					${layer.settings.htmlOverlayCss}
				}`)
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
