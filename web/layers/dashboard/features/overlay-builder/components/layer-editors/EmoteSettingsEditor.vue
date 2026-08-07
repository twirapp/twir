<script setup lang="ts">
import { toRef } from 'vue'
import EmotePicker from '~~/layers/dashboard/components/emote-picker/EmotePicker.vue'

import { Button } from '@/components/ui/button'

import type { Layer, LayerSettings } from '../../types'
import FieldInput from '../fields/FieldInput.vue'
import { useEmoteLayerEditor } from '../../composables/useEmoteLayerEditor'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()
const { t } = useI18n()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
}>()

function updateSettings(updates: Partial<LayerSettings>) {
	emit('update', {
		settings: {
			...props.layer.settings,
			...updates,
		},
	})
}

const { emoteUrl, selectEmote } = useEmoteLayerEditor(toRef(props, 'layer'), updateSettings)
const fieldId = (name: string) => `layer-${props.layer.id}-${name}`
</script>

<template>
	<div class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">{{ t('overlayBuilder.editors.emote.title') }}</h4>
		<EmotePicker @select="selectEmote">
			<template #trigger>
				<Button type="button" variant="outline" class="w-full">
					<Icon name="lucide:smile-plus" class="size-4" />
					{{ t('overlayBuilder.editors.emote.select') }}
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
				{{ t('overlayBuilder.editors.emote.manualUrl') }}
				<Icon name="lucide:chevron-down" class="ml-auto size-4 transition-transform group-open:rotate-180" />
			</summary>
			<FieldInput
				:id="fieldId('emote-url')"
				v-model="emoteUrl"
				:label="t('overlayBuilder.editors.emote.directUrl')"
				type="url"
				placeholder="https://.../emote.png"
				:description="t('overlayBuilder.editors.emote.urlDescription')"
				class="mt-3 flex flex-col gap-2"
			/>
		</details>
		<img v-if="emoteUrl" :src="emoteUrl" :alt="t('overlayBuilder.editors.emote.previewAlt')" class="max-h-32 max-w-full self-center object-contain" />
	</div>
</template>
