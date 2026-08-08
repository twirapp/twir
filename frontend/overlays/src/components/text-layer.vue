<script setup lang="ts">
import { generateFontKey, useFontSource } from '@twir/fontsource'
import { computed, ref, watch } from 'vue'

import type { Layer } from '@/composables/overlays/use-overlays.js'

const props = defineProps<{
	layer: Layer
	zIndex?: number
}>()

const fontSource = useFontSource(false)

const fontStyle = computed(() => {
	return props.layer.settings.textFontStyle === 'italic' ? 'italic' : 'normal'
})

const loadedFontKey = ref('')

watch(
	[() => props.layer.settings.textFontFamily, () => props.layer.settings.textFontWeight, fontStyle],
	async ([fontFamily, fontWeight, style]) => {
		loadedFontKey.value = ''
		if (!fontFamily) return

		const font = await fontSource.loadFont(fontFamily, fontWeight, style)
		if (!font) return

		loadedFontKey.value = generateFontKey(fontFamily, fontWeight, style)
	},
	{ immediate: true }
)

const fontFamily = computed(() => {
	if (loadedFontKey.value) {
		return `"${loadedFontKey.value}", sans-serif`
	}

	return props.layer.settings.textFontFamily || 'sans-serif'
})

const justifyContent = computed(() => {
	switch (props.layer.settings.textAlignVertical) {
		case 'center':
			return 'center'
		case 'bottom':
			return 'flex-end'
		default:
			return 'flex-start'
	}
})

const textStroke = computed(() => {
	const { textStrokeWidth, textStrokeColor } = props.layer.settings
	if (!textStrokeWidth) return undefined

	return `${textStrokeWidth}px ${textStrokeColor || '#000'}`
})

const textShadow = computed(() => {
	const { textShadowBlur, textShadowOffsetX, textShadowOffsetY, textShadowColor } =
		props.layer.settings
	if (!textShadowBlur && !textShadowOffsetX && !textShadowOffsetY) return undefined

	return `${textShadowOffsetX}px ${textShadowOffsetY}px ${textShadowBlur}px ${textShadowColor || 'rgba(0, 0, 0, 0.8)'}`
})
</script>

<template>
	<div
		:id="'layer' + layer.id"
		style="position: absolute; overflow: hidden"
		:style="{
			top: `${layer.posY}px`,
			left: `${layer.posX}px`,
			width: `${layer.width}px`,
			height: `${layer.height}px`,
			transform: `rotate(${layer.rotation || 0}deg)`,
			transformOrigin: 'center center',
			zIndex: zIndex ?? 0,
		}"
	>
		<div
			:style="{
				display: 'flex',
				flexDirection: 'column',
				justifyContent,
				width: '100%',
				height: '100%',
			}"
		>
			<div
				:style="{
					fontFamily,
					fontSize: `${layer.settings.textFontSize}px`,
					fontWeight: layer.settings.textFontWeight,
					fontStyle,
					color: layer.settings.textColor,
					textAlign: layer.settings.textAlign,
					lineHeight: layer.settings.textLineHeight || 1.2,
					letterSpacing: `${layer.settings.textLetterSpacing || 0}px`,
					textTransform: layer.settings.textTransform || 'none',
					WebkitTextStroke: textStroke,
					textShadow,
					whiteSpace: 'pre-wrap',
				}"
			>
				{{ layer.settings.textContent }}
			</div>
		</div>
	</div>
</template>
