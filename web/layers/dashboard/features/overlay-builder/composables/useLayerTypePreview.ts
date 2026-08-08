import { generateFontKey, useFontSource } from '@twir/fontsource'
import { type CSSProperties, type Ref, computed, onMounted, ref, watch } from 'vue'

import type { Layer } from '../types'

export function useLayerTypePreview(layer: Ref<Layer>) {
	const fontSource = useFontSource(false)
	const loadedFontKey = ref('')

	// document.fonts is client-only, hence the onMounted wrapper
	onMounted(() => {
		watch(
			[
				() => layer.value.settings.textFontFamily,
				() => layer.value.settings.textFontWeight,
				() => layer.value.settings.textFontStyle,
			],
			async ([fontFamily, fontWeight, fontStyle]) => {
				loadedFontKey.value = ''
				if (!fontFamily) return

				const style = fontStyle === 'italic' ? 'italic' : 'normal'
				const font = await fontSource.loadFont(fontFamily, fontWeight, style)
				if (!font) return

				loadedFontKey.value = generateFontKey(fontFamily, fontWeight, style)
			},
			{ immediate: true }
		)
	})

	const fontFamily = computed(() => {
		if (loadedFontKey.value) {
			return `"${loadedFontKey.value}", sans-serif`
		}

		return layer.value.settings.textFontFamily || 'sans-serif'
	})

	const textContainerStyle = computed<CSSProperties>(() => {
		const vertical = layer.value.settings.textAlignVertical

		return {
			justifyContent:
				vertical === 'center' ? 'center' : vertical === 'bottom' ? 'flex-end' : 'flex-start',
		}
	})

	const textStyle = computed<CSSProperties>(() => {
		const settings = layer.value.settings

		const strokeWidth = settings.textStrokeWidth
		const shadowBlur = settings.textShadowBlur
		const shadowX = settings.textShadowOffsetX
		const shadowY = settings.textShadowOffsetY

		return {
			color: settings.textColor,
			fontFamily: fontFamily.value,
			fontSize: `${settings.textFontSize}px`,
			fontWeight: settings.textFontWeight,
			fontStyle: settings.textFontStyle === 'italic' ? 'italic' : 'normal',
			textAlign: settings.textAlign as CSSProperties['textAlign'],
			lineHeight: settings.textLineHeight || 1.2,
			letterSpacing: `${settings.textLetterSpacing || 0}px`,
			textTransform: (settings.textTransform || 'none') as CSSProperties['textTransform'],
			WebkitTextStroke: strokeWidth
				? `${strokeWidth}px ${settings.textStrokeColor || '#000'}`
				: undefined,
			textShadow:
				shadowBlur || shadowX || shadowY
					? `${shadowX}px ${shadowY}px ${shadowBlur}px ${settings.textShadowColor || 'rgba(0, 0, 0, 0.8)'}`
					: undefined,
			whiteSpace: 'pre-wrap' as const,
			wordBreak: 'break-word' as const,
		}
	})

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

	return { textStyle, textContainerStyle, iframeStyle, youtubeEmbedUrl }
}
