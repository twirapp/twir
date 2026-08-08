<script setup lang="ts">
import { computed, ref, toRef, watch } from 'vue'

import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

import EmoteSettingsEditor from './EmoteSettingsEditor.vue'
import IframeSettingsEditor from './IframeSettingsEditor.vue'
import type { Layer, LayerSettings } from '../../types'
import FieldInput from '../fields/FieldInput.vue'
import FieldNumber from '../fields/FieldNumber.vue'
import FieldSelect from '../fields/FieldSelect.vue'
import FieldSwitch from '../fields/FieldSwitch.vue'
import { useLayerSettingsEditor } from '../../composables/useLayerSettingsEditor'
import { type Font, FontSelector } from '~~/layers/dashboard/lib/fontsource'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()
const { t } = useI18n()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
}>()

const layer = toRef(props, 'layer')
const fieldId = (name: string) => `layer-${layer.value.id}-${name}`

function updateSettings(updates: Partial<LayerSettings>) {
	emit('update', { settings: { ...layer.value.settings, ...updates } })
}

const {
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
} = useLayerSettingsEditor(layer, updateSettings)

const fontData = ref<Font | null>(null)
watch(fontData, (font) => {
	if (!font) return

	const updates: Partial<LayerSettings> = {}
	if (font.id !== layer.value.settings.textFontFamily) {
		updates.textFontFamily = font.id
	}
	if (!font.weights.includes(layer.value.settings.textFontWeight)) {
		updates.textFontWeight = font.weights.includes(400) ? 400 : font.weights[0]
	}
	const currentStyle = layer.value.settings.textFontStyle === 'italic' ? 'italic' : 'normal'
	if (!font.styles.includes(currentStyle)) {
		updates.textFontStyle = font.styles.includes('normal') ? 'normal' : font.styles[0]
	}

	if (Object.keys(updates).length > 0) {
		updateSettings(updates)
	}
})

const fontWeightOptions = computed(() => {
	const weights = fontData.value?.weights ?? [100, 200, 300, 400, 500, 600, 700, 800, 900]
	return weights.map((weight) => ({ value: String(weight), label: String(weight) }))
})
const fontStyleOptions = computed(() => {
	const styles = fontData.value?.styles ?? ['normal', 'italic']
	return styles.map((style) => ({
		value: style,
		labelKey:
			style === 'italic'
				? 'overlayBuilder.editors.text.styleItalic'
				: 'overlayBuilder.editors.text.styleNormal',
	}))
})
const textAlignOptions = [
	{ value: 'left', labelKey: 'overlayBuilder.editors.text.alignLeft' },
	{ value: 'center', labelKey: 'overlayBuilder.editors.text.alignCenter' },
	{ value: 'right', labelKey: 'overlayBuilder.editors.text.alignRight' },
] as const
const textAlignVerticalOptions = [
	{ value: 'top', labelKey: 'overlayBuilder.editors.text.alignTop' },
	{ value: 'center', labelKey: 'overlayBuilder.editors.text.alignMiddle' },
	{ value: 'bottom', labelKey: 'overlayBuilder.editors.text.alignBottom' },
] as const
const textTransformOptions = [
	{ value: 'none', labelKey: 'overlayBuilder.editors.text.transformNone' },
	{ value: 'uppercase', labelKey: 'overlayBuilder.editors.text.transformUppercase' },
	{ value: 'lowercase', labelKey: 'overlayBuilder.editors.text.transformLowercase' },
	{ value: 'capitalize', labelKey: 'overlayBuilder.editors.text.transformCapitalize' },
] as const
</script>

