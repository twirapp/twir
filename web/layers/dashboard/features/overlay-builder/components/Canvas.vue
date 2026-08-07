<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
// oxlint-disable-next-line consistent-type-imports
import Moveable from 'vue3-moveable'
import type { OnDrag, OnResize, OnRotate } from 'vue3-moveable'

import HtmlLayerPreview from './HtmlLayerPreview.vue'
import ImageLayerPreview from './ImageLayerPreview.vue'
import LayerTypePreview from './LayerTypePreview.vue'
import { getLayerTypeMeta } from '../layer-type-meta'
import type { AlignmentGuide, Layer } from '../types'

import { Button } from '@/components/ui/button'

interface Props {
	layers: Layer[]
	selectedLayerIds: string[]
	zoom: number
	panX: number
	panY: number
	canvasWidth: number
	canvasHeight: number
	showGrid: boolean
	snapToGrid: boolean
	gridSize: number
	alignmentGuides: AlignmentGuide[]
	snapToGuidesEnabled: boolean
}

const props = defineProps<Props>()

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

const canvasElement = ref<HTMLElement>()
const moveableRef = ref<InstanceType<typeof Moveable>>()

// Snapping state - tracks if layer is currently snapped
const snappedState = ref<{
	layerId: string | null
	snappedX: boolean
	snappedY: boolean
	snapPositionX: number | null
	snapPositionY: number | null
}>({
	layerId: null,
	snappedX: false,
	snappedY: false,
	snapPositionX: null,
	snapPositionY: null,
})

// Canvas and grid wrapper transformation
const canvasTransform = computed(() => {
	return `translate(${props.panX}px, ${props.panY}px)`
})

// Grid style applied directly to canvas background
const canvasStyle = computed(() => {
	const baseStyle = {
		width: `${props.canvasWidth}px`,
		height: `${props.canvasHeight}px`,
		transform: canvasTransform.value,
		transformOrigin: 'center',
		zIndex: 1,
	}

	if (!props.showGrid) {
		return baseStyle
	}

	// Add grid background
	const gridSize = props.gridSize
	return {
		...baseStyle,
		backgroundImage: `
			linear-gradient(to right, rgba(255, 255, 255, 0.1) 1px, transparent 1px),
			linear-gradient(to bottom, rgba(255, 255, 255, 0.1) 1px, transparent 1px)
		`,
		backgroundSize: `${gridSize}px ${gridSize}px`,
	}
})

// Canvas wrapper transformation (applies zoom)
const canvasWrapperStyle = computed(() => {
	return {
		transform: `scale(${props.zoom})`,
		transformOrigin: 'center',
	}
})

const selectedLayers = computed(() => {
	return props.layers.filter(layer => props.selectedLayerIds.includes(layer.id))
})

const moveableTargets = computed(() => {
	return props.selectedLayerIds.map(id => `#layer-${id}`)
})

const isDragging = ref(false)

// Deselect only when the whole press-release gesture started on empty canvas:
// selecting on mousedown can retarget the subsequent click to the canvas element
// (Moveable's drag area covers the layer mid-gesture), which must not clear selection.
let canvasMouseDownTarget: EventTarget | null = null

function handleCanvasMouseDown(event: MouseEvent) {
	canvasMouseDownTarget = event.target
}

function handleCanvasClick(event: MouseEvent) {
	if (event.target === canvasElement.value && canvasMouseDownTarget === canvasElement.value) {
		emit('deselectAll')
	}
}

