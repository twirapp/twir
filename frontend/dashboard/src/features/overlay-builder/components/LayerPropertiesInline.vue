<script setup lang="ts">
import { computed } from 'vue'
import { Code2, Eye, Lock, RefreshCw } from 'lucide-vue-next'

import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Slider } from '@/components/ui/slider'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import ImageLayerEditor from './layer-editors/ImageLayerEditor.vue'

import type { Layer } from '../types'

interface Props {
	layer: Layer
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
	<div class="space-y-4 bg-muted/30 rounded-md p-3">
		<!-- Layer Name -->
		<div class="space-y-2">
			<Label for="layer-name" class="text-xs">Name</Label>
			<Input
				id="layer-name"
				v-model="localName"
				placeholder="Layer name"
				class="h-8 text-xs"
				@keydown.stop
			/>
		</div>

		<Separator />

		<!-- Transform Tabs -->
		<Tabs default-value="position" class="w-full">
			<TabsList class="grid w-full grid-cols-2 h-8">
				<TabsTrigger value="position" class="text-xs">Position</TabsTrigger>
				<TabsTrigger value="appearance" class="text-xs">Appearance</TabsTrigger>
			</TabsList>

			<TabsContent value="position" class="space-y-3 mt-3">
				<div class="grid grid-cols-2 gap-2">
					<div class="space-y-1.5">
						<Label for="pos-x" class="text-xs">X</Label>
						<Input
							id="pos-x"
							v-model.number="localPosX"
							type="number"
							placeholder="0"
							class="h-8 text-xs"
							@keydown.stop
						/>
					</div>
					<div class="space-y-1.5">
						<Label for="pos-y" class="text-xs">Y</Label>
						<Input
							id="pos-y"
							v-model.number="localPosY"
							type="number"
							placeholder="0"
							class="h-8 text-xs"
							@keydown.stop
						/>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-2">
					<div class="space-y-1.5">
						<Label for="width" class="text-xs">Width</Label>
						<Input
							id="width"
							v-model.number="localWidth"
							type="number"
							:min="1"
							placeholder="200"
							class="h-8 text-xs"
							@keydown.stop
						/>
					</div>
					<div class="space-y-1.5">
						<Label for="height" class="text-xs">Height</Label>
						<Input
							id="height"
							v-model.number="localHeight"
							type="number"
							:min="1"
							placeholder="200"
							class="h-8 text-xs"
							@keydown.stop
						/>
					</div>
				</div>

				<div class="space-y-1.5">
					<div class="flex items-center justify-between">
						<Label for="rotation" class="text-xs">Rotation</Label>
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

			<TabsContent value="appearance" class="space-y-3 mt-3">
				<div class="space-y-1.5">
					<div class="flex items-center justify-between">
						<Label for="opacity" class="text-xs">Opacity</Label>
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
						<Label for="visible" class="text-xs">Visible</Label>
						<p class="text-xs text-muted-foreground">Show layer</p>
					</div>
					<Switch id="visible" v-model="localVisible">
						<Eye class="h-3 w-3" />
					</Switch>
				</div>

				<div class="flex items-center justify-between">
					<div class="space-y-0.5">
						<Label for="locked" class="text-xs">Locked</Label>
						<p class="text-xs text-muted-foreground">Prevent editing</p>
					</div>
					<Switch id="locked" v-model="localLocked">
						<Lock class="h-3 w-3" />
					</Switch>
				</div>
			</TabsContent>
		</Tabs>

		<Separator />

		<!-- HTML Layer Settings -->
		<div v-if="layer.type === 'HTML'" class="space-y-3">
			<div class="flex items-center justify-between">
				<h4 class="text-xs font-medium">HTML Settings</h4>
				<Button variant="outline" size="sm" class="h-7 text-xs" @click="emit('openCodeEditor')">
					<Code2 class="h-3 w-3 mr-1.5" />
					Code
				</Button>
			</div>

			<div class="flex items-center justify-between">
				<div class="space-y-0.5">
					<Label for="auto-refresh" class="text-xs">Auto Refresh</Label>
					<p class="text-xs text-muted-foreground">Periodically update</p>
				</div>
				<Switch id="auto-refresh" v-model="localPeriodicallyRefetch">
					<RefreshCw class="h-3 w-3" />
				</Switch>
			</div>

			<div v-if="localPeriodicallyRefetch" class="space-y-1.5">
				<Label for="poll-interval" class="text-xs">Interval (seconds)</Label>
				<Input
					id="poll-interval"
					v-model.number="localPollInterval"
					type="number"
					:min="5"
					:max="300"
					class="h-8 text-xs"
					@keydown.stop
				/>
			</div>
		</div>

		<!-- IMAGE Layer Settings -->
		<div v-if="layer.type === 'IMAGE'">
			<Separator class="mb-3" />
			<ImageLayerEditor
				:layer="layer"
				@update="emit('update', $event)"
			/>
		</div>
	</div>
</template>
