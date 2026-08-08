<script setup lang="ts">
// oxlint-disable-next-line consistent-type-imports
import Moveable from 'vue3-moveable'

import HtmlLayerPreview from './HtmlLayerPreview.vue'
import ImageLayerPreview from './ImageLayerPreview.vue'
import LayerTypePreview from './LayerTypePreview.vue'
import type { Layer } from '../types'
import { type CanvasInteractionProps, useCanvasInteraction } from '../composables/useCanvasInteraction'
import { useWidgetPreviewMode } from '../composables/useWidgetPreviewMode'

import { Button } from '@/components/ui/button'

const props = defineProps<CanvasInteractionProps>()
const { t } = useI18n()
const { enabled: widgetsPreview } = useWidgetPreviewMode()

const emit = defineEmits<{
	updateLayer: [layerId: string, updates: Partial<Layer>]
	selectLayer: [layerId: string, addToSelection: boolean]
	deselectAll: []
	dragStart: [layerId: string]
	dragEnd: [layerId: string]
	findGuides: [layer: Layer]
	clearGuides: []
	toggleVisibility: [layerId: string]
	toggleLock: [layerId: string]
	removeLayer: [layerId: string]
	openLayerSettings: [layerId: string]
}>()

const {
	canvasElement,
	moveableRef,
	canvasStyle,
	canvasWrapperStyle,
	selectedLayers,
	moveableTargets,
	quickActionsLayer,
	quickActionsStyle,
	quickActionsInnerStyle,
	quickActionsOffsetStyle,
	handleCanvasMouseDown,
	handleCanvasClick,
	handleLayerClick,
	handleLayerMouseDown,
	onDrag,
	onDragEnd,
	onResize,
	onRotate,
	getLayerStyle,
	isLayerSelected,
} = useCanvasInteraction(props, emit)
</script>

