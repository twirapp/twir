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
	{ value: 'left', label: 'Слева' },
	{ value: 'center', label: 'По центру' },
	{ value: 'right', label: 'Справа' },
] as const
</script>

<template>
	<div v-if="layer.type === 'TEXT'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">Текст</h4>

		<div class="flex flex-col gap-2">
			<Label :for="fieldId('text-content')">Содержимое</Label>
			<Textarea :id="fieldId('text-content')" v-model="textContent" rows="3" @keydown.stop />
		</div>

		<div class="grid grid-cols-2 gap-3">
			<FieldNumber :id="fieldId('text-size')" v-model="textFontSize" label="Размер" :min="1" :step="1" />
			<FieldSelect :id="fieldId('text-weight')" v-model="textFontWeight" label="Начертание" :options="fontWeightOptions" />
		</div>

		<FieldInput :id="fieldId('text-color')" v-model="textColor" label="Цвет" type="color" input-class="h-9 p-1" />
		<FieldSelect :id="fieldId('text-align')" v-model="textAlign" label="Выравнивание" :options="textAlignOptions" />
		<FieldInput :id="fieldId('text-family')" v-model="textFontFamily" label="Шрифт" list="overlay-font-families" />
		<datalist id="overlay-font-families">
			<option value="sans-serif" />
			<option value="serif" />
			<option value="monospace" />
			<option value="Inter" />
			<option value="Arial" />
		</datalist>
	</div>

	<div v-else-if="layer.type === 'VIDEO'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">Видео</h4>
		<FieldInput :id="fieldId('video-url')" v-model="videoUrl" label="URL видео" type="url" placeholder="https://example.com/video.mp4" />
		<FieldSwitch :id="fieldId('video-loop')" v-model="videoLoop" label="Повторять" />
		<FieldSwitch :id="fieldId('video-muted')" v-model="videoMuted" label="Без звука" />
	</div>

	<div v-else-if="layer.type === 'IFRAME'">
		<IframeSettingsEditor :layer="layer" @update="emit('update', $event)" />
	</div>

	<div v-else-if="layer.type === 'YOUTUBE'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">YouTube</h4>
		<FieldInput :id="fieldId('youtube-id')" v-model="youtubeVideoId" label="Ссылка или ID видео" placeholder="https://youtu.be/... или ID" @blur="youtubeVideoId = parseYoutubeVideoId(youtubeVideoId)" />
		<FieldSwitch :id="fieldId('youtube-autoplay')" v-model="youtubeAutoplay" label="Автовоспроизведение" />
		<FieldSwitch :id="fieldId('youtube-loop')" v-model="youtubeLoop" label="Повторять" />
		<FieldSwitch :id="fieldId('youtube-muted')" v-model="youtubeMuted" label="Без звука" />
	</div>

	<div v-else-if="layer.type === 'EMOTE'">
		<EmoteSettingsEditor :layer="layer" @update="emit('update', $event)" />
	</div>
</template>