// Local snapping function with hysteresis for sticky feel
function snapToGuides(layer: Layer, posX: number, posY: number): { x: number; y: number } {
	if (!props.snapToGuidesEnabled) return { x: posX, y: posY }

	const snapThreshold = 8 // Distance to snap in
	const releaseThreshold = 15 // Distance to release snap (higher = more sticky)
	let snappedX = posX
	let snappedY = posY
	let newSnapX = false
	let newSnapY = false
	let newSnapPosX: number | null = null
	let newSnapPosY: number | null = null

	const layerCenterX = posX + layer.width / 2
	const layerCenterY = posY + layer.height / 2
	const layerRight = posX + layer.width
	const layerBottom = posY + layer.height

	const canvasCenterX = props.canvasWidth / 2
	const canvasCenterY = props.canvasHeight / 2

	// Check if we're continuing a previous snap for this layer
	const isSameLayer = snappedState.value.layerId === layer.id
	const wasSnappedX = isSameLayer && snappedState.value.snappedX
	const wasSnappedY = isSameLayer && snappedState.value.snappedY

	// Vertical snapping with hysteresis
	const checkVerticalSnap = (targetPos: number, snapTo: number) => {
		const distance = Math.abs(targetPos - snapTo)
		if (wasSnappedX && snappedState.value.snapPositionX === snapTo) {
			// Already snapped - need more distance to release
			if (distance < releaseThreshold) {
				return true
			}
		} else {
			// Not snapped yet - easier to snap
			if (distance < snapThreshold) {
				return true
			}
		}
		return false
	}

	// Horizontal snapping with hysteresis
	const checkHorizontalSnap = (targetPos: number, snapTo: number) => {
		const distance = Math.abs(targetPos - snapTo)
		if (wasSnappedY && snappedState.value.snapPositionY === snapTo) {
			// Already snapped - need more distance to release
			if (distance < releaseThreshold) {
				return true
			}
		} else {
			// Not snapped yet - easier to snap
			if (distance < snapThreshold) {
				return true
			}
		}
		return false
	}

	// Snap to canvas edges and center - vertical
	if (checkVerticalSnap(posX, 0)) {
		snappedX = 0
		newSnapX = true
		newSnapPosX = 0
	} else if (checkVerticalSnap(layerCenterX, canvasCenterX)) {
		snappedX = canvasCenterX - layer.width / 2
		newSnapX = true
		newSnapPosX = canvasCenterX
	} else if (checkVerticalSnap(layerRight, props.canvasWidth)) {
		snappedX = props.canvasWidth - layer.width
		newSnapX = true
		newSnapPosX = props.canvasWidth
	}

	// Snap to canvas edges and center - horizontal
	if (checkHorizontalSnap(posY, 0)) {
		snappedY = 0
		newSnapY = true
		newSnapPosY = 0
	} else if (checkHorizontalSnap(layerCenterY, canvasCenterY)) {
		snappedY = canvasCenterY - layer.height / 2
		newSnapY = true
		newSnapPosY = canvasCenterY
	} else if (checkHorizontalSnap(layerBottom, props.canvasHeight)) {
		snappedY = props.canvasHeight - layer.height
		newSnapY = true
		newSnapPosY = props.canvasHeight
	}

	// Snap to other layers
	props.layers.forEach((other) => {
		if (other.id === layer.id || !other.visible) return

		const otherCenterX = other.posX + other.width / 2
		const otherCenterY = other.posY + other.height / 2
		const otherRight = other.posX + other.width
		const otherBottom = other.posY + other.height

		// Vertical snapping to other layers
		if (!newSnapX) {
			if (checkVerticalSnap(posX, other.posX)) {
				snappedX = other.posX
				newSnapX = true
				newSnapPosX = other.posX
			} else if (checkVerticalSnap(layerRight, otherRight)) {
				snappedX = otherRight - layer.width
				newSnapX = true
				newSnapPosX = otherRight
			} else if (checkVerticalSnap(layerCenterX, otherCenterX)) {
				snappedX = otherCenterX - layer.width / 2
				newSnapX = true
				newSnapPosX = otherCenterX
			}
		}

		// Horizontal snapping to other layers
		if (!newSnapY) {
			if (checkHorizontalSnap(posY, other.posY)) {
				snappedY = other.posY
				newSnapY = true
				newSnapPosY = other.posY
			} else if (checkHorizontalSnap(layerBottom, otherBottom)) {
				snappedY = otherBottom - layer.height
				newSnapY = true
				newSnapPosY = otherBottom
			} else if (checkHorizontalSnap(layerCenterY, otherCenterY)) {
				snappedY = otherCenterY - layer.height / 2
				newSnapY = true
				newSnapPosY = otherCenterY
			}
		}
	})

	// Update snapped state
	snappedState.value = {
		layerId: layer.id,
		snappedX: newSnapX,
		snappedY: newSnapY,
		snapPositionX: newSnapPosX,
		snapPositionY: newSnapPosY,
	}

	return { x: snappedX, y: snappedY }
}

