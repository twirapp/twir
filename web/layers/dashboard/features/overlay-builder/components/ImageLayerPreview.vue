<script setup lang="ts">
import { computed, ref, watch } from 'vue'

interface Props {
	imageUrl?: string
	width: number
	height: number
}

const props = defineProps<Props>()

const imageError = ref(false)
const imageLoaded = ref(false)

// Check if imageUrl is valid (not empty/null/undefined)
const hasValidUrl = computed(() => {
	return props.imageUrl && props.imageUrl.trim().length > 0
})

// Reset error state when URL changes
watch(() => props.imageUrl, (newUrl) => {
	console.log('[ImageLayerPreview] imageUrl changed:', newUrl, 'hasValidUrl:', hasValidUrl.value)
	imageError.value = false
	imageLoaded.value = false
})

// Log initial mount
console.log('[ImageLayerPreview] Component mounted with imageUrl:', props.imageUrl, 'hasValidUrl:', hasValidUrl.value)

const imageStyle = computed(() => ({
	width: '100%',
	height: '100%',
	objectFit: 'contain' as const,
}))

function handleImageLoad() {
	console.log('[ImageLayerPreview] Image loaded successfully:', props.imageUrl)
	imageLoaded.value = true
	imageError.value = false
}

function handleImageError() {
	console.error('[ImageLayerPreview] Failed to load image:', props.imageUrl)
	imageError.value = true
	imageLoaded.value = false
}
</script>

<template>
	<div class="w-full h-full flex items-center justify-center bg-slate-100 dark:bg-slate-800">
		<img
			v-if="hasValidUrl && !imageError"
			:src="imageUrl"
			:style="imageStyle"
			alt="Overlay image"
			@load="handleImageLoad"
			@error="handleImageError"
		/>
		<div v-else-if="imageError" class="text-center p-4">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-8 w-8 mx-auto mb-2 text-red-500"
				fill="none"
				viewBox="0 0 24 24"
				stroke="currentColor"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
				/>
			</svg>
			<p class="text-xs text-red-600 dark:text-red-400">Failed to load image</p>
			<p class="text-xs text-slate-500 mt-1">Check the image URL</p>
		</div>
		<div v-else class="text-center p-4">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-8 w-8 mx-auto mb-2 text-slate-400"
				fill="none"
				viewBox="0 0 24 24"
				stroke="currentColor"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
				/>
			</svg>
			<p class="text-xs text-slate-500">No image URL</p>
		</div>
	</div>
</template>
