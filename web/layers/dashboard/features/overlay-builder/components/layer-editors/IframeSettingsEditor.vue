<script setup lang="ts">
import { toRef } from 'vue'

import DialogOrSheet from '~~/layers/dashboard/components/dialog-or-sheet.vue'

import { Button } from '@/components/ui/button'
import { Dialog, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

import type { Layer } from '../../types'
import { useWidgetLayer } from '../../composables/useWidgetLayer'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
}>()

const {
	fieldId,
	overlayApiKey,
	selectedWidget,
	selectedWidgetSettings,
	widgetSettingsOpen,
	iframeSource,
	iframeUrl,
	iframeScale,
	handleIframeSourceChange,
	handleWidgetChange,
	handleWidgetPresetSelect,
	overlayWidgetRegistry,
} = useWidgetLayer(toRef(props, 'layer'), (updates) => emit('update', updates))
</script>

<template>
	<div class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">Виджет</h4>
		<div class="flex flex-col gap-2">
			<Label :for="fieldId('iframe-source')">Источник</Label>
			<Select :model-value="iframeSource" @update:model-value="handleIframeSourceChange">
				<SelectTrigger :id="fieldId('iframe-source')" class="w-full"><SelectValue /></SelectTrigger>
				<SelectContent>
					<SelectItem value="custom">Свой URL</SelectItem>
					<SelectItem value="twir">Виджет Twir</SelectItem>
				</SelectContent>
			</Select>
		</div>

		<template v-if="iframeSource === 'twir'">
			<div class="flex flex-col gap-2">
				<Label :for="fieldId('widget')">Виджет Twir</Label>
				<Select :model-value="layer.settings.widgetKey" @update:model-value="handleWidgetChange">
					<SelectTrigger :id="fieldId('widget')" class="w-full"><SelectValue placeholder="Выберите виджет" /></SelectTrigger>
					<SelectContent>
						<SelectItem v-for="widget in overlayWidgetRegistry" :key="widget.key" :value="widget.key">
							<div class="flex items-center gap-2">
								<Icon :name="widget.icon" class="h-4 w-4" />
								<span>{{ widget.name }}</span>
							</div>
						</SelectItem>
					</SelectContent>
				</Select>
				<p class="text-xs text-muted-foreground">{{ selectedWidget?.description || 'Выберите встроенный виджет Twir.' }}</p>
			</div>

			<Button
				v-if="selectedWidgetSettings"
				type="button"
				variant="outline"
				class="w-full"
				@click="widgetSettingsOpen = true"
			>
				<Icon name="lucide:settings-2" class="mr-2 h-4 w-4" />
				Настроить
			</Button>

			<p v-if="layer.settings.widgetKey && !overlayApiKey" class="text-xs text-destructive">
				Не найден API-ключ канала. Ссылка на виджет не будет работать.
			</p>
		</template>

		<div v-else class="flex flex-col gap-2">
			<Label :for="fieldId('iframe-url')">URL</Label>
			<Input :id="fieldId('iframe-url')" v-model="iframeUrl" type="url" placeholder="https://example.com/widget" @keydown.stop />
		</div>

		<div class="flex flex-col gap-2">
			<Label :for="fieldId('iframe-scale')">Масштаб</Label>
			<Input :id="fieldId('iframe-scale')" v-model.number="iframeScale" type="number" min="0.1" max="4" step="0.1" @keydown.stop />
		</div>
	</div>

	<Dialog v-model:open="widgetSettingsOpen">
		<DialogOrSheet v-if="selectedWidgetSettings" class="gap-0 p-0 sm:max-w-5xl">
			<DialogHeader class="border-b px-6 py-4">
				<DialogTitle>{{ selectedWidget?.name }}</DialogTitle>
				<DialogDescription class="sr-only">Настройки виджета {{ selectedWidget?.name }}</DialogDescription>
			</DialogHeader>
			<div class="min-w-0 p-6">
				<component :is="selectedWidgetSettings" @select-preset="handleWidgetPresetSelect" />
			</div>
		</DialogOrSheet>
	</Dialog>
</template>
