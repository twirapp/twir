<script setup lang="ts">
import { computed } from 'vue'
import { Code2, Eye, Lock, RefreshCw } from 'lucide-vue-next'


import ImageLayerEditor from './layer-editors/ImageLayerEditor.vue'

import type { Layer } from '../types'

interface Props {
	layer: Layer | null
	multipleSelected: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
	openCodeEditor: []
}>()

// Local reactive values for inputs
const localName = computed({
	get: () => props.layer?.name ?? '',
	set: (value: string) => emit('update', { name: value }),
})

const localPosX = computed({
	get: () => props.layer?.posX ?? 0,
	set: (value: number) => emit('update', { posX: value }),
})

const localPosY = computed({
	get: () => props.layer?.posY ?? 0,
	set: (value: number) => emit('update', { posY: value }),
})

const localWidth = computed({
	get: () => props.layer?.width ?? 0,
	set: (value: number) => emit('update', { width: value }),
})

const localHeight = computed({
	get: () => props.layer?.height ?? 0,
	set: (value: number) => emit('update', { height: value }),
})

const localRotation = computed({
	get: () => props.layer?.rotation ?? 0,
	set: (value: number) => emit('update', { rotation: value }),
})

const localOpacity = computed({
	get: () => (props.layer?.opacity ?? 1) * 100,
	set: (value: number) => emit('update', { opacity: value / 100 }),
})

const localVisible = computed({
	get: () => props.layer?.visible ?? true,
	set: (value: boolean) => emit('update', { visible: value }),
})

const localLocked = computed({
	get: () => props.layer?.locked ?? false,
	set: (value: boolean) => emit('update', { locked: value }),
})

const localPeriodicallyRefetch = computed({
	get: () => props.layer?.periodicallyRefetchData ?? true,
	set: (value: boolean) => emit('update', { periodicallyRefetchData: value }),
})

const localPollInterval = computed({
	get: () => props.layer?.settings?.htmlOverlayDataPollSecondsInterval ?? 5,
	set: (value: number) => {
		if (!props.layer) return
		emit('update', {
			settings: {
				...props.layer.settings,
				htmlOverlayDataPollSecondsInterval: value,
			},
		})
	},
})
</script>

