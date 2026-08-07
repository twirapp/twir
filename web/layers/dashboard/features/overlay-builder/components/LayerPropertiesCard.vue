<script setup lang="ts">
import { toRef } from 'vue'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Slider } from '@/components/ui/slider'
import { Switch } from '@/components/ui/switch'

import type { Layer } from '../types'
import { getLayerTypeMeta } from '../layer-type-meta'
import { useLayerProperties } from '../composables/useLayerProperties'

import ImageLayerEditor from './layer-editors/ImageLayerEditor.vue'
import LayerSettingsEditor from './layer-editors/LayerSettingsEditor.vue'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()
const { t } = useI18n()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
	openCodeEditor: []
}>()

const {
	fieldId,
	localName,
	localPosX,
	localPosY,
	localWidth,
	localHeight,
	localRotation,
	localOpacity,
	localVisible,
	localLocked,
	localPeriodicallyRefetch,
	localPollInterval,
} = useLayerProperties(toRef(props, 'layer'), (updates) => emit('update', updates))
</script>

<template>
	<Card
		id="layer-properties-card"
		class="flex min-h-0 flex-1 flex-col border-0 p-0 outline-none"
		tabindex="-1"
	>
		<div class="flex items-center justify-between border-b p-2">
			<CardTitle class="text-sm font-medium">{{ t('overlayBuilder.properties.title') }}</CardTitle>
			<span
				class="flex items-center gap-1.5 rounded-md px-2 py-1 text-[10px] font-medium"
				:class="getLayerTypeMeta(layer.type).chipClass"
			>
				<Icon :name="getLayerTypeMeta(layer.type).icon" class="h-3 w-3" />
				{{ t(getLayerTypeMeta(layer.type).labelKey) }}
			</span>
		</div>
		<CardContent class="min-h-0 flex-1 overflow-hidden p-0">
			<ScrollArea class="h-full">
				<div class="space-y-4 p-3">
					<!-- Layer Name -->
					<div class="space-y-2">
						<Label :for="fieldId('name')" class="text-xs">{{ t('overlayBuilder.properties.name') }}</Label>
						<Input
							:id="fieldId('name')"
							v-model="localName"
							:placeholder="t('overlayBuilder.properties.namePlaceholder')"
							class="h-8 text-xs"
							@keydown.stop
						/>
					</div>

					<Separator />

					<!-- Position & Size -->
					<div class="space-y-2">
						<span class="block text-[10px] font-medium uppercase tracking-wider text-muted-foreground">
							{{ t('overlayBuilder.properties.positionSize') }}
						</span>
						<div class="grid grid-cols-4 gap-1.5">
							<div class="space-y-1">
								<Label :for="fieldId('pos-x')" class="text-[10px] text-muted-foreground">X</Label>
								<Input
									:id="fieldId('pos-x')"
									v-model.number="localPosX"
									type="number"
									placeholder="0"
									class="h-7 px-1.5 text-xs"
									@keydown.stop
								/>
							</div>
							<div class="space-y-1">
								<Label :for="fieldId('pos-y')" class="text-[10px] text-muted-foreground">Y</Label>
								<Input
									:id="fieldId('pos-y')"
									v-model.number="localPosY"
									type="number"
									placeholder="0"
									class="h-7 px-1.5 text-xs"
									@keydown.stop
								/>
							</div>
							<div class="space-y-1">
								<Label :for="fieldId('width')" class="text-[10px] text-muted-foreground">{{ t('overlayBuilder.properties.width') }}</Label>
								<Input
									:id="fieldId('width')"
									v-model.number="localWidth"
									type="number"
									:min="1"
									placeholder="200"
									class="h-7 px-1.5 text-xs"
									@keydown.stop
								/>
							</div>
							<div class="space-y-1">
								<Label :for="fieldId('height')" class="text-[10px] text-muted-foreground">{{ t('overlayBuilder.properties.height') }}</Label>
								<Input
									:id="fieldId('height')"
									v-model.number="localHeight"
									type="number"
									:min="1"
									placeholder="200"
									class="h-7 px-1.5 text-xs"
									@keydown.stop
								/>
							</div>
						</div>
					</div>

					<!-- Rotation -->
					<div class="space-y-1.5">
						<div class="flex items-center justify-between">
						<Label :for="fieldId('rotation')" class="text-xs">{{ t('overlayBuilder.properties.rotation') }}</Label>
							<span class="text-xs text-muted-foreground">{{ localRotation }}°</span>
						</div>
						<Slider
							:id="fieldId('rotation')"
							:model-value="[localRotation]"
							:min="0"
							:max="360"
							:step="1"
							@update:model-value="(newValue) => {
								if (!newValue) return;
								localRotation = newValue[0] ?? 0
							}"
						/>
					</div>

					<!-- Opacity -->
					<div class="space-y-1.5">
						<div class="flex items-center justify-between">
							<Label :for="fieldId('opacity')" class="text-xs">{{ t('overlayBuilder.properties.opacity') }}</Label>
							<span class="text-xs text-muted-foreground">{{ Math.round(localOpacity) }}%</span>
						</div>
						<Slider
							:id="fieldId('opacity')"
							:model-value="[localOpacity]"
							:min="0"
							:max="100"
							:step="1"
							@update:model-value="(newValue) => {
								if (!newValue) return;
								localOpacity = newValue[0] ?? 0
							}"
						/>
					</div>

					<div class="flex items-center justify-between">
						<div class="space-y-0.5">
							<Label :for="fieldId('visible')" class="text-xs">{{ t('overlayBuilder.properties.visible') }}</Label>
							<p class="text-xs text-muted-foreground">{{ t('overlayBuilder.properties.visibleDescription') }}</p>
						</div>
						<Switch :id="fieldId('visible')" v-model="localVisible" />
					</div>

					<div class="flex items-center justify-between">
						<div class="space-y-0.5">
							<Label :for="fieldId('locked')" class="text-xs">{{ t('overlayBuilder.properties.locked') }}</Label>
							<p class="text-xs text-muted-foreground">{{ t('overlayBuilder.properties.lockedDescription') }}</p>
						</div>
						<Switch :id="fieldId('locked')" v-model="localLocked" />
					</div>

					<Separator />

					<!-- HTML Layer Settings -->
					<div v-if="layer.type === 'HTML'" class="space-y-3">
						<div class="flex items-center justify-between">
							<h4 class="text-xs font-medium">HTML</h4>
							<Button variant="outline" size="sm" class="h-7 text-xs" @click="emit('openCodeEditor')">
								<Icon name="lucide:code-xml" class="h-3 w-3 mr-1.5" />
								{{ t('overlayBuilder.properties.code') }}
							</Button>
						</div>

						<div class="flex items-center justify-between">
							<div class="space-y-0.5">
								<Label :for="fieldId('auto-refresh')" class="text-xs">{{ t('overlayBuilder.properties.autoRefresh') }}</Label>
								<p class="text-xs text-muted-foreground">{{ t('overlayBuilder.properties.autoRefreshDescription') }}</p>
							</div>
							<Switch :id="fieldId('auto-refresh')" v-model="localPeriodicallyRefetch" />
						</div>

						<div v-if="localPeriodicallyRefetch" class="space-y-1.5">
							<Label :for="fieldId('poll-interval')" class="text-xs">{{ t('overlayBuilder.properties.pollInterval') }}</Label>
							<Input
								:id="fieldId('poll-interval')"
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
					<ImageLayerEditor
						v-else-if="layer.type === 'IMAGE'"
						:layer="layer"
						@update="emit('update', $event)"
					/>

					<LayerSettingsEditor
						v-else-if="layer.type === 'TEXT' || layer.type === 'VIDEO' || layer.type === 'IFRAME' || layer.type === 'YOUTUBE' || layer.type === 'EMOTE'"
						:layer="layer"
						@update="emit('update', $event)"
					/>
				</div>
			</ScrollArea>
		</CardContent>
	</Card>
</template>
