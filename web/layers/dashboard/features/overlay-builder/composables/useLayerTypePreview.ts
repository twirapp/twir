import { type CSSProperties, type Ref, computed } from 'vue'

import type { Layer } from '../types'

export function useLayerTypePreview(layer: Ref<Layer>) {
	const textStyle = computed<CSSProperties>(() => ({
		color: layer.value.settings.textColor,
		fontFamily: layer.value.settings.textFontFamily,
		fontSize: `${layer.value.settings.textFontSize}px`,
		fontWeight: layer.value.settings.textFontWeight,
		textAlign: layer.value.settings.textAlign as CSSProperties['textAlign'],
		whiteSpace: 'pre-wrap' as const,
		wordBreak: 'break-word' as const,
	}))

	const iframeStyle = computed(() => {
		const scale = layer.value.settings.iframeScale || 1
		return {
			width: `${100 / scale}%`,
			height: `${100 / scale}%`,
			transform: `scale(${scale})`,
			transformOrigin: 'top left',
		}
	})

	const youtubeEmbedUrl = computed(() => {
		const videoId = layer.value.settings.youtubeVideoId.trim()
		if (!videoId) return ''

		const params = new URLSearchParams({
			autoplay: layer.value.settings.youtubeAutoplay ? '1' : '0',
			mute: layer.value.settings.youtubeMuted ? '1' : '0',
			loop: layer.value.settings.youtubeLoop ? '1' : '0',
			controls: '0',
		})
		if (layer.value.settings.youtubeLoop) params.set('playlist', videoId)

		return `https://www.youtube-nocookie.com/embed/${encodeURIComponent(videoId)}?${params.toString()}`
	})

	return { textStyle, iframeStyle, youtubeEmbedUrl }
}
