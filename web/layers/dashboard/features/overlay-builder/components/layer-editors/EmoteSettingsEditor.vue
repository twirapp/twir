<script setup lang="ts">
import EmotePicker from '~~/layers/dashboard/components/emote-picker/EmotePicker.vue'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

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

const emoteUrl = computed({
	get: () => props.layer.settings.emoteUrl,
	set: (value: string) => updateSettings({ emoteUrl: value, emoteName: '', emoteProvider: '' }),
})

interface SelectedEmote {
	readonly url: string
	readonly name: string
	readonly provider: '7TV'
}

function selectEmote(emote: SelectedEmote) {
	updateSettings({
		emoteUrl: emote.url,
		emoteName: emote.name,
		emoteProvider: emote.provider,
	})
}
</script>

<template>
	<div class="flex flex-col gap-4">
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