<template>
	<div
		class="relative flex-1 overflow-hidden bg-slate-900"
		@click="handleCanvasClick"
		@mousedown="handleCanvasMouseDown"
	>
		<div class="flex items-center justify-center w-full h-full p-8">
			<div class="relative" :style="canvasWrapperStyle">
				<div
					ref="canvasElement"
					class="relative bg-[#121212] shadow-2xl border border-slate-700"
					:style="canvasStyle"
				>
					<div
						v-for="(guide, index) in alignmentGuides"
						:key="`guide-${index}`"
						class="absolute pointer-events-none"
						:class="{
							'border-l-2 border-blue-500': guide.type === 'vertical',
							'border-t-2 border-blue-500': guide.type === 'horizontal',
						}"
						:style="{
							left: guide.type === 'vertical' ? `${guide.position}px` : '0',
							top: guide.type === 'horizontal' ? `${guide.position}px` : '0',
							width: guide.type === 'vertical' ? '0' : '100%',
							height: guide.type === 'horizontal' ? '0' : '100%',
							zIndex: 9999,
						}"
					/>

					<div
						v-for="layer in layers"
						:id="`layer-${layer.id}`"
						:key="layer.id"
						class="absolute border-2 transition-colors"
						:class="{
							'border-primary bg-primary/5': isLayerSelected(layer.id),
							'border-transparent hover:border-slate-500': !isLayerSelected(layer.id) && !layer.locked && layer.visible,
							'border-dashed border-zinc-600 hover:border-zinc-500': !isLayerSelected(layer.id) && !layer.locked && !layer.visible,
							'border-slate-700': layer.locked && layer.visible,
							'border-dashed border-zinc-700': layer.locked && !layer.visible,
						}"
						:style="getLayerStyle(layer)"
						@click="handleLayerClick(layer.id, $event)"
						@mousedown="handleLayerMouseDown(layer.id, $event)"
					>
						<div class="w-full h-full overflow-hidden transition-opacity" :class="{ 'opacity-40': !layer.visible }">
							<HtmlLayerPreview
								v-if="layer.type === 'HTML'"
								:html="layer.settings?.htmlOverlayHtml"
								:css="layer.settings?.htmlOverlayCss"
								:js="layer.settings?.htmlOverlayJs"
								:width="layer.width"
								:height="layer.height"
								:refresh-interval="layer.settings?.htmlOverlayDataPollSecondsInterval"
							/>
							<ImageLayerPreview
								v-else-if="layer.type === 'IMAGE'"
								:image-url="layer.settings?.imageUrl"
								:width="layer.width"
								:height="layer.height"
							/>
							<LayerTypePreview
								v-else-if="layer.type === 'TEXT' || layer.type === 'VIDEO' || layer.type === 'IFRAME' || layer.type === 'YOUTUBE' || layer.type === 'EMOTE'"
								:layer="layer"
							/>
							<div v-else class="w-full h-full flex items-center justify-center">
								<slot name="layer-content" :layer="layer">
									<div class="text-white/50 text-center p-2">
										<p class="text-xs font-medium">{{ layer.name }}</p>
										<p class="text-xs mt-1">{{ layer.type }}</p>
									</div>
								</slot>
							</div>
						</div>

						<div
							v-if="!layer.visible"
							class="absolute left-1 top-1 flex items-center gap-1 rounded bg-slate-900/85 px-1.5 py-1 text-zinc-400 pointer-events-none"
						>
							<Icon name="lucide:eye-off" class="h-3 w-3" />
					<span class="text-[10px] leading-none">{{ t('overlayBuilder.canvas.hidden') }}</span>
						</div>

						<div
							v-if="isLayerSelected(layer.id)"
							class="absolute -top-6 left-0 px-2 py-1 bg-primary text-primary-foreground text-xs rounded pointer-events-none whitespace-nowrap"
						>
							{{ layer.name }}
						</div>

						<div
							v-if="layer.locked"
							class="absolute top-1 right-1 bg-slate-900/80 text-white p-1 rounded"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-3 w-3"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
								/>
							</svg>
						</div>
					</div>

					<Moveable
						v-if="selectedLayerIds.length > 0 && selectedLayers.every(l => !l.locked)"
						ref="moveableRef"
						:target="moveableTargets"
						:draggable="true"
						:resizable="true"
						:rotatable="true"
						:snappable="snapToGrid"
						:snapThreshold="5"
						:origin="false"
						:renderDirections="['nw', 'n', 'ne', 'w', 'e', 'sw', 's', 'se']"
						:keepRatio="false"
						:edge-draggable="false"
						:throttleDrag="0"
						:throttleResize="0"
						@drag="onDrag"
						@drag-end="onDragEnd"
						@resize="onResize"
						@rotate="onRotate"
					/>

					<div
						v-if="quickActionsLayer"
						class="pointer-events-none absolute"
						:style="quickActionsStyle"
					>
						<div :style="quickActionsInnerStyle">
							<div
								class="bg-popover pointer-events-auto flex items-center gap-0.5 rounded-md border p-1 shadow-lg"
								:style="quickActionsOffsetStyle"
								@mousedown.stop
								@click.stop
							>
								<Button
									variant="ghost"
									size="icon"
									class="h-7 w-7"
					:title="quickActionsLayer.visible ? t('overlayBuilder.layers.visibility.hide') : t('overlayBuilder.layers.visibility.show')"
									@click="emit('toggleVisibility', quickActionsLayer.id)"
								>
									<Icon v-if="quickActionsLayer.visible" name="lucide:eye" class="h-3.5 w-3.5" />
									<Icon v-else name="lucide:eye-off" class="text-muted-foreground h-3.5 w-3.5" />
								</Button>
								<Button
									variant="ghost"
									size="icon"
									class="h-7 w-7"
					:title="quickActionsLayer.locked ? t('overlayBuilder.layers.lock.unlock') : t('overlayBuilder.layers.lock.lock')"
									@click="emit('toggleLock', quickActionsLayer.id)"
								>
									<Icon v-if="!quickActionsLayer.locked" name="lucide:lock-open" class="h-3.5 w-3.5" />
									<Icon v-else name="lucide:lock" class="text-muted-foreground h-3.5 w-3.5" />
								</Button>
								<Button
									v-if="quickActionsLayer.type === 'IFRAME' && quickActionsLayer.settings.widgetKey"
									variant="ghost"
									size="icon"
									class="h-7 w-7"
									:class="{ 'bg-accent': widgetsPreview }"
									:title="t('overlayBuilder.toolbar.widgetsPreview')"
									@click="widgetsPreview = !widgetsPreview"
								>
									<Icon name="lucide:play" class="h-3.5 w-3.5" />
								</Button>
								<Button
									variant="ghost"
									size="icon"
									class="h-7 w-7"
									:title="t('overlayBuilder.layers.settings')"
									@click="emit('openLayerSettings', quickActionsLayer.id)"
								>
									<Icon name="lucide:settings-2" class="h-3.5 w-3.5" />
								</Button>
								<Button
									variant="ghost"
									size="icon"
									class="text-destructive hover:text-destructive h-7 w-7"
					:title="t('overlayBuilder.layers.delete')"
									@click="emit('removeLayer', quickActionsLayer.id)"
								>
									<Icon name="lucide:trash" class="h-3.5 w-3.5" />
								</Button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