<template>
	<div v-if="layer.type === 'TEXT'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">{{ t('overlayBuilder.editors.text.title') }}</h4>

		<div class="flex flex-col gap-2">
			<Label :for="fieldId('text-content')">{{ t('overlayBuilder.editors.text.content') }}</Label>
			<Textarea :id="fieldId('text-content')" v-model="textContent" rows="3" @keydown.stop />
		</div>

		<h4 class="text-sm font-medium">{{ t('overlayBuilder.editors.text.font') }}</h4>
		<FontSelector
			v-model:font="fontData"
			:font-family="textFontFamily"
			:font-weight="layer.settings.textFontWeight"
			:font-style="layer.settings.textFontStyle"
		/>
		<div class="grid grid-cols-2 gap-3">
			<FieldNumber :id="fieldId('text-size')" v-model="textFontSize" :label="t('overlayBuilder.editors.text.size')" :min="1" :step="1" />
			<FieldSelect :id="fieldId('text-weight')" v-model="textFontWeight" :label="t('overlayBuilder.editors.text.weight')" :options="fontWeightOptions" />
		</div>
		<div class="grid grid-cols-2 gap-3">
			<FieldSelect :id="fieldId('text-style')" v-model="textFontStyle" :label="t('overlayBuilder.editors.text.fontStyle')" :options="fontStyleOptions.map((option) => ({ value: option.value, label: t(option.labelKey) }))" />
			<FieldInput :id="fieldId('text-color')" v-model="textColor" :label="t('overlayBuilder.editors.text.color')" type="color" input-class="h-9 p-1" />
		</div>

		<h4 class="text-sm font-medium">{{ t('overlayBuilder.editors.text.layout') }}</h4>
		<div class="grid grid-cols-2 gap-3">
			<FieldSelect :id="fieldId('text-align')" v-model="textAlign" :label="t('overlayBuilder.editors.text.align')" :options="textAlignOptions.map((option) => ({ value: option.value, label: t(option.labelKey) }))" />
			<FieldSelect :id="fieldId('text-align-vertical')" v-model="textAlignVertical" :label="t('overlayBuilder.editors.text.verticalAlign')" :options="textAlignVerticalOptions.map((option) => ({ value: option.value, label: t(option.labelKey) }))" />
		</div>
		<div class="grid grid-cols-2 gap-3">
			<FieldNumber :id="fieldId('text-line-height')" v-model="textLineHeight" :label="t('overlayBuilder.editors.text.lineHeight')" :min="0.5" :max="5" :step="0.1" />
			<FieldNumber :id="fieldId('text-letter-spacing')" v-model="textLetterSpacing" :label="t('overlayBuilder.editors.text.letterSpacing')" :step="0.5" />
		</div>
		<FieldSelect :id="fieldId('text-transform')" v-model="textTransform" :label="t('overlayBuilder.editors.text.transform')" :options="textTransformOptions.map((option) => ({ value: option.value, label: t(option.labelKey) }))" />

		<h4 class="text-sm font-medium">{{ t('overlayBuilder.editors.text.stroke') }}</h4>
		<div class="grid grid-cols-2 gap-3">
			<FieldNumber :id="fieldId('text-stroke-width')" v-model="textStrokeWidth" :label="t('overlayBuilder.editors.text.strokeWidth')" :min="0" :step="0.5" />
			<FieldInput :id="fieldId('text-stroke-color')" v-model="textStrokeColor" :label="t('overlayBuilder.editors.text.strokeColor')" type="color" input-class="h-9 p-1" />
		</div>

		<h4 class="text-sm font-medium">{{ t('overlayBuilder.editors.text.shadow') }}</h4>
		<div class="grid grid-cols-2 gap-3">
			<FieldInput :id="fieldId('text-shadow-color')" v-model="textShadowColor" :label="t('overlayBuilder.editors.text.shadowColor')" type="color" input-class="h-9 p-1" />
			<FieldNumber :id="fieldId('text-shadow-blur')" v-model="textShadowBlur" :label="t('overlayBuilder.editors.text.shadowBlur')" :min="0" :step="1" />
		</div>
		<div class="grid grid-cols-2 gap-3">
			<FieldNumber :id="fieldId('text-shadow-offset-x')" v-model="textShadowOffsetX" :label="t('overlayBuilder.editors.text.shadowOffsetX')" :step="1" />
			<FieldNumber :id="fieldId('text-shadow-offset-y')" v-model="textShadowOffsetY" :label="t('overlayBuilder.editors.text.shadowOffsetY')" :step="1" />
		</div>
	</div>

	<div v-else-if="layer.type === 'VIDEO'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">{{ t('overlayBuilder.editors.video.title') }}</h4>
		<FieldInput :id="fieldId('video-url')" v-model="videoUrl" :label="t('overlayBuilder.editors.video.url')" type="url" placeholder="https://example.com/video.mp4" />
		<FieldSwitch :id="fieldId('video-loop')" v-model="videoLoop" :label="t('overlayBuilder.editors.video.loop')" />
		<FieldSwitch :id="fieldId('video-muted')" v-model="videoMuted" :label="t('overlayBuilder.editors.video.muted')" />
	</div>

	<div v-else-if="layer.type === 'IFRAME'">
		<IframeSettingsEditor :layer="layer" @update="emit('update', $event)" />
	</div>

	<div v-else-if="layer.type === 'YOUTUBE'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">YouTube</h4>
		<FieldInput :id="fieldId('youtube-id')" v-model="youtubeVideoId" :label="t('overlayBuilder.editors.youtube.videoLink')" :placeholder="t('overlayBuilder.editors.youtube.videoPlaceholder')" @blur="youtubeVideoId = parseYoutubeVideoId(youtubeVideoId)" />
		<FieldSwitch :id="fieldId('youtube-autoplay')" v-model="youtubeAutoplay" :label="t('overlayBuilder.editors.youtube.autoplay')" />
		<FieldSwitch :id="fieldId('youtube-loop')" v-model="youtubeLoop" :label="t('overlayBuilder.editors.youtube.loop')" />
		<FieldSwitch :id="fieldId('youtube-muted')" v-model="youtubeMuted" :label="t('overlayBuilder.editors.youtube.muted')" />
	</div>

	<div v-else-if="layer.type === 'EMOTE'">
		<EmoteSettingsEditor :layer="layer" @update="emit('update', $event)" />
	</div>
</template>
