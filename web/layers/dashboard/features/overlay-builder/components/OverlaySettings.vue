<script setup lang="ts">
import type { AcceptableValue } from 'reka-ui'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'

const overlayName = defineModel<string>('overlayName', {
	type: String,
	required: true,
})
const instaSave = defineModel<boolean>('instaSave', {
	type: Boolean,
	required: true,
})
const canvasWidth = defineModel<number>('canvasWidth', { type: Number, required: true })
const canvasHeight = defineModel<number>('canvasHeight', { type: Number, required: true })
const { t } = useI18n()

const MAX_CANVAS_DIMENSION = 7680
const canvasPresets = [
	{ value: 'full-hd', width: 1920, height: 1080, labelKey: 'overlayBuilder.canvasSize.fullHd' },
	{ value: 'hd', width: 1280, height: 720, labelKey: 'overlayBuilder.canvasSize.hd' },
	{ value: 'qhd', width: 2560, height: 1440, labelKey: 'overlayBuilder.canvasSize.qhd' },
	{ value: '4k', width: 3840, height: 2160, labelKey: 'overlayBuilder.canvasSize.fourK' },
] as const
const selectedCanvasPreset = computed(() => canvasPresets.find((preset) => preset.width === canvasWidth.value && preset.height === canvasHeight.value)?.value ?? 'custom')

function updateCanvasPreset(value: AcceptableValue) {
	if (typeof value !== 'string') return
	const preset = canvasPresets.find((item) => item.value === value)
	if (!preset) return
	canvasWidth.value = preset.width
	canvasHeight.value = preset.height
}

function updateCanvasDimension(dimension: 'width' | 'height', event: Event) {
	if (!(event.target instanceof HTMLInputElement) || !Number.isFinite(event.target.valueAsNumber)) return
	const value = Math.min(MAX_CANVAS_DIMENSION, Math.max(1, Math.round(event.target.valueAsNumber)))
	if (dimension === 'width') canvasWidth.value = value
	else canvasHeight.value = value
}
</script>

<template>
	<Card class="border-0 shadow-none p-0">
		<div class="border-b p-2">
			<CardTitle class="text-sm font-medium">{{ t('overlayBuilder.overlaySettings.title') }}</CardTitle>
		</div>
		<CardContent class="px-3 pb-3 space-y-3">
			<!-- Overlay Name -->
			<div class="space-y-1.5">
				<Label for="overlay-name" class="text-xs">
					{{ t('overlayBuilder.overlaySettings.name') }} <span class="text-destructive">*</span>
				</Label>
				<Input
					id="overlay-name"
					v-model="overlayName"
					:placeholder="t('overlayBuilder.overlaySettings.namePlaceholder')"
					maxlength="30"
					class="h-8 text-sm"
					@keydown.stop
				/>
				<p class="text-xs text-muted-foreground">
					{{ t('overlayBuilder.overlaySettings.characterCount', { count: overlayName.length }) }}
				</p>
			</div>

			<div class="space-y-1.5">
				<Label for="canvas-size" class="text-xs">{{ t('overlayBuilder.canvasSize.label') }}</Label>
				<Select :model-value="selectedCanvasPreset" @update:model-value="updateCanvasPreset">
					<SelectTrigger id="canvas-size" class="h-8 w-full text-sm"><SelectValue /></SelectTrigger>
					<SelectContent>
						<SelectItem v-for="preset in canvasPresets" :key="preset.value" :value="preset.value">{{ t(preset.labelKey) }}</SelectItem>
						<SelectItem value="custom">{{ t('overlayBuilder.canvasSize.custom') }}</SelectItem>
					</SelectContent>
				</Select>
				<div v-if="selectedCanvasPreset === 'custom'" class="grid grid-cols-2 gap-2">
					<div class="space-y-1.5"><Label for="canvas-width" class="text-xs">{{ t('overlayBuilder.canvasSize.width') }}</Label><Input id="canvas-width" :model-value="canvasWidth" type="number" min="1" :max="MAX_CANVAS_DIMENSION" class="h-8 text-sm" @change="updateCanvasDimension('width', $event)" @keydown.stop /></div>
					<div class="space-y-1.5"><Label for="canvas-height" class="text-xs">{{ t('overlayBuilder.canvasSize.height') }}</Label><Input id="canvas-height" :model-value="canvasHeight" type="number" min="1" :max="MAX_CANVAS_DIMENSION" class="h-8 text-sm" @change="updateCanvasDimension('height', $event)" @keydown.stop /></div>
				</div>
			</div>

			<!-- Insta Save -->
			<div class="space-y-1.5">
				<div class="flex items-center justify-between">
					<Label for="insta-save" class="text-xs">
						{{ t('overlayBuilder.overlaySettings.instantSave') }}
					</Label>
					<Switch
						id="insta-save"
						:model-value="instaSave"
						@update:model-value="(val: boolean) => (instaSave = val)"
					/>
				</div>
				<p class="text-xs text-muted-foreground">
					{{ t('overlayBuilder.overlaySettings.instantSaveDescription') }}
				</p>
			</div>
		</CardContent>
	</Card>
</template>
