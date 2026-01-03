<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'

import type { Layer } from '@/composables/overlays/use-overlays.js'

const props = defineProps<{
	layer: Layer
	parsedData?: string
}>()

const containerRef = ref<HTMLDivElement>()
const shadowRoot = ref<ShadowRoot>()

// Initialize shadow DOM
function initializeShadowDOM() {
	if (!containerRef.value || shadowRoot.value) return

	try {
		shadowRoot.value = containerRef.value.attachShadow({ mode: 'open' })
		renderContent()
	} catch (e) {
		console.error('Failed to create shadow DOM for layer:', props.layer.id, e)
	}
}

// Render content in shadow DOM
function renderContent() {
	if (!shadowRoot.value) return

	// Base styles for shadow DOM
	const baseStyles = `
		* {
			box-sizing: border-box;
		}
		:host {
			display: block;
			width: 100%;
			height: 100%;
			overflow: hidden;
		}
		.layer-content {
			width: 100%;
			height: 100%;
			white-space: nowrap;
		}
	`

	// Combine base styles with user CSS
	const combinedStyles = baseStyles + (props.layer.settings.htmlOverlayCss || '')

	// Build shadow DOM content
	const htmlContent = props.parsedData || ''
	shadowRoot.value.innerHTML = `
		<style>${combinedStyles}</style>
		<div class="layer-content">${htmlContent}</div>
	`

	// Execute user JavaScript after DOM is ready
	setTimeout(() => {
		executeScript()
	}, 0)
}

// Execute JavaScript in shadow DOM context
function executeScript() {
	if (!props.layer.settings.htmlOverlayJs || !shadowRoot.value) return

	try {
		const contentElement = shadowRoot.value.querySelector('.layer-content')
		if (!contentElement) return

		// Create onDataUpdate function that user code can call
		const onDataUpdate = () => {
			// User's onDataUpdate callback
		}

		// Execute user's JavaScript with access to container and onDataUpdate
		// eslint-disable-next-line no-new-func
		const scriptFunc = new Function('container', 'onDataUpdate', `
			${props.layer.settings.htmlOverlayJs}
			if (typeof onDataUpdate === 'function') {
				onDataUpdate();
			}
		`)
		scriptFunc(contentElement, onDataUpdate)
	} catch (e) {
		console.error('Layer JS Error:', props.layer.id, e)
	}
}

// Watch for data updates
watch(
	() => props.parsedData,
	() => {
		renderContent()
	}
)

// Watch for settings changes
watch(
	() => [props.layer.settings.htmlOverlayCss, props.layer.settings.htmlOverlayJs],
	() => {
		renderContent()
	},
	{ deep: true }
)

onMounted(() => {
	initializeShadowDOM()
})

onUnmounted(() => {
	if (shadowRoot.value) {
		shadowRoot.value.innerHTML = ''
	}
})
</script>

<template>
	<div
		ref="containerRef"
		:id="'layer' + layer.id"
		style="position: absolute; overflow: hidden;"
		:style="{
			top: `${layer.posY}px`,
			left: `${layer.posX}px`,
			width: `${layer.width}px`,
			height: `${layer.height}px`,
			transform: `rotate(${layer.rotation || 0}deg)`,
			transformOrigin: 'center center',
		}"
	/>
</template>
