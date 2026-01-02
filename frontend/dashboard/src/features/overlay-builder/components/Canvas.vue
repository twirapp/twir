<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import Moveable from 'vue3-moveable'
import type { OnDrag, OnResize, OnRotate } from 'vue3-moveable'

import HtmlLayerPreview from './HtmlLayerPreview.vue'
import type { AlignmentGuide, Layer } from '../types'

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
}>()

const canvasElement = ref<HTMLElement>()
const moveableRef = ref<InstanceType<typeof Moveable>>()

const canvasTransform = computed(() => {
	// Don't scale the canvas itself, let it stay at logical size
	// Only apply pan, Moveable will work with actual coordinates
	return `translate(${props.panX}px, ${props.panY}px)`
})

const gridStyle = computed(() => {
	if (!props.showGrid) return {}

	// Use logical grid size, no zoom multiplication needed
	const size = props.gridSize
	return {
		backgroundImage: `
			linear-gradient(to right, rgba(255, 255, 255, 0.05) 1px, transparent 1px),
			linear-gradient(to bottom, rgba(255, 255, 255, 0.05) 1px, transparent 1px)
		`,
		backgroundSize: `${size}px ${size}px`,
	}
})

const selectedLayers = computed(() => {
	return props.layers.filter(layer => props.selectedLayerIds.includes(layer.id))
})

const moveableTargets = computed(() => {
	return props.selectedLayerIds.map(id => `#layer-${id}`)
})

const isDragging = ref(false)

function handleCanvasClick(event: MouseEvent) {
	if (event.target === canvasElement.value) {
		emit('deselectAll')
	}
}

function handleLayerClick(layerId: string, event: MouseEvent) {
	event.stopPropagation()
	const layer = props.layers.find(l => l.id === layerId)
	if (layer?.locked) return

	const addToSelection = event.ctrlKey || event.metaKey
	emit('selectLayer', layerId, addToSelection)
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

	// Update the element transform with clamped position
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
	const visibility = layer.visible ? 'visible' : 'hidden'
	return {
		position: 'absolute' as const,
		left: '0px',
		top: '0px',
		width: `${layer.width}px`,
		height: `${layer.height}px`,
		transform: `translate(${layer.posX}px, ${layer.posY}px) rotate(${layer.rotation}deg)`,
		transformOrigin: 'center center',
		opacity: layer.opacity,
		visibility: visibility as 'visible' | 'hidden',
		zIndex: layer.zIndex,
		cursor: layer.locked ? 'not-allowed' : 'move',
	}
}

function isLayerSelected(layerId: string) {
	return props.selectedLayerIds.includes(layerId)
}

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
	>
		<div
			class="flex items-center justify-center w-full h-full p-8"
			:style="{
				transform: `scale(${zoom})`,
				transformOrigin: 'center',
			}"
		>
			<div
				ref="canvasElement"
				class="relative bg-[#121212] shadow-2xl border border-slate-700"
				:style="{
					width: `${canvasWidth}px`,
					height: `${canvasHeight}px`,
					transform: canvasTransform,
					transformOrigin: 'center',
					...gridStyle,
				}"
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
						'border-transparent hover:border-slate-500': !isLayerSelected(layer.id) && !layer.locked,
						'border-slate-700': layer.locked,
					}"
					:style="getLayerStyle(layer)"
					@click="handleLayerClick(layer.id, $event)"
				>
					<div class="w-full h-full overflow-hidden">
						<!-- HTML Layer Preview -->
						<HtmlLayerPreview
							v-if="layer.type === 'HTML'"
							:html="layer.settings?.htmlOverlayHtml"
							:css="layer.settings?.htmlOverlayCss"
							:js="layer.settings?.htmlOverlayJs"
							:width="layer.width"
							:height="layer.height"
							:refresh-interval="layer.settings?.htmlOverlayDataPollSecondsInterval"
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
			</div>
		</div>
	</div>
</template>
