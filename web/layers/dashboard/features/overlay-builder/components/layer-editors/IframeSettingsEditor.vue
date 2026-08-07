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
const { t } = useI18n()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
}>()

const {
	fieldId,
	overlayApiKey,
	selectedWidget,
	selectedWidgetSettings,
	selectedWidgetParams,
	widgetSettingsOpen,
	iframeSource,
	iframeUrl,
	iframeScale,
	handleIframeSourceChange,
	handleWidgetChange,
	handleWidgetPresetSelect,
	handleWidgetParamsUpdate,
	overlayWidgetRegistry,
} = useWidgetLayer(toRef(props, 'layer'), (updates) => emit('update', updates))
</script>

<template>
	<div class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">{{ t('overlayBuilder.editors.iframe.title') }}</h4>
		<div class="flex flex-col gap-2">
			<Label :for="fieldId('iframe-source')">{{ t('overlayBuilder.editors.iframe.source') }}</Label>
			<Select :model-value="iframeSource" @update:model-value="handleIframeSourceChange">
				<SelectTrigger :id="fieldId('iframe-source')" class="w-full"><SelectValue /></SelectTrigger>
				<SelectContent>
					<SelectItem value="custom">{{ t('overlayBuilder.editors.iframe.customUrl') }}</SelectItem>
					<SelectItem value="twir">{{ t('overlayBuilder.editors.iframe.twirWidget') }}</SelectItem>
				</SelectContent>
			</Select>
		</div>

		<template v-if="iframeSource === 'twir'">
			<div class="flex flex-col gap-2">
				<Label :for="fieldId('widget')">{{ t('overlayBuilder.editors.iframe.twirWidget') }}</Label>
				<Select :model-value="layer.settings.widgetKey" @update:model-value="handleWidgetChange">
					<SelectTrigger :id="fieldId('widget')" class="w-full"><SelectValue :placeholder="t('overlayBuilder.editors.iframe.widgetPlaceholder')" /></SelectTrigger>
					<SelectContent>
						<SelectItem v-for="widget in overlayWidgetRegistry" :key="widget.key" :value="widget.key">
							<div class="flex items-center gap-2">
								<Icon :name="widget.icon" class="h-4 w-4" />
								<span>{{ t(widget.nameKey) }}</span>
							</div>
						</SelectItem>
					</SelectContent>
				</Select>
				<p class="text-xs text-muted-foreground">{{ selectedWidget ? t(selectedWidget.descriptionKey) : t('overlayBuilder.editors.iframe.widgetDescription') }}</p>
			</div>

			<Button
				v-if="selectedWidgetSettings"
				type="button"
				variant="outline"
				class="w-full"
				@click="widgetSettingsOpen = true"
			>
				<Icon name="lucide:settings-2" class="mr-2 h-4 w-4" />
				{{ t('overlayBuilder.editors.iframe.configure') }}
			</Button>

			<p v-if="layer.settings.widgetKey && !overlayApiKey" class="text-xs text-destructive">
				{{ t('overlayBuilder.editors.iframe.apiKeyMissing') }}
			</p>
		</template>

		<div v-else class="flex flex-col gap-2">
			<Label :for="fieldId('iframe-url')">{{ t('overlayBuilder.editors.iframe.url') }}</Label>
			<Input :id="fieldId('iframe-url')" v-model="iframeUrl" type="url" placeholder="https://example.com/widget" @keydown.stop />
		</div>

		<div class="flex flex-col gap-2">
			<Label :for="fieldId('iframe-scale')">{{ t('overlayBuilder.editors.iframe.scale') }}</Label>
			<Input :id="fieldId('iframe-scale')" v-model.number="iframeScale" type="number" min="0.1" max="4" step="0.1" @keydown.stop />
		</div>
	</div>

	<Dialog v-model:open="widgetSettingsOpen">
		<DialogOrSheet v-if="selectedWidgetSettings" class="gap-0 p-0 sm:max-w-5xl">
			<DialogHeader class="border-b px-6 py-4">
				<DialogTitle>{{ selectedWidget ? t(selectedWidget.nameKey) : '' }}</DialogTitle>
				<DialogDescription class="sr-only">{{ t('overlayBuilder.editors.iframe.dialogDescription', { name: selectedWidget ? t(selectedWidget.nameKey) : '' }) }}</DialogDescription>
			</DialogHeader>
			<div class="min-w-0 p-6">
				<component
					:is="selectedWidgetSettings"
					v-bind="selectedWidgetParams"
					@select-preset="handleWidgetPresetSelect"
					@update-params="handleWidgetParamsUpdate"
				/>
			</div>
		</DialogOrSheet>
	</Dialog>
</template>
