<script setup lang="ts">
import EmotePicker from '~~/layers/dashboard/components/emote-picker/EmotePicker.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'

import type { Layer, LayerSettings } from '../../types'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
}>()

const fieldId = (name: string) => `layer-${props.layer.id}-${name}`

function updateSettings(updates: Partial<LayerSettings>) {
	emit('update', {
		settings: {
			...props.layer.settings,
			...updates,
		},
	})
}

const textContent = computed({
	get: () => props.layer.settings.textContent,
	set: (value: string | number) => updateSettings({ textContent: String(value) }),
})

const textFontSize = computed({
	get: () => props.layer.settings.textFontSize,
	set: (value: string | number) => {
		const parsed = Number(value)
		if (Number.isFinite(parsed)) updateSettings({ textFontSize: Math.max(1, Math.round(parsed)) })
	},
})

const textFontWeight = computed({
	get: () => String(props.layer.settings.textFontWeight),
	set: (value: string) => updateSettings({ textFontWeight: Number(value) }),
})

const textColor = computed({
	get: () => props.layer.settings.textColor,
	set: (value: string) => updateSettings({ textColor: value }),
})

const textAlign = computed({
	get: () => props.layer.settings.textAlign,
	set: (value: string) => updateSettings({ textAlign: value }),
})

const textFontFamily = computed({
	get: () => props.layer.settings.textFontFamily,
	set: (value: string) => updateSettings({ textFontFamily: value }),
})

const videoUrl = computed({
	get: () => props.layer.settings.videoUrl,
	set: (value: string) => updateSettings({ videoUrl: value }),
})

const videoLoop = computed({
	get: () => props.layer.settings.videoLoop,
	set: (value: boolean) => updateSettings({ videoLoop: value }),
})

const videoMuted = computed({
	get: () => props.layer.settings.videoMuted,
	set: (value: boolean) => updateSettings({ videoMuted: value }),
})

const iframeUrl = computed({
	get: () => props.layer.settings.iframeUrl,
	set: (value: string) => updateSettings({ iframeUrl: value }),
})

const iframeScale = computed({
	get: () => props.layer.settings.iframeScale,
	set: (value: string | number) => {
		const parsed = Number(value)
		if (Number.isFinite(parsed)) updateSettings({ iframeScale: Math.min(4, Math.max(0.1, parsed)) })
	},
})

const youtubeVideoId = computed({
	get: () => props.layer.settings.youtubeVideoId,
	set: (value: string) => updateSettings({ youtubeVideoId: value }),
})

const youtubeAutoplay = computed({
	get: () => props.layer.settings.youtubeAutoplay,
	set: (value: boolean) => updateSettings({ youtubeAutoplay: value }),
})

const youtubeLoop = computed({
	get: () => props.layer.settings.youtubeLoop,
	set: (value: boolean) => updateSettings({ youtubeLoop: value }),
})

const youtubeMuted = computed({
	get: () => props.layer.settings.youtubeMuted,
	set: (value: boolean) => updateSettings({ youtubeMuted: value }),
})

const emoteUrl = computed({
	get: () => props.layer.settings.emoteUrl,
	set: (value: string) => updateSettings({ emoteUrl: value, emoteName: '', emoteProvider: '' }),
})

interface SelectedEmote {
	url: string
	name: string
	provider: '7TV'
}

function selectEmote(emote: SelectedEmote) {
	updateSettings({
		emoteUrl: emote.url,
		emoteName: emote.name,
		emoteProvider: emote.provider,
	})
}

