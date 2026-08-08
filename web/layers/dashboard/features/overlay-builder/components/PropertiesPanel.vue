<script setup lang="ts">
import { toRef } from 'vue'
import { Card, CardContent, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Slider } from '@/components/ui/slider'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import ImageLayerEditor from './layer-editors/ImageLayerEditor.vue'
import LayerSettingsEditor from './layer-editors/LayerSettingsEditor.vue'

import type { Layer } from '../types'
import { useLayerProperties } from '../composables/useLayerProperties'

interface Props {
	layer: Layer | null
	multipleSelected: boolean
}

const props = defineProps<Props>()
const { t } = useI18n()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
	openCodeEditor: []
}>()

const {
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
	<Card class="h-full flex flex-col border-0 p-0">
		<div class="border-b p-2">
			<CardTitle class="text-sm font-medium">
				{{ t('overlayBuilder.propertyPanel.title') }}
			</CardTitle>
		</div>
		<CardContent class="flex-1 p-0 overflow-hidden">
			<ScrollArea class="h-full">
				<div v-if="!layer && !multipleSelected" class="p-4 text-center text-muted-foreground">
					<p class="text-sm">{{ t('overlayBuilder.propertyPanel.noLayer') }}</p>
					<p class="text-xs mt-1">{{ t('overlayBuilder.propertyPanel.noLayerHint') }}</p>
				</div>

				<div v-else-if="multipleSelected" class="p-4 text-center text-muted-foreground">
					<p class="text-sm">{{ t('overlayBuilder.propertyPanel.multipleLayers') }}</p>
					<p class="text-xs mt-1">{{ t('overlayBuilder.propertyPanel.multipleLayersHint') }}</p>
				</div>

				<div v-else-if="layer" class="p-4 space-y-6">
					<!-- Layer Info -->
					<div class="space-y-3">
						<div class="space-y-2">
							<Label for="layer-name">{{ t('overlayBuilder.propertyPanel.name') }}</Label>
							<Input id="layer-name" v-model="localName" :placeholder="t('overlayBuilder.propertyPanel.namePlaceholder')" @keydown.stop />
						</div>
					</div>

					<Separator />

					<!-- Transform -->
					<Tabs default-value="position" class="w-full">
						<TabsList class="grid w-full grid-cols-2">
							<TabsTrigger value="position">{{ t('overlayBuilder.propertyPanel.position') }}</TabsTrigger>
							<TabsTrigger value="appearance">{{ t('overlayBuilder.propertyPanel.appearance') }}</TabsTrigger>
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
										@keydown.stop
									/>
								</div>
								<div class="space-y-2">
									<Label for="pos-y">Y</Label>
									<Input
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
									<Label for="width">{{ t('overlayBuilder.propertyPanel.width') }}</Label>
									<Input
										id="width"
										v-model.number="localWidth"
										type="number"
										:min="1"
										placeholder="200"
										@keydown.stop
									/>
								</div>
								<div class="space-y-2">
									<Label for="height">{{ t('overlayBuilder.propertyPanel.height') }}</Label>
									<Input
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
									<Label for="rotation">{{ t('overlayBuilder.propertyPanel.rotation') }}</Label>
									<span class="text-xs text-muted-foreground">{{ localRotation }}°</span>
								</div>
								<Slider
									id="rotation"
									@update:model-value="(newValue) => {
										if (!newValue) return;
										localRotation = newValue[0] ?? 0
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
									<Label for="opacity">{{ t('overlayBuilder.propertyPanel.opacity') }}</Label>
									<span class="text-xs text-muted-foreground">{{ Math.round(localOpacity) }}%</span>
								</div>
								<Slider
									id="opacity"
									@update:model-value="(newValue) => {
										if (!newValue) return;
										localOpacity = newValue[0] ?? 0
									}"
									:model-value="[localOpacity]"
									:min="0"
									:max="100"
									:step="1"
								/>
							</div>

							<div class="flex items-center justify-between">
								<div class="space-y-0.5">
									<Label for="visible">{{ t('overlayBuilder.propertyPanel.visible') }}</Label>
									<p class="text-xs text-muted-foreground">{{ t('overlayBuilder.propertyPanel.visibleDescription') }}</p>
								</div>
								<Switch id="visible" v-model="localVisible">
									<Icon name="lucide:eye" class="h-4 w-4" />
								</Switch>
							</div>

							<div class="flex items-center justify-between">
								<div class="space-y-0.5">
									<Label for="locked">{{ t('overlayBuilder.propertyPanel.locked') }}</Label>
									<p class="text-xs text-muted-foreground">{{ t('overlayBuilder.propertyPanel.lockedDescription') }}</p>
								</div>
								<Switch id="locked" v-model="localLocked">
									<Icon name="lucide:lock" class="h-4 w-4" />
								</Switch>
							</div>
						</TabsContent>
					</Tabs>

					<Separator />

					<!-- HTML Layer Settings -->
					<div v-if="layer.type === 'HTML'" class="space-y-4">
						<div class="flex items-center justify-between">
							<h4 class="text-sm font-medium">{{ t('overlayBuilder.propertyPanel.htmlSettings') }}</h4>
							<Button variant="outline" size="sm" @click="emit('openCodeEditor')">
								<Icon name="lucide:code-xml" class="h-4 w-4 mr-2" />
								{{ t('overlayBuilder.propertyPanel.editCode') }}
							</Button>
						</div>

						<div class="flex items-center justify-between">
							<div class="space-y-0.5">
								<Label for="auto-refresh">{{ t('overlayBuilder.propertyPanel.autoRefresh') }}</Label>
								<p class="text-xs text-muted-foreground">{{ t('overlayBuilder.propertyPanel.autoRefreshDescription') }}</p>
							</div>
							<Switch id="auto-refresh" v-model="localPeriodicallyRefetch">
								<Icon name="lucide:refresh-cw" class="h-4 w-4" />
							</Switch>
						</div>

						<div v-if="localPeriodicallyRefetch" class="space-y-2">
							<Label for="poll-interval">{{ t('overlayBuilder.propertyPanel.updateInterval') }}</Label>
							<Input
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

					<LayerSettingsEditor
						v-if="layer.type === 'TEXT' || layer.type === 'VIDEO' || layer.type === 'IFRAME' || layer.type === 'YOUTUBE' || layer.type === 'EMOTE'"
						:layer="layer"
						@update="emit('update', $event)"
					/>
				</div>
			</ScrollArea>
		</CardContent>
	</Card>
</template>