function handleLayerClick(layerId: string, event: MouseEvent) {
	event.stopPropagation()
	const layer = props.layers.find(l => l.id === layerId)
	if (layer?.locked) return

	const addToSelection = event.ctrlKey || event.metaKey
	// Selection already applied on mousedown - skip resetting Moveable's target mid-gesture
	if (!addToSelection && props.selectedLayerIds.length === 1 && props.selectedLayerIds[0] === layerId) {
		return
	}
	emit('selectLayer', layerId, addToSelection)
}

// Select on press so an unselected layer can be dragged in a single gesture:
// after Vue retargets Moveable, kick off its drag with the same mousedown event.
function handleLayerMouseDown(layerId: string, event: MouseEvent) {
	if (event.button !== 0) return
	const layer = props.layers.find(l => l.id === layerId)
	if (!layer || layer.locked) return
	if (isLayerSelected(layerId)) return

	emit('selectLayer', layerId, event.ctrlKey || event.metaKey)
	nextTick(() => {
		moveableRef.value?.dragStart(event)
	})
}

function onDrag(e: OnDrag) {
	isDragging.value = true
	const target = e.target as HTMLElement
	const layerId = target.id.replace('layer-', '')
	const layer = props.layers.find(l => l.id === layerId)
	if (!layer || layer.locked) return

	// e.translate is already in logical coordinates since canvas is not scaled
	let newPosX = Math.round(e.translate[0])
	let newPosY = Math.round(e.translate[1])

	// Calculate bounds
	const maxX = props.canvasWidth - layer.width
	const maxY = props.canvasHeight - layer.height

	// Clamp position within canvas bounds
	newPosX = Math.max(0, Math.min(newPosX, maxX))
	newPosY = Math.max(0, Math.min(newPosY, maxY))

	// Apply snapping
	const snappedPos = snapToGuides(layer, newPosX, newPosY)
	newPosX = snappedPos.x
	newPosY = snappedPos.y

	// Update the element transform with snapped position
	target.style.transform = `translate(${newPosX}px, ${newPosY}px) rotate(${layer.rotation}deg)`

	// Update layer position
	emit('updateLayer', layerId, {
		posX: newPosX,
		posY: newPosY,
	})

	const updatedLayer = { ...layer, posX: newPosX, posY: newPosY }
	emit('findGuides', updatedLayer)
}

function onDragEnd() {
	isDragging.value = false
	emit('clearGuides')

	// Reset snapping state when drag ends
	snappedState.value = {
		layerId: null,
		snappedX: false,
		snappedY: false,
		snapPositionX: null,
		snapPositionY: null,
	}
}

function onResize(e: OnResize) {
	const target = e.target as HTMLElement
	const layerId = target.id.replace('layer-', '')
	const layer = props.layers.find(l => l.id === layerId)
	if (!layer || layer.locked) return

	// e.width and e.drag.translate are in logical coordinates
	const minSize = 10
	let width = Math.round(e.width)
	let height = Math.round(e.height)
	let posX = Math.round(e.drag.translate[0])
	let posY = Math.round(e.drag.translate[1])

	// Ensure minimum size
	width = Math.max(minSize, width)
	height = Math.max(minSize, height)

	// Ensure size doesn't exceed canvas bounds
	width = Math.min(width, props.canvasWidth)
	height = Math.min(height, props.canvasHeight)

	// Ensure position + size stays within canvas bounds
	const maxX = props.canvasWidth - width
	const maxY = props.canvasHeight - height
	posX = Math.max(0, Math.min(posX, maxX))
	posY = Math.max(0, Math.min(posY, maxY))

	// Update element styles with clamped values
	target.style.width = `${width}px`
	target.style.height = `${height}px`
	target.style.transform = `translate(${posX}px, ${posY}px) rotate(${layer.rotation}deg)`

	// Update layer
	emit('updateLayer', layerId, {
		width,
		height,
		posX,
		posY,
	})
}

