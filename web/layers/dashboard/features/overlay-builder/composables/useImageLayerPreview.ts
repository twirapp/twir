import { type Ref, computed, onMounted, ref, watch } from 'vue'

export function useImageLayerPreview(imageUrl: Ref<string | undefined>) {
	const imageError = ref(false)
	const imageLoaded = ref(false)
	const hasValidUrl = computed(() => Boolean(imageUrl.value?.trim()))

	watch(imageUrl, (newUrl) => {
		console.log('[ImageLayerPreview] imageUrl changed:', newUrl, 'hasValidUrl:', hasValidUrl.value)
		imageError.value = false
		imageLoaded.value = false
	})
	onMounted(() => console.log('[ImageLayerPreview] Component mounted with imageUrl:', imageUrl.value, 'hasValidUrl:', hasValidUrl.value))

	const imageStyle = computed(() => ({
		width: '100%',
		height: '100%',
		objectFit: 'contain' as const,
	}))

	function handleImageLoad() {
		console.log('[ImageLayerPreview] Image loaded successfully:', imageUrl.value)
		imageLoaded.value = true
		imageError.value = false
	}

	function handleImageError() {
		console.error('[ImageLayerPreview] Failed to load image:', imageUrl.value)
		imageError.value = true
		imageLoaded.value = false
	}

	return { hasValidUrl, imageError, imageStyle, handleImageLoad, handleImageError }
}
