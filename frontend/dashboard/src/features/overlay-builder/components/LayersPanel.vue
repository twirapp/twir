<script setup lang="ts">
import {
	ChevronDown,
	ChevronUp,
	Copy,
	Eye,
	EyeOff,
	Lock,
	LockOpen,
	Plus,
	Trash2,
} from 'lucide-vue-next'
import { ref, watch } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from '@/components/ui/tooltip'

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
}>()

const { t } = useI18n()

// Reverse layers for display (top layer shown first)
const displayLayers = ref<Layer[]>([])

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
	emit('select', layerId, addToSelection)
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
	<Card class="h-full flex flex-col border-0 p-0">
		<div class="border-b p-2 flex flex-row items-center justify-between space-y-0">
			<CardTitle class="text-sm font-medium">{{ t('overlaysRegistry.layers') || 'Layers' }}</CardTitle>
			<Button
				variant="default"
				size="sm"
				class="h-7 text-xs"
				@click="emit('addLayer')"
			>
				<Plus class="h-3 w-3 mr-1" />
				Add
			</Button>
		</div>
		<CardContent class="flex-1 p-0 overflow-hidden">
			<ScrollArea class="h-full">
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
					<div
						v-for="layer in displayLayers"
						:key="layer.id"
						class="layer-item group relative"
					>
						<div
							class="flex items-center gap-2 px-2 py-2 rounded-md border transition-all cursor-pointer drag-handle"
							:class="{
								'bg-accent border-primary': isLayerSelected(layer.id),
								'hover:bg-accent/50': !isLayerSelected(layer.id) && !layer.locked,
								'opacity-50': !layer.visible || layer.locked,
							}"
							@click="handleLayerClick(layer.id, $event)"
						>
							<!-- Layer Type Icon -->
							<span class="text-lg select-none">{{ getLayerTypeIcon(layer.type) }}</span>

							<!-- Layer Name -->
							<div class="flex-1 min-w-0">
								<p class="text-sm font-medium truncate">{{ layer.name }}</p>
								<p class="text-xs text-muted-foreground">
									{{ layer.width }}x{{ layer.height }}
								</p>
							</div>

							<!-- Actions -->
							<div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
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
												<Eye v-if="layer.visible" class="h-3.5 w-3.5" />
												<EyeOff v-else class="h-3.5 w-3.5 text-muted-foreground" />
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
												<LockOpen v-if="!layer.locked" class="h-3.5 w-3.5" />
												<Lock v-else class="h-3.5 w-3.5 text-muted-foreground" />
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
												<Copy class="h-3.5 w-3.5" />
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
												class="h-7 w-7 text-destructive hover:text-destructive"
												@click.stop="emit('remove', layer.id)"
											>
												<Trash2 class="h-3.5 w-3.5" />
											</Button>
										</TooltipTrigger>
										<TooltipContent>
											<p>Delete</p>
										</TooltipContent>
									</Tooltip>
								</TooltipProvider>
							</div>
						</div>

						<!-- Layer Order Controls (Only visible when selected) -->
						<div
							v-if="isLayerSelected(layer.id)"
							class="absolute -right-2 top-1/2 -translate-y-1/2 flex flex-col gap-0.5"
						>
							<TooltipProvider>
								<Tooltip>
									<TooltipTrigger as-child>
										<Button
											variant="secondary"
											size="icon"
											class="h-5 w-5 rounded-sm"
											@click.stop="emit('moveUp', layer.id)"
										>
											<ChevronUp class="h-3 w-3" />
										</Button>
									</TooltipTrigger>
									<TooltipContent>
										<p>Move Up</p>
									</TooltipContent>
								</Tooltip>
							</TooltipProvider>

							<TooltipProvider>
								<Tooltip>
									<TooltipTrigger as-child>
										<Button
											variant="secondary"
											size="icon"
											class="h-5 w-5 rounded-sm"
											@click.stop="emit('moveDown', layer.id)"
										>
											<ChevronDown class="h-3 w-3" />
										</Button>
									</TooltipTrigger>
									<TooltipContent>
										<p>Move Down</p>
									</TooltipContent>
								</Tooltip>
							</TooltipProvider>
						</div>
					</div>
				</VueDraggable>
			</ScrollArea>
		</CardContent>
	</Card>
</template>

<style scoped>
.layer-item {
	position: relative;
}

.drag-handle {
	cursor: grab;
}

.drag-handle:active {
	cursor: grabbing;
}
</style>
