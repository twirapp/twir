<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'

import { Accordion, AccordionContent, AccordionItem } from '@/components/ui/accordion'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardTitle } from '@/components/ui/card'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { ChannelOverlayLayerType } from '~/gql/graphql.js'

import type { Layer } from '../types'

import LayerPropertiesInline from './LayerPropertiesInline.vue'

interface Props {
	layers: Layer[]
	selectedLayerIds: string[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
	select: [layerId: string, addToSelection: boolean]
	toggleVisibility: [layerId: string]
	toggleLock: [layerId: string]
	duplicate: [layerId: string]
	remove: [layerId: string]
	moveUp: [layerId: string]
	moveDown: [layerId: string]
	reorder: [layers: Layer[]]
	addLayer: [type: ChannelOverlayLayerType]
	updateLayerProperties: [layerId: string, updates: Partial<Layer>]
	openCodeEditor: []
}>()

// Reverse layers for display (top layer shown first)
const displayLayers = ref<Layer[]>([])

// Track expanded accordion items
const expandedLayerId = ref<string>()

// Add-layer popover
const isAddPopoverOpen = ref(false)

const layerTypeOptions: { type: ChannelOverlayLayerType; icon: string; label: string; description: string }[] = [
	{ type: ChannelOverlayLayerType.Image, icon: 'lucide:image', label: 'Картинка', description: 'Изображение из URL' },
	{ type: ChannelOverlayLayerType.Video, icon: 'lucide:video', label: 'Видео', description: 'Видео из URL' },
	{ type: ChannelOverlayLayerType.Youtube, icon: 'simple-icons:youtube', label: 'YouTube', description: 'Видео с YouTube' },
	{ type: ChannelOverlayLayerType.Text, icon: 'lucide:type', label: 'Текст', description: 'Текстовый слой' },
	{ type: ChannelOverlayLayerType.Html, icon: 'lucide:code-xml', label: 'HTML', description: 'HTML, CSS и JavaScript' },
	{ type: ChannelOverlayLayerType.Iframe, icon: 'lucide:panels-top-left', label: 'Виджет', description: 'Встраиваемый URL' },
	{ type: ChannelOverlayLayerType.Emote, icon: 'lucide:smile', label: 'Эмоции', description: 'Один эмоут на слой' },
]

function handleAddLayerType(type: ChannelOverlayLayerType) {
	isAddPopoverOpen.value = false
	emit('addLayer', type)
}

function expandLayer(layerId: string) {
	expandedLayerId.value = layerId
	nextTick(() => {
		document.getElementById(`layer-row-${layerId}`)?.scrollIntoView({ block: 'nearest' })
	})
}

defineExpose({ expandLayer })

// Watch for prop changes and update local ref
watch(
	() => props.layers,
	(newLayers) => {
		displayLayers.value = [...newLayers].reverse()
	},
	{ immediate: true, deep: true }
)

// Handle reordering when drag ends
function handleReorder() {
	// Reverse back to original order before emitting
	const newOrder = [...displayLayers.value].reverse()
	emit('reorder', newOrder)
}

function handleLayerClick(layerId: string, event: MouseEvent) {
	const addToSelection = event.ctrlKey || event.metaKey
	const wasSelected = isLayerSelected(layerId)

	emit('select', layerId, addToSelection)

	// Toggle accordion: close if already open and selected, open if not
	if (wasSelected && expandedLayerId.value === layerId) {
		expandedLayerId.value = undefined
	} else if (!addToSelection) {
		expandedLayerId.value = layerId
	}
}

function isLayerSelected(layerId: string) {
	return props.selectedLayerIds.includes(layerId)
}

function getLayerTypeIcon(type: string): string {
	switch (type) {
		case 'HTML':
			return 'lucide:code-xml'
		case 'IMAGE':
			return 'lucide:image'
		case 'TEXT':
			return 'lucide:type'
		case 'VIDEO':
			return 'lucide:video'
		case 'IFRAME':
			return 'lucide:panels-top-left'
		case 'YOUTUBE':
			return 'simple-icons:youtube'
		case 'EMOTE':
			return 'lucide:smile'
		default:
			return 'lucide:file'
	}
}
</script>

<template>
	<Card class="flex h-full flex-col border-0 p-0">
		<div class="flex flex-row items-center justify-between space-y-0 border-b p-2">
			<CardTitle class="text-sm font-medium">Layers</CardTitle>
			<Popover v-model:open="isAddPopoverOpen">
				<PopoverTrigger as-child>
					<Button
						variant="default"
						size="sm"
						class="h-7 text-xs"
					>
						<Icon
							name="lucide:plus"
							class="mr-1 h-3 w-3"
						/>
						Add
					</Button>
				</PopoverTrigger>
				<PopoverContent align="end" class="w-80 p-2">
					<div class="grid grid-cols-2 gap-1">
						<button
							v-for="option in layerTypeOptions"
							:key="option.type"
							type="button"
							class="hover:bg-accent flex flex-col items-start gap-1 rounded-md p-2 text-left transition-colors"
							@click="handleAddLayerType(option.type)"
						>
							<div class="flex items-center gap-1.5">
								<Icon :name="option.icon" class="h-4 w-4 shrink-0" />
								<span class="text-sm font-medium">{{ option.label }}</span>
							</div>
							<p class="text-muted-foreground text-xs">{{ option.description }}</p>
						</button>
					</div>
				</PopoverContent>
			</Popover>
		</div>
		<CardContent class="flex-1 overflow-hidden p-0">
			<ScrollArea class="h-full">
				<div
					v-if="layers.length === 0"
					class="text-muted-foreground p-8 text-center"
				>
					<p class="text-sm">No layers yet</p>
					<p class="mt-1 text-xs">Click "Add Layer" to get started</p>
				</div>
				<VueDraggable
					v-if="displayLayers.length > 0"
					v-model="displayLayers"
					:animation="150"
					handle=".drag-handle"
					ghost-class="opacity-30"
					class="space-y-1 p-2"
					@end="handleReorder"
				>
					<Accordion
						v-for="layer in displayLayers"
						:key="layer.id"
						type="single"
						collapsible
						:model-value="expandedLayerId === layer.id ? layer.id : undefined"
						class="layer-item"
					>
						<AccordionItem
							:value="layer.id"
							class="border-0"
						>
							<div :id="`layer-row-${layer.id}`" class="group relative">
								<div
									class="flex items-center gap-2 rounded-md border px-2 py-2 transition-all"
									:class="{
										'bg-accent border-primary': isLayerSelected(layer.id),
										'hover:bg-accent/50': !isLayerSelected(layer.id) && !layer.locked,
										'opacity-50': !layer.visible || layer.locked,
									}"
								>
									<!-- Drag Handle -->
									<div class="drag-handle cursor-grab active:cursor-grabbing">
										<Icon
											name="lucide:grip-vertical"
											class="text-muted-foreground h-4 w-4"
										/>
									</div>

									<!-- Layer Type Icon -->
									<span
										class="cursor-pointer text-lg select-none"
										@click="handleLayerClick(layer.id, $event)"
									>
										<Icon :name="getLayerTypeIcon(layer.type)" class="size-4" />
									</span>

									<!-- Layer Name -->
									<div
										class="min-w-0 flex-1 cursor-pointer"
										@click="handleLayerClick(layer.id, $event)"
									>
										<p class="truncate text-sm font-medium">{{ layer.name }}</p>
									<p class="text-muted-foreground text-xs">
										{{ layer.width }}x{{ layer.height }}
										<span v-if="!layer.visible" class="ml-2 rounded bg-muted px-1.5 py-0.5 text-[10px]">Скрыт</span>
									</p>
									</div>

									<!-- Actions -->
									<div
										class="flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100"
									>
										<!-- Visibility Toggle -->
										<TooltipProvider>
											<Tooltip>
												<TooltipTrigger as-child>
													<Button
														variant="ghost"
														size="icon"
														class="h-7 w-7"
														@click.stop="emit('toggleVisibility', layer.id)"
													>
														<Icon
															name="lucide:eye"
															v-if="layer.visible"
															class="h-3.5 w-3.5"
														/>
														<Icon
															name="lucide:eye-off"
															v-else
															class="text-muted-foreground h-3.5 w-3.5"
														/>
													</Button>
												</TooltipTrigger>
												<TooltipContent>
													<p>{{ layer.visible ? 'Hide' : 'Show' }}</p>
												</TooltipContent>
											</Tooltip>
										</TooltipProvider>

										<!-- Lock Toggle -->
										<TooltipProvider>
											<Tooltip>
												<TooltipTrigger as-child>
													<Button
														variant="ghost"
														size="icon"
														class="h-7 w-7"
														@click.stop="emit('toggleLock', layer.id)"
													>
														<Icon
															name="lucide:lock-open"
															v-if="!layer.locked"
															class="h-3.5 w-3.5"
														/>
														<Icon
															name="lucide:lock"
															v-else
															class="text-muted-foreground h-3.5 w-3.5"
														/>
													</Button>
												</TooltipTrigger>
												<TooltipContent>
													<p>{{ layer.locked ? 'Unlock' : 'Lock' }}</p>
												</TooltipContent>
											</Tooltip>
										</TooltipProvider>

										<!-- Duplicate -->
										<TooltipProvider>
											<Tooltip>
												<TooltipTrigger as-child>
													<Button
														variant="ghost"
														size="icon"
														class="h-7 w-7"
														@click.stop="emit('duplicate', layer.id)"
													>
														<Icon
															name="lucide:copy"
															class="h-3.5 w-3.5"
														/>
													</Button>
												</TooltipTrigger>
												<TooltipContent>
													<p>Duplicate</p>
												</TooltipContent>
											</Tooltip>
										</TooltipProvider>

										<!-- Delete -->
										<TooltipProvider>
											<Tooltip>
												<TooltipTrigger as-child>
													<Button
														variant="ghost"
														size="icon"
														class="text-destructive hover:text-destructive h-7 w-7"
														@click.stop="emit('remove', layer.id)"
													>
														<Icon
															name="lucide:trash"
															class="h-3.5 w-3.5"
														/>
													</Button>
												</TooltipTrigger>
												<TooltipContent>
													<p>Delete</p>
												</TooltipContent>
											</Tooltip>
										</TooltipProvider>
									</div>
								</div>
							</div>

							<!-- Properties in Accordion Content -->
							<AccordionContent class="pt-2 pb-0">
								<div class="pr-2 pl-6">
									<LayerPropertiesInline
										:layer="layer"
										@update="emit('updateLayerProperties', layer.id, $event)"
										@open-code-editor="emit('openCodeEditor')"
									/>
								</div>
							</AccordionContent>
						</AccordionItem>
					</Accordion>
				</VueDraggable>
			</ScrollArea>
		</CardContent>
	</Card>
</template>

<style scoped>
.layer-item {
	position: relative;
	margin-bottom: 0.25rem;
}
</style>
