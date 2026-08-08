<script setup lang="ts">
import { toRef } from 'vue'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'

import EmoteSettingsEditor from './EmoteSettingsEditor.vue'
import IframeSettingsEditor from './IframeSettingsEditor.vue'
import type { Layer, LayerSettings } from '../../types'
import FieldInput from '../fields/FieldInput.vue'
import FieldNumber from '../fields/FieldNumber.vue'
import FieldSelect from '../fields/FieldSelect.vue'
import FieldSwitch from '../fields/FieldSwitch.vue'
import { useLayerSettingsEditor } from '../../composables/useLayerSettingsEditor'

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
	textColor,
	textAlign,
	textFontFamily,
	videoUrl,
	videoLoop,
	videoMuted,
	youtubeVideoId,
	youtubeAutoplay,
	youtubeLoop,
	youtubeMuted,
	parseYoutubeVideoId,
} = useLayerSettingsEditor(layer, updateSettings)

const fontWeightOptions = [
	{ value: '400', label: '400' },
	{ value: '500', label: '500' },
	{ value: '600', label: '600' },
	{ value: '700', label: '700' },
	{ value: '800', label: '800' },
] as const
const textAlignOptions = [
	{ value: 'left', labelKey: 'overlayBuilder.editors.text.alignLeft' },
	{ value: 'center', labelKey: 'overlayBuilder.editors.text.alignCenter' },
	{ value: 'right', labelKey: 'overlayBuilder.editors.text.alignRight' },
] as const
</script>

<template>
	<div v-if="layer.type === 'TEXT'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">{{ t('overlayBuilder.editors.text.title') }}</h4>

		<div class="flex flex-col gap-2">
			<Label :for="fieldId('text-content')">{{ t('overlayBuilder.editors.text.content') }}</Label>
			<Textarea :id="fieldId('text-content')" v-model="textContent" rows="3" @keydown.stop />
		</div>

		<div class="grid grid-cols-2 gap-3">
			<FieldNumber :id="fieldId('text-size')" v-model="textFontSize" :label="t('overlayBuilder.editors.text.size')" :min="1" :step="1" />
			<FieldSelect :id="fieldId('text-weight')" v-model="textFontWeight" :label="t('overlayBuilder.editors.text.weight')" :options="fontWeightOptions" />
		</div>

		<FieldInput :id="fieldId('text-color')" v-model="textColor" :label="t('overlayBuilder.editors.text.color')" type="color" input-class="h-9 p-1" />
		<FieldSelect :id="fieldId('text-align')" v-model="textAlign" :label="t('overlayBuilder.editors.text.align')" :options="textAlignOptions.map((option) => ({ value: option.value, label: t(option.labelKey) }))" />
		<FieldInput :id="fieldId('text-family')" v-model="textFontFamily" :label="t('overlayBuilder.editors.text.font')" list="overlay-font-families" />
		<datalist id="overlay-font-families">
			<option value="sans-serif" />
			<option value="serif" />
			<option value="monospace" />
			<option value="Inter" />
			<option value="Arial" />
		</datalist>
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
