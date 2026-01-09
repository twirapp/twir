<script setup lang="ts">
import {
	Copy,
	Eye,
	EyeOff,
	GripVertical,
	Lock,
	LockOpen,
	Plus,
	Trash2,
} from 'lucide-vue-next'
import { ref, watch } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'



import LayerPropertiesInline from './LayerPropertiesInline.vue'
import type { Layer } from '../types'

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
	addLayer: []
	updateLayerProperties: [layerId: string, updates: Partial<Layer>]
	openCodeEditor: []
}>()

// Reverse layers for display (top layer shown first)
const displayLayers = ref<Layer[]>([])

// Track expanded accordion items
const expandedLayerId = ref<string>()

// Watch for prop changes and update local ref
watch(() => props.layers, (newLayers) => {
	displayLayers.value = [...newLayers].reverse()
}, { immediate: true, deep: true })



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
			return '🌐'
		default:
			return '📄'
	}
}
</script>

<template>
	<UiCard class="h-full flex flex-col border-0 p-0">
		<div class="border-b p-2 flex flex-row items-center justify-between space-y-0">
			<UiCardTitle class="text-sm font-medium">Layers</UiCardTitle>
			<UiButton
				variant="default"
				size="sm"
				class="h-7 text-xs"
				@click="emit('addLayer')"
			>
				<Plus class="h-3 w-3 mr-1" />
				Add
			</UiButton>
		</div>
		<UiCardContent class="flex-1 p-0 overflow-hidden">
			<UiScrollArea class="h-full">
				<div v-if="layers.length === 0" class="p-8 text-center text-muted-foreground">
					<p class="text-sm">No layers yet</p>
					<p class="text-xs mt-1">Click "Add Layer" to get started</p>
				</div>
				<VueDraggable
					v-if="displayLayers.length > 0"
					v-model="displayLayers"
					:animation="150"
					handle=".drag-handle"
					ghost-class="opacity-30"
					class="p-2 space-y-1"
					@end="handleReorder"
				>
					<UiAccordion
						v-for="layer in displayLayers"
						:key="layer.id"
						type="single"
						collapsible
						:model-value="expandedLayerId === layer.id ? layer.id : undefined"
						class="layer-item"
					>
						<UiAccordionItem :value="layer.id" class="border-0">
							<div class="relative group">
								<div
									class="flex items-center gap-2 px-2 py-2 rounded-md border transition-all"
									:class="{
										'bg-accent border-primary': isLayerSelected(layer.id),
										'hover:bg-accent/50': !isLayerSelected(layer.id) && !layer.locked,
										'opacity-50': !layer.visible || layer.locked,
									}"
								>
									<!-- Drag Handle -->
									<div class="drag-handle cursor-grab active:cursor-grabbing">
										<GripVertical class="h-4 w-4 text-muted-foreground" />
									</div>

									<!-- Layer Type Icon -->
									<span
										class="text-lg select-none cursor-pointer"
										@click="handleLayerClick(layer.id, $event)"
									>
										{{ getLayerTypeIcon(layer.type) }}
									</span>

									<!-- Layer Name -->
									<div
										class="flex-1 min-w-0 cursor-pointer"
										@click="handleLayerClick(layer.id, $event)"
									>
										<p class="text-sm font-medium truncate">{{ layer.name }}</p>
										<p class="text-xs text-muted-foreground">
											{{ layer.width }}x{{ layer.height }}
										</p>
									</div>

									<!-- Actions -->
									<div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
										<!-- Visibility Toggle -->
										<UiTooltipProvider>
											<UiTooltip>
												<UiTooltipTrigger as-child>
													<UiButton
														variant="ghost"
														size="icon"
														class="h-7 w-7"
														@click.stop="emit('toggleVisibility', layer.id)"
													>
														<Eye v-if="layer.visible" class="h-3.5 w-3.5" />
														<EyeOff v-else class="h-3.5 w-3.5 text-muted-foreground" />
													</UiButton>
												</UiTooltipTrigger>
												<UiTooltipContent>
													<p>{{ layer.visible ? 'Hide' : 'Show' }}</p>
												</UiTooltipContent>
											</UiTooltip>
										</UiTooltipProvider>

										<!-- Lock Toggle -->
										<UiTooltipProvider>
											<UiTooltip>
												<UiTooltipTrigger as-child>
													<UiButton
														variant="ghost"
														size="icon"
														class="h-7 w-7"
														@click.stop="emit('toggleLock', layer.id)"
													>
														<LockOpen v-if="!layer.locked" class="h-3.5 w-3.5" />
														<Lock v-else class="h-3.5 w-3.5 text-muted-foreground" />
													</UiButton>
												</UiTooltipTrigger>
												<UiTooltipContent>
													<p>{{ layer.locked ? 'Unlock' : 'Lock' }}</p>
												</UiTooltipContent>
											</UiTooltip>
										</UiTooltipProvider>

										<!-- Duplicate -->
										<UiTooltipProvider>
											<UiTooltip>
												<UiTooltipTrigger as-child>
													<UiButton
														variant="ghost"
														size="icon"
														class="h-7 w-7"
														@click.stop="emit('duplicate', layer.id)"
													>
														<Copy class="h-3.5 w-3.5" />
													</UiButton>
												</UiTooltipTrigger>
												<UiTooltipContent>
													<p>Duplicate</p>
												</UiTooltipContent>
											</UiTooltip>
										</UiTooltipProvider>

										<!-- Delete -->
										<UiTooltipProvider>
											<UiTooltip>
												<UiTooltipTrigger as-child>
													<UiButton
														variant="ghost"
														size="icon"
														class="h-7 w-7 text-destructive hover:text-destructive"
														@click.stop="emit('remove', layer.id)"
													>
														<Trash2 class="h-3.5 w-3.5" />
													</UiButton>
												</UiTooltipTrigger>
												<UiTooltipContent>
													<p>Delete</p>
												</UiTooltipContent>
											</UiTooltip>
										</UiTooltipProvider>
									</div>
								</div>
							</div>

							<!-- Properties in Accordion Content -->
							<UiAccordionContent class="pt-2 pb-0">
								<div class="pl-6 pr-2">
									<LayerPropertiesInline
										:layer="layer"
										@update="emit('updateLayerProperties', layer.id, $event)"
										@open-code-editor="emit('openCodeEditor')"
									/>
								</div>
							</UiAccordionContent>
						</UiAccordionItem>
					</UiAccordion>
				</VueDraggable>
			</UiScrollArea>
		</UiCardContent>
	</UiCard>
</template>

<style scoped>
.layer-item {
	position: relative;
	margin-bottom: 0.25rem;
}
</style>
