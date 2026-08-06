<script setup lang="ts">
import type { AcceptableValue } from 'reka-ui'

import DialogOrSheet from '~~/layers/dashboard/components/dialog-or-sheet.vue'
import { useProfile } from '~~/layers/dashboard/api/auth.js'

import { Button } from '@/components/ui/button'
import { Dialog, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

import type { Layer, LayerSettings } from '../../types'
import { overlayWidgetRegistry } from '../../widgets-registry'

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

const { data: profile } = useProfile()
const requestUrl = useRequestURL()

const selectedDashboard = computed(() => {
	return profile.value?.availableDashboards.find(
		(dashboard) => dashboard.id === profile.value?.selectedDashboardId
	)
})

const overlayApiKey = computed(() => {
	return selectedDashboard.value?.channelApiKey || profile.value?.channelApiKey || ''
})

const selectedWidget = computed(() => {
	return overlayWidgetRegistry.find((widget) => widget.key === props.layer.settings.widgetKey)
})

const selectedWidgetSettings = computed(() => selectedWidget.value?.settingsComponent)
const widgetSettingsOpen = ref(false)
type IframeSource = 'custom' | 'twir'
const iframeSource = ref<IframeSource>(props.layer.settings.widgetKey ? 'twir' : 'custom')

watch(
	() => props.layer.settings.widgetKey,
	(widgetKey) => {
		iframeSource.value = widgetKey ? 'twir' : 'custom'
		if (!widgetKey) widgetSettingsOpen.value = false
	}
)

function handleIframeSourceChange(value: AcceptableValue) {
	if (typeof value !== 'string') return

	if (value === 'custom') {
		iframeSource.value = 'custom'
		updateSettings({ widgetKey: '' })
		return
	}

	if (value === 'twir') iframeSource.value = 'twir'
}

function handleWidgetChange(key: AcceptableValue) {
	if (typeof key !== 'string') return

	const widget = overlayWidgetRegistry.find((entry) => entry.key === key)
	if (!widget) return

	iframeSource.value = 'twir'
	updateSettings({
		widgetKey: widget.key,
		iframeUrl: widget.buildUrl({
			origin: requestUrl.origin,
			apiKey: overlayApiKey.value,
		}),
	})
}

const iframeUrl = computed({
	get: () => props.layer.settings.iframeUrl,
	set: (value: string) => {
		iframeSource.value = 'custom'
		updateSettings({ iframeUrl: value, widgetKey: '' })
	},
})

const iframeScale = computed({
	get: () => props.layer.settings.iframeScale,
	set: (value: string | number) => {
		const parsed = Number(value)
		if (Number.isFinite(parsed)) updateSettings({ iframeScale: Math.min(4, Math.max(0.1, parsed)) })
	},
})
</script>

<template>
	<div class="flex flex-col gap-4">
		<h4 class="text-sm font-medium">Виджет</h4>
		<div class="flex flex-col gap-2">
			<Label :for="fieldId('iframe-source')">Источник</Label>
			<Select :model-value="iframeSource" @update:model-value="handleIframeSourceChange">
				<SelectTrigger :id="fieldId('iframe-source')" class="w-full">
					<SelectValue />
				</SelectTrigger>
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
					<SelectTrigger :id="fieldId('widget')" class="w-full">
						<SelectValue placeholder="Выберите виджет" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem v-for="widget in overlayWidgetRegistry" :key="widget.key" :value="widget.key">
							<div class="flex items-center gap-2">
								<Icon :name="widget.icon" class="h-4 w-4" />
								<span>{{ widget.name }}</span>
							</div>
						</SelectItem>
					</SelectContent>
				</Select>
				<p class="text-xs text-muted-foreground">
					{{ selectedWidget?.description || 'Выберите встроенный виджет Twir.' }}
				</p>
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
				<component :is="selectedWidgetSettings" />
			</div>
		</DialogOrSheet>
	</Dialog>
</template>
