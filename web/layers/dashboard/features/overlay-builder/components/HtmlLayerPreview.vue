<script setup lang="ts">
import { type HtmlLayerPreviewProps, useHtmlLayerPreview } from '../composables/useHtmlLayerPreview'

const props = withDefaults(defineProps<HtmlLayerPreviewProps>(), {
	html: '',
	css: '',
	js: '',
	refreshInterval: 5,
})

const { containerRef, renderKey } = useHtmlLayerPreview(props)
</script>

<template>
	<div
		:key="renderKey"
		ref="containerRef"
		class="html-layer-preview"
		:style="{ width: `${props.width}px`, height: `${props.height}px` }"
	/>
</template>

<style scoped>
.html-layer-preview {
	display: block;
	overflow: hidden;
	background: transparent;
}
</style>
