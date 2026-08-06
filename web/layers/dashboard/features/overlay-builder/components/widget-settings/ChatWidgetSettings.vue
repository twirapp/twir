<script setup lang="ts">
import { until } from '@vueuse/core'
import type { AcceptableValue } from 'reka-ui'

import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useChatOverlayApi } from '~~/layers/dashboard/api/overlays/chat.ts'
import Form from '~~/layers/dashboard/pages/dashboard/overlays/chat/components/Form.vue'
import { useChatOverlayForm } from '~~/layers/dashboard/pages/dashboard/overlays/chat/components/form.ts'

const emit = defineEmits<{
	// Fired when the active chat preset changes (initial selection included), so the
	// layer editor can keep the widget URL (?id=<presetId>) in sync.
	'select-preset': [id: string]
}>()

const chatOverlaysManager = useChatOverlayApi()
const creator = chatOverlaysManager.useOverlayCreate()
const { data: chatOverlaysData, fetching: fetchingOverlays } = chatOverlaysManager.useOverlaysQuery()
const { setData, getDefaultSettings } = useChatOverlayForm()

const selectedPresetId = ref<string>()
const presets = computed(() => chatOverlaysData.value?.chatOverlays ?? [])

watch(
	() => chatOverlaysData.value?.chatOverlays,
	(overlays) => {
		if (!overlays?.length) {
			selectedPresetId.value = undefined
			return
		}

		if (!overlays.some((overlay) => overlay.id === selectedPresetId.value)) {
			selectedPresetId.value = overlays[0]?.id ?? undefined
		}
	},
	{ immediate: true }
)

watch(selectedPresetId, (id) => {
	const preset = presets.value.find((overlay) => overlay.id === id)
	if (preset) setData(preset)
	if (id) emit('select-preset', id)
})

function handlePresetChange(id: AcceptableValue) {
	if (typeof id !== 'string') return

	selectedPresetId.value = id
	const preset = presets.value.find((overlay) => overlay.id === id)
	if (preset) setData(preset)
}

async function createPreset() {
	const previousLength = presets.value.length

	await creator.executeMutation({ input: getDefaultSettings() })
	await until(() => chatOverlaysData.value?.chatOverlays).changed()

	const overlays = chatOverlaysData.value?.chatOverlays ?? []
	if (overlays.length > previousLength) {
		selectedPresetId.value = overlays[overlays.length - 1]?.id ?? undefined
	}
}
</script>

<template>
	<div class="min-w-0">
		<div v-if="fetchingOverlays" class="flex min-h-48 flex-col items-center justify-center gap-3">
			<div class="border-primary h-7 w-7 animate-spin rounded-full border-4 border-t-transparent" />
			<p class="text-sm text-muted-foreground">Загрузка пресетов...</p>
		</div>

		<div v-else-if="presets.length === 0" class="flex min-h-48 flex-col items-center justify-center gap-4 text-center">
			<div class="space-y-1">
				<p class="font-medium">Нет пресетов чата</p>
				<p class="text-sm text-muted-foreground">Создайте пресет, чтобы настроить виджет.</p>
			</div>
			<Button type="button" @click="createPreset">
				<Icon name="lucide:plus" class="mr-2 h-4 w-4" />
				Создать пресет
			</Button>
		</div>

		<div v-else class="flex min-w-0 flex-col gap-4">
			<div v-if="presets.length > 1" class="flex flex-col gap-2">
				<Label for="chat-widget-preset">Пресет</Label>
				<Select :model-value="selectedPresetId" @update:model-value="handlePresetChange">
					<SelectTrigger id="chat-widget-preset" class="w-full">
						<SelectValue placeholder="Выберите пресет" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem v-for="(preset, index) in presets" :key="preset.id" :value="preset.id">
							Пресет #{{ index + 1 }}
						</SelectItem>
					</SelectContent>
				</Select>
			</div>

			<Form embedded />
		</div>
	</div>
</template>