function onRotate(e: OnRotate) {
	const target = e.target as HTMLElement
	const layerId = target.id.replace('layer-', '')
	const layer = props.layers.find(l => l.id === layerId)
	if (!layer || layer.locked) return

	const rotation = Math.round(e.rotate)
	let posX = Math.round(e.drag.translate[0])
	let posY = Math.round(e.drag.translate[1])

	// Constrain position within canvas bounds during rotation
	const maxX = props.canvasWidth - layer.width
	const maxY = props.canvasHeight - layer.height
	posX = Math.max(0, Math.min(posX, maxX))
	posY = Math.max(0, Math.min(posY, maxY))

	// Update element transform with clamped values
	target.style.transform = `translate(${posX}px, ${posY}px) rotate(${rotation}deg)`

	// Update layer
	emit('updateLayer', layerId, {
		rotation,
		posX,
		posY,
	})
}

function getLayerStyle(layer: Layer) {
	// Hidden layers stay rendered as ghost placeholders so they remain selectable
	return {
		position: 'absolute' as const,
		left: '0px',
		top: '0px',
		width: `${layer.width}px`,
		height: `${layer.height}px`,
		transform: `translate(${layer.posX}px, ${layer.posY}px) rotate(${layer.rotation}deg)`,
		transformOrigin: 'center center',
		opacity: layer.opacity,
		zIndex: layer.zIndex,
		cursor: layer.locked ? 'not-allowed' : 'move',
	}
}

function isLayerSelected(layerId: string) {
	return props.selectedLayerIds.includes(layerId)
}

// Quick actions floating bar for the single selected layer.
// Rendered inside the scaled canvas, so position is in canvas coordinates and
// the bar content is counter-scaled to keep a constant on-screen size.
const quickActionsLayer = computed(() => {
	if (selectedLayers.value.length !== 1) return null
	return selectedLayers.value[0] ?? null
})

// Estimated on-screen size of the bar, only used for bounds clamping
const QUICK_ACTIONS_WIDTH = 150
const QUICK_ACTIONS_HEIGHT = 36
const QUICK_ACTIONS_GAP = 8

const quickActionsFlipAbove = computed(() => {
	const layer = quickActionsLayer.value
	if (!layer) return false
	const gap = QUICK_ACTIONS_GAP / props.zoom
	const barHeight = QUICK_ACTIONS_HEIGHT / props.zoom
	const spaceBelow = props.canvasHeight - (layer.posY + layer.height)
	return spaceBelow < gap + barHeight && layer.posY >= gap + barHeight
})

const quickActionsStyle = computed(() => {
	const layer = quickActionsLayer.value
	if (!layer) return {}
	const barWidth = QUICK_ACTIONS_WIDTH / props.zoom
	const left = Math.max(0, Math.min(layer.posX, props.canvasWidth - barWidth))
	const top = quickActionsFlipAbove.value ? layer.posY : layer.posY + layer.height
	return {
		left: `${left}px`,
		top: `${top}px`,
		zIndex: 10000,
	}
})

const quickActionsInnerStyle = computed(() => {
	return {
		transform: `scale(${1 / props.zoom})`,
		transformOrigin: '0 0',
	}
})

const quickActionsOffsetStyle = computed(() => {
	return {
		transform: quickActionsFlipAbove.value
			? 'translateY(calc(-100% - 8px))'
			: 'translateY(8px)',
	}
})

function handleKeyDown(event: KeyboardEvent) {
	if (event.key === 'Delete' || event.key === 'Backspace') {
		if (props.selectedLayerIds.length > 0) {
			event.preventDefault()
		}
	}
}

