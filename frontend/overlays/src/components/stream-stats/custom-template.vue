<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'

const props = defineProps<{
	html: string
	css: string
	values: Record<string, string>
}>()

const baseStyles = `
	* {
		box-sizing: border-box;
	}
`

const containerRef = ref<HTMLDivElement>()
const shadowRoot = ref<ShadowRoot>()

function renderContent() {
	if (!shadowRoot.value) return

	let html = props.html
	for (const [key, value] of Object.entries(props.values)) {
		html = html.replaceAll(`{{${key}}}`, value)
	}

	shadowRoot.value.innerHTML = `<style>${baseStyles}${props.css}</style>${html}`
}

watch(
	() => [props.html, props.css, props.values],
	() => {
		renderContent()
	},
	{ deep: true }
)

onMounted(() => {
	if (!containerRef.value || shadowRoot.value) return

	try {
		shadowRoot.value = containerRef.value.attachShadow({ mode: 'open' })
		renderContent()
	} catch (e) {
		console.error('Failed to create shadow DOM for stream stats custom template:', e)
	}
})

onUnmounted(() => {
	if (shadowRoot.value) {
		shadowRoot.value.innerHTML = ''
	}
})
</script>

<template>
	<div ref="containerRef" />
</template>