<template>
	<UiCard class="h-full flex flex-col border-0 p-0">
		<div class="border-b p-2">
			<UiCardTitle class="text-sm font-medium">
				Properties
			</UiCardTitle>
		</div>
		<UiCardContent class="flex-1 p-0 overflow-hidden">
			<UiScrollArea class="h-full">
				<div v-if="!layer && !multipleSelected" class="p-4 text-center text-muted-foreground">
					<p class="text-sm">No layer selected</p>
					<p class="text-xs mt-1">Select a layer to edit its properties</p>
				</div>

				<div v-else-if="multipleSelected" class="p-4 text-center text-muted-foreground">
					<p class="text-sm">Multiple layers selected</p>
					<p class="text-xs mt-1">Select a single layer to edit properties</p>
				</div>

				<div v-else-if="layer" class="p-4 space-y-6">
					<!-- Layer Info -->
					<div class="space-y-3">
						<div class="space-y-2">
							<UiLabel for="layer-name">Name</UiLabel>
							<UiInput id="layer-name" v-model="localName" placeholder="Layer name" @keydown.stop />
						</div>
					</div>

					<UiSeparator />

					<!-- Transform -->
					<UiTabs default-value="position" class="w-full">
						<UiTabsList class="grid w-full grid-cols-2">
							<UiTabsTrigger value="position">Position</UiTabsTrigger>
							<UiTabsTrigger value="appearance">Appearance</UiTabsTrigger>
						</UiTabsList>

						<UiTabsContent value="position" class="space-y-4 mt-4">
							<div class="grid grid-cols-2 gap-3">
								<div class="space-y-2">
									<UiLabel for="pos-x">X</UiLabel>
									<UiInput
										id="pos-x"
										v-model.number="localPosX"
										type="number"
										placeholder="0"
										@keydown.stop
									/>
								</div>
								<div class="space-y-2">
									<UiLabel for="pos-y">Y</UiLabel>
									<UiInput
										id="pos-y"
										v-model.number="localPosY"
										type="number"
										placeholder="0"
										@keydown.stop
									/>
								</div>
							</div>

							<div class="grid grid-cols-2 gap-3">
								<div class="space-y-2">
									<UiLabel for="width">Width</UiLabel>
									<UiInput
										id="width"
										v-model.number="localWidth"
										type="number"
										:min="1"
										placeholder="200"
										@keydown.stop
									/>
								</div>
								<div class="space-y-2">
									<UiLabel for="height">Height</UiLabel>
									<UiInput
										id="height"
										v-model.number="localHeight"
										type="number"
										:min="1"
										placeholder="200"
										@keydown.stop
									/>
								</div>
							</div>

							<div class="space-y-2">
								<div class="flex items-center justify-between">
									<UiLabel for="rotation">Rotation</UiLabel>
									<span class="text-xs text-muted-foreground">{{ localRotation }}°</span>
								</div>
								<UiSlider
									id="rotation"
									@update:model-value="(newValue) => {
										if (!newValue) return;
										localRotation = newValue[0]
									}"
									:model-value="[localRotation]"
									:min="0"
									:max="360"
									:step="1"
								/>
							</div>
						</UiTabsContent>

						<UiTabsContent value="appearance" class="space-y-4 mt-4">
							<div class="space-y-2">
								<div class="flex items-center justify-between">
									<UiLabel for="opacity">Opacity</UiLabel>
									<span class="text-xs text-muted-foreground">{{ Math.round(localOpacity) }}%</span>
								</div>
								<UiSlider
									id="opacity"
									@update:model-value="(newValue) => {
										if (!newValue) return;
										localOpacity = newValue[0]
									}"
									:model-value="[localOpacity]"
									:min="0"
									:max="100"
									:step="1"
								/>
							</div>

							<div class="flex items-center justify-between">
								<div class="space-y-0.5">
									<UiLabel for="visible">Visible</UiLabel>
									<p class="text-xs text-muted-foreground">Show layer on canvas</p>
								</div>
								<UiSwitch id="visible" v-model="localVisible">
									<Eye class="h-4 w-4" />
								</UiSwitch>
							</div>

							<div class="flex items-center justify-between">
								<div class="space-y-0.5">
									<UiLabel for="locked">Locked</UiLabel>
									<p class="text-xs text-muted-foreground">Prevent editing</p>
								</div>
								<UiSwitch id="locked" v-model="localLocked">
									<Lock class="h-4 w-4" />
								</UiSwitch>
							</div>
						</UiTabsContent>
					</UiTabs>

					<UiSeparator />

					<!-- HTML Layer Settings -->
					<div v-if="layer.type === 'HTML'" class="space-y-4">
						<div class="flex items-center justify-between">
							<h4 class="text-sm font-medium">HTML Settings</h4>
							<UiButton variant="outline" size="sm" @click="emit('openCodeEditor')">
								<Code2 class="h-4 w-4 mr-2" />
								Edit Code
							</UiButton>
						</div>

						<div class="flex items-center justify-between">
							<div class="space-y-0.5">
								<UiLabel for="auto-refresh">Auto Refresh</UiLabel>
								<p class="text-xs text-muted-foreground">Periodically update data</p>
							</div>
							<UiSwitch id="auto-refresh" v-model="localPeriodicallyRefetch">
								<RefreshCw class="h-4 w-4" />
							</UiSwitch>
						</div>

						<div v-if="localPeriodicallyRefetch" class="space-y-2">
							<UiLabel for="poll-interval">Update Interval (seconds)</UiLabel>
							<UiInput
								id="poll-interval"
								v-model.number="localPollInterval"
								type="number"
								:min="5"
								:max="300"
								@keydown.stop
							/>
						</div>
					</div>

					<!-- IMAGE Layer Settings -->
					<div v-if="layer.type === 'IMAGE'">
						<ImageLayerEditor
							:layer="layer"
							@update="emit('update', $event)"
						/>
					</div>
				</div>
			</UiScrollArea>
		</UiCardContent>
	</UiCard>
</template>
