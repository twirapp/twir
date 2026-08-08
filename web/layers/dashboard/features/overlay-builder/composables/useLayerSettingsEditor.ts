import { type Ref, computed } from 'vue'

import type { Layer, LayerSettings } from '../types'

type UpdateSettings = (updates: Partial<LayerSettings>) => void

export function useLayerSettingsEditor(layer: Ref<Layer>, updateSettings: UpdateSettings) {
	const textContent = computed({
		get: () => layer.value.settings.textContent,
		set: (value: string | number) => updateSettings({ textContent: String(value) }),
	})

	const textFontSize = computed({
		get: () => layer.value.settings.textFontSize,
		set: (value: string | number) => {
			const parsed = Number(value)
			if (Number.isFinite(parsed)) updateSettings({ textFontSize: Math.max(1, Math.round(parsed)) })
		},
	})

	const textFontWeight = computed({
		get: () => String(layer.value.settings.textFontWeight),
		set: (value: string) => updateSettings({ textFontWeight: Number(value) }),
	})

	const textColor = computed({
		get: () => layer.value.settings.textColor,
		set: (value: string) => updateSettings({ textColor: value }),
	})

	const textAlign = computed({
		get: () => layer.value.settings.textAlign,
		set: (value: string) => updateSettings({ textAlign: value }),
	})

	const textFontFamily = computed({
		get: () => layer.value.settings.textFontFamily,
		set: (value: string) => updateSettings({ textFontFamily: value }),
	})

	const textFontStyle = computed({
		get: () => layer.value.settings.textFontStyle,
		set: (value: string) => updateSettings({ textFontStyle: value }),
	})

	const textAlignVertical = computed({
		get: () => layer.value.settings.textAlignVertical,
		set: (value: string) => updateSettings({ textAlignVertical: value }),
	})

	const textStrokeWidth = computed({
		get: () => layer.value.settings.textStrokeWidth,
		set: (value: string | number) => {
			const parsed = Number(value)
			if (Number.isFinite(parsed)) updateSettings({ textStrokeWidth: Math.max(0, parsed) })
		},
	})

	const textStrokeColor = computed({
		get: () => layer.value.settings.textStrokeColor,
		set: (value: string) => updateSettings({ textStrokeColor: value }),
	})

	const textShadowColor = computed({
		get: () => layer.value.settings.textShadowColor,
		set: (value: string) => updateSettings({ textShadowColor: value }),
	})

	const textShadowBlur = computed({
		get: () => layer.value.settings.textShadowBlur,
		set: (value: string | number) => {
			const parsed = Number(value)
			if (Number.isFinite(parsed))
				updateSettings({ textShadowBlur: Math.max(0, Math.round(parsed)) })
		},
	})

	const textShadowOffsetX = computed({
		get: () => layer.value.settings.textShadowOffsetX,
		set: (value: string | number) => {
			const parsed = Number(value)
			if (Number.isFinite(parsed)) updateSettings({ textShadowOffsetX: Math.round(parsed) })
		},
	})

	const textShadowOffsetY = computed({
		get: () => layer.value.settings.textShadowOffsetY,
		set: (value: string | number) => {
			const parsed = Number(value)
			if (Number.isFinite(parsed)) updateSettings({ textShadowOffsetY: Math.round(parsed) })
		},
	})

	const textLineHeight = computed({
		get: () => layer.value.settings.textLineHeight,
		set: (value: string | number) => {
			const parsed = Number(value)
			if (Number.isFinite(parsed))
				updateSettings({ textLineHeight: Math.min(5, Math.max(0.5, parsed)) })
		},
	})

	const textLetterSpacing = computed({
		get: () => layer.value.settings.textLetterSpacing,
		set: (value: string | number) => {
			const parsed = Number(value)
			if (Number.isFinite(parsed)) updateSettings({ textLetterSpacing: parsed })
		},
	})

	const textTransform = computed({
		get: () => layer.value.settings.textTransform,
		set: (value: string) => updateSettings({ textTransform: value }),
	})

	const videoUrl = computed({
		get: () => layer.value.settings.videoUrl,
		set: (value: string) => updateSettings({ videoUrl: value }),
	})

	const videoLoop = computed({
		get: () => layer.value.settings.videoLoop,
		set: (value: boolean) => updateSettings({ videoLoop: value }),
	})

	const videoMuted = computed({
		get: () => layer.value.settings.videoMuted,
		set: (value: boolean) => updateSettings({ videoMuted: value }),
	})

	const youtubeVideoId = computed({
		get: () => layer.value.settings.youtubeVideoId,
		set: (value: string) => updateSettings({ youtubeVideoId: value }),
	})

	const youtubeAutoplay = computed({
		get: () => layer.value.settings.youtubeAutoplay,
		set: (value: boolean) => updateSettings({ youtubeAutoplay: value }),
	})

	const youtubeLoop = computed({
		get: () => layer.value.settings.youtubeLoop,
		set: (value: boolean) => updateSettings({ youtubeLoop: value }),
	})

	const youtubeMuted = computed({
		get: () => layer.value.settings.youtubeMuted,
		set: (value: boolean) => updateSettings({ youtubeMuted: value }),
	})

	function parseYoutubeVideoId(value: string): string {
		const trimmed = value.trim()
		if (!trimmed) return ''

		const match = trimmed.match(
			/(?:youtube\.com\/watch\?[^#\s]*?v=|youtu\.be\/|youtube\.com\/(?:embed|shorts|live)\/)([^?&#/\s]+)/i
		)

		return match?.[1] ?? trimmed
	}

	return {
		textContent,
		textFontSize,
		textFontWeight,
		textFontStyle,
		textColor,
		textAlign,
		textAlignVertical,
		textFontFamily,
		textStrokeWidth,
		textStrokeColor,
		textShadowColor,
		textShadowBlur,
		textShadowOffsetX,
		textShadowOffsetY,
		textLineHeight,
		textLetterSpacing,
		textTransform,
		videoUrl,
		videoLoop,
		videoMuted,
		youtubeVideoId,
		youtubeAutoplay,
		youtubeLoop,
		youtubeMuted,
		parseYoutubeVideoId,
	}
}
