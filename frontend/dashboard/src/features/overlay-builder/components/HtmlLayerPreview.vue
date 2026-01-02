<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

interface Props {
	html?: string
	css?: string
	js?: string
	width: number
	height: number
}

const props = withDefaults(defineProps<Props>(), {
	html: '',
	css: '',
	js: '',
})

const containerRef = ref<HTMLDivElement>()
const contentRef = ref<HTMLDivElement>()
const styleElement = ref<HTMLStyleElement>()
const renderKey = ref(0)

const sanitizedHtml = computed(() => {
	return props.html || '<div style="display: flex; align-items: center; justify-content: center; height: 100%; color: rgba(255,255,255,0.5); font-size: 14px;">Empty HTML Layer</div>'
})

// Apply CSS by injecting a style element
function updateStyles() {
	if (!containerRef.value) return

	// Remove old style element if exists
	if (styleElement.value) {
		styleElement.value.remove()
		styleElement.value = undefined
	}

	// Create new style element with scoped styles
	if (props.css) {
		const style = document.createElement('style')
		style.textContent = props.css
		styleElement.value = style
		containerRef.value.appendChild(style)
	}
}

// Execute JavaScript
function executeScript() {
	if (!props.js) return

	try {
		// Create a safer execution context
		const scriptFunc = new Function('container', props.js)
		scriptFunc(contentRef.value)
	} catch (e) {
		console.error('Layer JS Error:', e)
	}
}

// Force re-render when content changes
function forceUpdate() {
	renderKey.value++
	// Wait for DOM update before applying styles and scripts
	setTimeout(() => {
		updateStyles()
		executeScript()
	}, 0)
}

// Watch for prop changes
watch(() => props.html, () => {
	forceUpdate()
})

watch(() => props.css, () => {
	updateStyles()
})

watch(() => props.js, () => {
	executeScript()
})

onMounted(() => {
	updateStyles()
	executeScript()
})

onUnmounted(() => {
	if (styleElement.value) {
		styleElement.value.remove()
	}
})
</script>

<template>
	<div
		ref="containerRef"
		class="html-layer-preview w-full h-full relative overflow-hidden"
		:style="{
			width: props.width + 'px',
			height: props.height + 'px',
		}"
	>
		<div
			:key="renderKey"
			ref="contentRef"
			class="html-content w-full h-full"
			v-html="sanitizedHtml"
		/>
	</div>
</template>

<style scoped>
.html-layer-preview {
	background: transparent;
	color: #fff;
	font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.html-layer-preview :deep(*) {
	box-sizing: border-box;
}

.html-content {
	pointer-events: none;
}
</style>