function parseYoutubeVideoId(value: string): string {
	const trimmed = value.trim()
	if (!trimmed) return ''

	const match = trimmed.match(
		/(?:youtube\.com\/watch\?[^#\s]*?v=|youtu\.be\/|youtube\.com\/(?:embed|shorts|live)\/)([^?&#/\s]+)/i
	)

	return match?.[1] ?? trimmed
}
</script>

<template>
	<div v-if="layer.type === 'TEXT'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">Текст</h4>

		<div class="flex flex-col gap-2">
			<Label :for="fieldId('text-content')">Содержимое</Label>
			<Textarea :id="fieldId('text-content')" v-model="textContent" rows="3" @keydown.stop />
		</div>

		<div class="grid grid-cols-2 gap-3">
			<div class="flex flex-col gap-2">
				<Label :for="fieldId('text-size')">Размер</Label>
				<Input :id="fieldId('text-size')" v-model.number="textFontSize" type="number" :min="1" :step="1" @keydown.stop />
			</div>
			<div class="flex flex-col gap-2">
				<Label :for="fieldId('text-weight')">Начертание</Label>
				<Select v-model="textFontWeight">
					<SelectTrigger :id="fieldId('text-weight')" class="w-full"><SelectValue /></SelectTrigger>
					<SelectContent>
						<SelectItem value="400">400</SelectItem>
						<SelectItem value="500">500</SelectItem>
						<SelectItem value="600">600</SelectItem>
						<SelectItem value="700">700</SelectItem>
						<SelectItem value="800">800</SelectItem>
					</SelectContent>
				</Select>
			</div>
		</div>

		<div class="flex flex-col gap-2">
			<Label :for="fieldId('text-color')">Цвет</Label>
			<Input :id="fieldId('text-color')" v-model="textColor" type="color" class="h-9 p-1" />
		</div>

		<div class="flex flex-col gap-2">
			<Label :for="fieldId('text-align')">Выравнивание</Label>
			<Select v-model="textAlign">
				<SelectTrigger :id="fieldId('text-align')" class="w-full"><SelectValue /></SelectTrigger>
				<SelectContent>
					<SelectItem value="left">Слева</SelectItem>
					<SelectItem value="center">По центру</SelectItem>
					<SelectItem value="right">Справа</SelectItem>
				</SelectContent>
			</Select>
		</div>

		<div class="flex flex-col gap-2">
			<Label :for="fieldId('text-family')">Шрифт</Label>
			<Input :id="fieldId('text-family')" v-model="textFontFamily" list="overlay-font-families" @keydown.stop />
			<datalist id="overlay-font-families">
				<option value="sans-serif" />
				<option value="serif" />
				<option value="monospace" />
				<option value="Inter" />
				<option value="Arial" />
			</datalist>
		</div>
	</div>

	<div v-else-if="layer.type === 'VIDEO'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">Видео</h4>
		<div class="flex flex-col gap-2">
			<Label :for="fieldId('video-url')">URL видео</Label>
			<Input :id="fieldId('video-url')" v-model="videoUrl" type="url" placeholder="https://example.com/video.mp4" @keydown.stop />
		</div>
		<div class="flex items-center justify-between">
			<Label :for="fieldId('video-loop')">Повторять</Label>
			<Switch :id="fieldId('video-loop')" v-model="videoLoop" />
		</div>
		<div class="flex items-center justify-between">
			<Label :for="fieldId('video-muted')">Без звука</Label>
			<Switch :id="fieldId('video-muted')" v-model="videoMuted" />
		</div>
	</div>

	<div v-else-if="layer.type === 'IFRAME'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">Виджет</h4>
		<div class="flex flex-col gap-2">
			<Label :for="fieldId('iframe-url')">URL</Label>
			<Input :id="fieldId('iframe-url')" v-model="iframeUrl" type="url" placeholder="https://example.com/widget" @keydown.stop />
		</div>
		<div class="flex flex-col gap-2">
			<Label :for="fieldId('iframe-scale')">Масштаб</Label>
			<Input :id="fieldId('iframe-scale')" v-model.number="iframeScale" type="number" min="0.1" max="4" step="0.1" @keydown.stop />
		</div>
	</div>

	<div v-else-if="layer.type === 'YOUTUBE'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">YouTube</h4>
		<div class="flex flex-col gap-2">
			<Label :for="fieldId('youtube-id')">Ссылка или ID видео</Label>
			<Input :id="fieldId('youtube-id')" v-model="youtubeVideoId" placeholder="https://youtu.be/... или ID" @blur="youtubeVideoId = parseYoutubeVideoId(youtubeVideoId)" @keydown.stop />
		</div>
		<div class="flex items-center justify-between">
			<Label :for="fieldId('youtube-autoplay')">Автовоспроизведение</Label>
			<Switch :id="fieldId('youtube-autoplay')" v-model="youtubeAutoplay" />
		</div>
		<div class="flex items-center justify-between">
			<Label :for="fieldId('youtube-loop')">Повторять</Label>
			<Switch :id="fieldId('youtube-loop')" v-model="youtubeLoop" />
		</div>
		<div class="flex items-center justify-between">
			<Label :for="fieldId('youtube-muted')">Без звука</Label>
			<Switch :id="fieldId('youtube-muted')" v-model="youtubeMuted" />
		</div>
	</div>

	<div v-else-if="layer.type === 'EMOTE'" class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">Эмоция</h4>
		<EmotePicker @select="selectEmote">
			<template #trigger>
				<Button type="button" variant="outline" class="w-full">
					<Icon name="lucide:smile-plus" class="size-4" />
					Выбрать эмоцию
				</Button>
			</template>
		</EmotePicker>
		<div v-if="layer.settings.emoteName" class="flex items-center gap-2 rounded-md border bg-muted/40 px-3 py-2 text-sm">
			<Icon name="lucide:badge-check" class="size-4 text-muted-foreground" />
			<span class="truncate">{{ layer.settings.emoteName }}</span>
			<span class="ml-auto shrink-0 text-xs text-muted-foreground">{{ layer.settings.emoteProvider }}</span>
		</div>
		<details class="group rounded-md border px-3 py-2">
			<summary class="flex cursor-pointer list-none items-center gap-2 text-sm text-muted-foreground">
				<Icon name="lucide:link" class="size-4" />
				Ввести URL вручную
				<Icon name="lucide:chevron-down" class="ml-auto size-4 transition-transform group-open:rotate-180" />
			</summary>
			<div class="mt-3 flex flex-col gap-2">
				<Label :for="fieldId('emote-url')">Прямая ссылка</Label>
				<Input :id="fieldId('emote-url')" v-model="emoteUrl" type="url" placeholder="https://.../emote.png" @keydown.stop />
				<p class="text-xs text-muted-foreground">Если нужной эмоции нет в 7TV, вставьте ссылку на изображение вручную.</p>
			</div>
		</details>
		<img v-if="emoteUrl" :src="emoteUrl" alt="Предпросмотр эмоции" class="max-h-32 max-w-full self-center object-contain" />
	</div>
</template>
