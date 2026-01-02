<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

import { useChannelOverlayParseHtml } from '@/api/overlays/custom'

interface Props {
	html?: string
	css?: string
	js?: string
	width: number
	height: number
	refreshInterval?: number
}

const props = withDefaults(defineProps<Props>(), {
	html: '',
	css: '',
	js: '',
	refreshInterval: 5,
})

const containerRef = ref<HTMLDivElement>()
const contentRef = ref<HTMLDivElement>()
const styleElement = ref<HTMLStyleElement>()
const renderKey = ref(0)
const parsedHtml = ref('')
const pollInterval = ref<ReturnType<typeof setInterval>>()

const parseHtmlMutation = useChannelOverlayParseHtml()

const sanitizedHtml = computed(() => {
	// Use parsed HTML if available, otherwise use raw HTML
	const html = parsedHtml.value || props.html
	return html || '<div style="display: flex; align-items: center; justify-content: center; height: 100%; color: rgba(255,255,255,0.5); font-size: 14px;">Empty HTML Layer</div>'
})

// Parse HTML with variables
async function parseHtml() {
	if (!props.html) {
		parsedHtml.value = ''
		return
	}

	try {
		const result = await parseHtmlMutation.executeMutation({ html: props.html })
		parsedHtml.value = result.data?.channelOverlayParseHtml ?? props.html

		// Call onDataUpdate after parsing
		executeScript()
	} catch (e) {
		console.error('Canvas Layer: Failed to parse HTML:', e)
		parsedHtml.value = props.html
	}
}

// Start periodic polling
function startPolling() {
	stopPolling()

	// Initial parse
	parseHtml()

	// Set up interval
	if (props.refreshInterval && props.refreshInterval > 0) {
		pollInterval.value = setInterval(() => {
			parseHtml()
		}, props.refreshInterval * 1000)
	}
}

// Stop polling
function stopPolling() {
	if (pollInterval.value) {
		clearInterval(pollInterval.value)
		pollInterval.value = undefined
	}
}

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
		// eslint-disable-next-line no-new-func
		const scriptFunc = new Function('container', props.js)
		scriptFunc(contentRef.value)
	} catch (e) {
		console.error('Canvas Layer JS Error:', e)
	}
}

// Force re-render when content changes
function forceUpdate() {
	renderKey.value++
	// Wait for DOM update before applying styles and scripts
	setTimeout(() => {
		updateStyles()
	}, 0)
}

// Watch for prop changes
watch(() => props.html, () => {
	forceUpdate()
	startPolling()
})

watch(() => props.css, () => {
	updateStyles()
})

watch(() => props.js, () => {
	executeScript()
})

watch(() => props.refreshInterval, () => {
	startPolling()
})

onMounted(() => {
	updateStyles()
	startPolling()
})

onUnmounted(() => {
	stopPolling()
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
