<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import Form from '~~/layers/dashboard/pages/dashboard/overlays/chat/components/Form.vue'

import { useChatOverlayPresets } from '../../composables/useChatOverlayPresets'

const emit = defineEmits<{
	// Fired when the active chat preset changes (initial selection included), so the
	// layer editor can keep the widget URL (?id=<presetId>) in sync.
	'select-preset': [id: string]
}>()

const { fetchingOverlays, presets, selectedPresetId, handlePresetChange, createPreset } = useChatOverlayPresets(
	(id) => emit('select-preset', id),
)
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