// Watch for external layer changes and update Moveable position
watch(() => props.layers.map(l => ({ id: l.id, posX: l.posX, posY: l.posY, width: l.width, height: l.height, rotation: l.rotation })), () => {
	// Only update if not currently dragging
	if (!isDragging.value && moveableRef.value) {
		nextTick(() => {
			moveableRef.value?.updateRect()
		})
	}
}, { deep: true })

onMounted(() => {
	window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
	window.removeEventListener('keydown', handleKeyDown)
})
</script>

<template>
	<div
		class="relative flex-1 overflow-hidden bg-slate-900"
		@click="handleCanvasClick"
		@mousedown="handleCanvasMouseDown"
	>
		<!-- Container for canvas and grid -->
		<div class="flex items-center justify-center w-full h-full p-8">
			<!-- Wrapper that scales canvas -->
			<div class="relative" :style="canvasWrapperStyle">
				<!-- Main canvas with grid background -->
				<div
					ref="canvasElement"
					class="relative bg-[#121212] shadow-2xl border border-slate-700"
					:style="canvasStyle"
				>
					<!-- Alignment guides -->
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
					<div class="w-full h-full overflow-hidden">
						<!-- Hidden layer: gray ghost placeholder, no content preview -->
						<div
							v-if="!layer.visible"
							class="flex h-full w-full flex-col items-center justify-center gap-1 bg-white/5 p-2 text-zinc-500"
						>
							<Icon
								:name="getLayerTypeMeta(layer.type).icon"
								class="h-5 w-5 flex-none"
							/>
							<span class="max-w-full truncate text-xs">{{ layer.name }}</span>
						</div>
						<!-- HTML Layer Preview -->
						<HtmlLayerPreview
							v-else-if="layer.type === 'HTML'"
							:html="layer.settings?.htmlOverlayHtml"
							:css="layer.settings?.htmlOverlayCss"
							:js="layer.settings?.htmlOverlayJs"
							:width="layer.width"
							:height="layer.height"
							:refresh-interval="layer.settings?.htmlOverlayDataPollSecondsInterval"
						/>
						<!-- IMAGE Layer Preview -->
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
						<!-- Fallback content -->
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

				<!-- Quick actions for the single selected layer -->
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
								:title="quickActionsLayer.visible ? 'Hide' : 'Show'"
								@click="emit('toggleVisibility', quickActionsLayer.id)"
							>
								<Icon
									v-if="quickActionsLayer.visible"
									name="lucide:eye"
									class="h-3.5 w-3.5"
								/>
								<Icon
									v-else
									name="lucide:eye-off"
									class="text-muted-foreground h-3.5 w-3.5"
								/>
							</Button>
							<Button
								variant="ghost"
								size="icon"
								class="h-7 w-7"
								:title="quickActionsLayer.locked ? 'Unlock' : 'Lock'"
								@click="emit('toggleLock', quickActionsLayer.id)"
							>
								<Icon
									v-if="!quickActionsLayer.locked"
									name="lucide:lock-open"
									class="h-3.5 w-3.5"
								/>
								<Icon
									v-else
									name="lucide:lock"
									class="text-muted-foreground h-3.5 w-3.5"
								/>
							</Button>
							<Button
								variant="ghost"
								size="icon"
								class="h-7 w-7"
								title="Settings"
								@click="emit('openLayerSettings', quickActionsLayer.id)"
							>
								<Icon
									name="lucide:settings-2"
									class="h-3.5 w-3.5"
								/>
							</Button>
							<Button
								variant="ghost"
								size="icon"
								class="text-destructive hover:text-destructive h-7 w-7"
								title="Delete"
								@click="emit('removeLayer', quickActionsLayer.id)"
							>
								<Icon
									name="lucide:trash"
									class="h-3.5 w-3.5"
								/>
							</Button>
						</div>
					</div>
				</div>
				</div>
			</div>
		</div>
	</div>
</template>
