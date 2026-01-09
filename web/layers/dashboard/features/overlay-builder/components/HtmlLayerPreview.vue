<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

import { useChannelOverlayParseHtml } from '#layers/dashboard/api/overlays/custom'

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
const shadowRoot = ref<ShadowRoot>()
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

// Render content in shadow DOM
function renderInShadowDOM() {
	if (!shadowRoot.value) return

	// Create base styles for shadow DOM
	const baseStyles = `
		* {
			box-sizing: border-box;
		}
		:host {
			display: block;
			width: 100%;
			height: 100%;
			overflow: hidden;
			background: transparent;
			color: #fff;
			font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
		}
		.html-content {
			width: 100%;
			height: 100%;
			pointer-events: none;
		}
	`

	// Combine base styles with user CSS
	const combinedStyles = baseStyles + (props.css || '')

	// Build shadow DOM content
	shadowRoot.value.innerHTML = `
		<style>${combinedStyles}</style>
		<div class="html-content">${sanitizedHtml.value}</div>
	`

	// Execute user JavaScript after DOM is ready
	setTimeout(() => {
		executeScript()
	}, 0)
}

// Execute JavaScript in shadow DOM context
function executeScript() {
	if (!props.js || !shadowRoot.value) return

	try {
		const contentElement = shadowRoot.value.querySelector('.html-content')
		if (!contentElement) return

		// Create a function that has access to the content element
		// eslint-disable-next-line no-new-func
		const scriptFunc = new Function('container', props.js)
		scriptFunc(contentElement)
	} catch (e) {
		console.error('Canvas Layer JS Error:', e)
	}
}

// Force re-render when content changes
function forceUpdate() {
	renderKey.value++
	setTimeout(() => {
		renderInShadowDOM()
	}, 0)
}

// Initialize shadow DOM
function initializeShadowDOM() {
	if (!containerRef.value) return

	// Create shadow root if it doesn't exist
	if (!shadowRoot.value) {
		try {
			shadowRoot.value = containerRef.value.attachShadow({ mode: 'open' })
		} catch (e) {
			console.error('Failed to create shadow DOM:', e)
			return
		}
	}

	renderInShadowDOM()
}

// Watch for prop changes
watch(() => props.html, () => {
	forceUpdate()
	startPolling()
})

watch(() => props.css, () => {
	renderInShadowDOM()
})

watch(() => props.js, () => {
	renderInShadowDOM()
})

watch(() => props.refreshInterval, () => {
	startPolling()
})

watch(() => sanitizedHtml.value, () => {
	renderInShadowDOM()
})

onMounted(() => {
	initializeShadowDOM()
	startPolling()
})

onUnmounted(() => {
	stopPolling()
	if (shadowRoot.value) {
		shadowRoot.value.innerHTML = ''
	}
})
</script>

<template>
	<div
		:key="renderKey"
		ref="containerRef"
		class="html-layer-preview"
		:style="{
			width: props.width + 'px',
			height: props.height + 'px',
		}"
	/>
</template>

<style scoped>
.html-layer-preview {
	display: block;
	overflow: hidden;
	background: transparent;
}
</style>
