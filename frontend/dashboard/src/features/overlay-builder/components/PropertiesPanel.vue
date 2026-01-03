<script setup lang="ts">
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Code2, Eye, Lock, RefreshCw } from 'lucide-vue-next'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Slider } from '@/components/ui/slider'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
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

const { t } = useI18n()

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
	<Card class="h-full flex flex-col border-0 p-0">
		<div class="border-b p-2">
			<CardTitle class="text-sm font-medium">
				Properties
			</CardTitle>
		</div>
		<CardContent class="flex-1 p-0 overflow-hidden">
			<ScrollArea class="h-full">
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
							<Label for="layer-name">Name</Label>
							<Input id="layer-name" v-model="localName" placeholder="Layer name" />
						</div>
					</div>

					<Separator />

					<!-- Transform -->
					<Tabs default-value="position" class="w-full">
						<TabsList class="grid w-full grid-cols-2">
							<TabsTrigger value="position">Position</TabsTrigger>
							<TabsTrigger value="appearance">Appearance</TabsTrigger>
						</TabsList>

						<TabsContent value="position" class="space-y-4 mt-4">
							<div class="grid grid-cols-2 gap-3">
								<div class="space-y-2">
									<Label for="pos-x">X</Label>
									<Input
										id="pos-x"
										v-model.number="localPosX"
										type="number"
										placeholder="0"
									/>
								</div>
								<div class="space-y-2">
									<Label for="pos-y">Y</Label>
									<Input
										id="pos-y"
										v-model.number="localPosY"
										type="number"
										placeholder="0"
									/>
								</div>
							</div>

							<div class="grid grid-cols-2 gap-3">
								<div class="space-y-2">
									<Label for="width">Width</Label>
									<Input
										id="width"
										v-model.number="localWidth"
										type="number"
										:min="1"
										placeholder="200"
									/>
								</div>
								<div class="space-y-2">
									<Label for="height">Height</Label>
									<Input
										id="height"
										v-model.number="localHeight"
										type="number"
										:min="1"
										placeholder="200"
									/>
								</div>
							</div>

							<div class="space-y-2">
								<div class="flex items-center justify-between">
									<Label for="rotation">Rotation</Label>
									<span class="text-xs text-muted-foreground">{{ localRotation }}°</span>
								</div>
								<Slider
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
						</TabsContent>

						<TabsContent value="appearance" class="space-y-4 mt-4">
							<div class="space-y-2">
								<div class="flex items-center justify-between">
									<Label for="opacity">Opacity</Label>
									<span class="text-xs text-muted-foreground">{{ Math.round(localOpacity) }}%</span>
								</div>
								<Slider
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
									<Label for="visible">Visible</Label>
									<p class="text-xs text-muted-foreground">Show layer on canvas</p>
								</div>
								<Switch id="visible" v-model:checked="localVisible">
									<Eye class="h-4 w-4" />
								</Switch>
							</div>

							<div class="flex items-center justify-between">
								<div class="space-y-0.5">
									<Label for="locked">Locked</Label>
									<p class="text-xs text-muted-foreground">Prevent editing</p>
								</div>
								<Switch id="locked" v-model:checked="localLocked">
									<Lock class="h-4 w-4" />
								</Switch>
							</div>
						</TabsContent>
					</Tabs>

					<Separator />

					<!-- HTML Layer Settings -->
					<div v-if="layer.type === 'HTML'" class="space-y-4">
						<div class="flex items-center justify-between">
							<h4 class="text-sm font-medium">HTML Settings</h4>
							<Button variant="outline" size="sm" @click="emit('openCodeEditor')">
								<Code2 class="h-4 w-4 mr-2" />
								Edit Code
							</Button>
						</div>

						<div class="flex items-center justify-between">
							<div class="space-y-0.5">
								<Label for="auto-refresh">Auto Refresh</Label>
								<p class="text-xs text-muted-foreground">Periodically update data</p>
							</div>
							<Switch id="auto-refresh" v-model:checked="localPeriodicallyRefetch">
								<RefreshCw class="h-4 w-4" />
							</Switch>
						</div>

						<div v-if="localPeriodicallyRefetch" class="space-y-2">
							<Label for="poll-interval">Update Interval (seconds)</Label>
							<Input
								id="poll-interval"
								v-model.number="localPollInterval"
								type="number"
								:min="5"
								:max="300"
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
			</ScrollArea>
		</CardContent>
	</Card>
</template>
