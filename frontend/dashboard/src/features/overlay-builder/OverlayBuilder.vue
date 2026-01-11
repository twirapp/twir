<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

import BuilderToolbar from './components/BuilderToolbar.vue'
import LayersPanel from './components/LayersPanel.vue'
import Canvas from './components/Canvas.vue'
import CodeEditorDialog from './components/CodeEditorDialog.vue'
import OverlaySettings from './components/OverlaySettings.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import { useOverlayBuilder } from './composables/useOverlayBuilder'
import { type ChannelOverlayLayer, ChannelOverlayLayerType } from '@/gql/graphql'
import type { Layer, OverlayProject } from './types'

interface InitialProjectLayer {
	id: string
	type: ChannelOverlayLayer['type']
	name?: string
	posX: number
	posY: number
	width: number
	height: number
	rotation: number
	opacity?: number
	visible?: boolean
	locked?: boolean
	periodicallyRefetchData: boolean
	settings?: {
		htmlOverlayHtml?: string
		htmlOverlayCss?: string
		htmlOverlayJs?: string
		htmlOverlayDataPollSecondsInterval?: number
		imageUrl?: string
	}
}

interface Props {
	initialProject?: {
		id: string
		name: string
		width: number
		height: number
		instaSave?: boolean
		layers: InitialProjectLayer[]
	}
}

const props = defineProps<Props>()

const emit = defineEmits<{
	save: [project: OverlayProject]
	instantSave: [project: OverlayProject]
}>()

// Initialize builder
const builder = useOverlayBuilder()

// Overlay name state
const overlayName = ref('')
const instaSave = ref(false)
const canvasAreaRef = ref<HTMLElement>()
const loadedProjectId = ref<string>('')

// Auto-fit zoom calculation
function calculateFitZoom() {
	if (!canvasAreaRef.value) return

	// Get the actual DOM element from Vue component ref
	const canvasArea = (canvasAreaRef.value as any)?.$el as HTMLElement
	if (!canvasArea) return

	const availableWidth = canvasArea.clientWidth - 64 // padding
	const availableHeight = canvasArea.clientHeight - 64

	const scaleX = availableWidth / builder.project.width
	const scaleY = availableHeight / builder.project.height
	const fitZoom = Math.min(scaleX, scaleY) // Always fit to viewport

	// Set zoom to 80% of fit for more comfortable working space
	builder.canvasState.zoom = Math.max(0.1, fitZoom * 0.8)
}

// Recalculate zoom on mount and window resize only (canvas size is fixed at 1920x1080)

// Load initial project when it becomes available (handles async data loading)
function loadInitialProject() {
	if (!props.initialProject) return

	// Don't reload if it's the same project (prevents reset when editing name)
	if (loadedProjectId.value === props.initialProject.id && props.initialProject.id !== '') {
		return
	}

	loadedProjectId.value = props.initialProject.id
	overlayName.value = props.initialProject.name || ''
	instaSave.value = props.initialProject.instaSave || false

	const layers = props.initialProject.layers.map((layer, index) => {
		return {
			id: layer.id || `layer-${index}`,
			type: layer.type,
			name: layer.name || `Layer ${index + 1}`,
			posX: layer.posX,
			posY: layer.posY,
			width: layer.width,
			height: layer.height,
			rotation: Number(layer.rotation) || 0,
			opacity: layer.opacity || 1,
			visible: layer.visible !== undefined ? layer.visible : true,
			locked: layer.locked || false,
			zIndex: index,
			periodicallyRefetchData: layer.periodicallyRefetchData,
			settings: {
				htmlOverlayHtml: layer.settings?.htmlOverlayHtml || '',
				htmlOverlayCss: layer.settings?.htmlOverlayCss || '',
				htmlOverlayJs: layer.settings?.htmlOverlayJs || '',
				htmlOverlayDataPollSecondsInterval: layer.settings?.htmlOverlayDataPollSecondsInterval || 5,
				imageUrl: layer.settings?.imageUrl || '',
			},
		}
	})

	builder.loadProject({
		id: props.initialProject.id,
		name: props.initialProject.name,
		width: 1920,
		height: 1080,
		instaSave: props.initialProject.instaSave || false,
		layers,
	})

	// Recalculate zoom after loading
	nextTick(() => {
		calculateFitZoom()
	})
}

// Watch for initialProject changes (handles async loading)
// Only watch the ID to avoid reloading when other props change
watch(() => props.initialProject?.id, (newId) => {
	if (newId !== undefined) {
		loadInitialProject()
	}
}, { immediate: true })

// Watch for instaSave changes from props without reloading the whole project
watch(() => props.initialProject?.instaSave, (newInstaSave) => {
	// Only update if we have a loaded project and the value actually changed
	if (loadedProjectId.value && newInstaSave !== undefined && newInstaSave !== instaSave.value) {
		instaSave.value = newInstaSave
	}
})

// Load initial project on mount as well
onMounted(() => {
	loadInitialProject()
})

// Handle window resize for auto-fit zoom
onMounted(() => {
	window.addEventListener('resize', calculateFitZoom)
})

onUnmounted(() => {
	window.removeEventListener('resize', calculateFitZoom)
})

// Add layer dialog
const showAddLayerDialog = ref(false)

function handleAddLayer() {
	showAddLayerDialog.value = true
}

function addHtmlLayer() {
	builder.addLayer(ChannelOverlayLayerType.Html)
	showAddLayerDialog.value = false
}

function addImageLayer() {
	builder.addLayer(ChannelOverlayLayerType.Image)
	showAddLayerDialog.value = false
}

// Toolbar handlers
function handleSave() {
	const project = builder.exportProject()
	project.name = overlayName.value
	project.instaSave = instaSave.value

	emit('save', project)
}

async function handleLayerUpdate() {
	// Trigger instant save if enabled (debouncing handled by parent)
	if (instaSave.value) {
		// Wait for Vue to apply the layer updates
		await nextTick()

		const project = builder.exportProject()
		project.name = overlayName.value
		project.instaSave = instaSave.value

		emit('instantSave', project)
	}
}

// Watch for instaSave changes to save the setting immediately
watch(instaSave, (newValue, oldValue) => {
	// Only trigger if value actually changed and we have a loaded project
	if (newValue !== oldValue && loadedProjectId.value) {
		const project = builder.exportProject()
		project.name = overlayName.value
		project.instaSave = newValue

		emit('instantSave', project)
	}
})

// Canvas handlers
function handleUpdateLayer(layerId: string, updates: Partial<Layer>) {
	builder.updateLayer(layerId, updates)

	// Check if this is a position, size, rotation, opacity, or visibility update that should trigger instant save
	if (updates.posX !== undefined || updates.posY !== undefined || updates.rotation !== undefined || updates.width !== undefined || updates.height !== undefined || updates.opacity !== undefined || updates.visible !== undefined) {
		handleLayerUpdate()
	}
}

function handleSelectLayer(layerId: string, addToSelection: boolean) {
	builder.selectLayers([layerId], addToSelection)
}

function handleDeselectAll() {
	builder.deselectAll()
}

function handleFindGuides(layer: Layer) {
	builder.alignmentGuides.value = builder.findAlignmentGuides(layer)
}

function handleClearGuides() {
	builder.alignmentGuides.value = []
}

// Layers panel handlers
function handleLayerSelect(layerId: string, addToSelection: boolean) {
	builder.selectLayers([layerId], addToSelection)
}

function handleToggleVisibility(layerId: string) {
	const layer = builder.project.layers.find(l => l.id === layerId)
	if (layer) {
		builder.updateLayer(layerId, { visible: !layer.visible })
		// Trigger instant save for visibility change (async to ensure update is applied)
		handleLayerUpdate()
	}
}

function handleToggleLock(layerId: string) {
	const layer = builder.project.layers.find(l => l.id === layerId)
	if (layer) {
		builder.updateLayer(layerId, { locked: !layer.locked })
	}
}

function handleDuplicateLayer(layerId: string) {
	builder.duplicateLayer(layerId)
}

function handleRemoveLayer(layerId: string) {
	builder.removeLayer(layerId)
}

function handleMoveLayerUp(layerId: string) {
	builder.moveLayerUp(layerId)
}

function handleMoveLayerDown(layerId: string) {
	builder.moveLayerDown(layerId)
}

function handleReorderLayers(layers: Layer[]) {
	builder.reorderLayers(layers)
}

// Properties panel handlers
function handleUpdateLayerProperties(layerId: string, updates: Partial<Layer>) {
	builder.updateLayer(layerId, updates)

	// Check if this update should trigger instant save
	if (updates.posX !== undefined || updates.posY !== undefined || updates.rotation !== undefined || updates.width !== undefined || updates.height !== undefined || updates.opacity !== undefined || updates.visible !== undefined) {
		handleLayerUpdate()
	}
}

// Code editor dialog
const showCodeEditor = ref(false)
const editorLayer = ref<Layer | null>(null)

function handleOpenCodeEditor() {
	if (builder.activeLayer.value) {
		editorLayer.value = builder.activeLayer.value
		showCodeEditor.value = true
	}
}

function handleSaveCode(data: { html: string; css: string; js: string; refreshInterval: number }) {
	if (!editorLayer.value) return

	builder.updateLayer(editorLayer.value.id, {
		settings: {
			htmlOverlayHtml: data.html,
			htmlOverlayCss: data.css,
			htmlOverlayJs: data.js,
			htmlOverlayDataPollSecondsInterval: data.refreshInterval,
		},
	})
}

// Keyboard shortcuts
function handleKeyDown(event: KeyboardEvent) {
	// Check if focus is on an input element - if so, ignore most shortcuts
	const target = event.target as HTMLElement
	const isInputFocused = target.tagName === 'INPUT' ||
		target.tagName === 'TEXTAREA' ||
		target.isContentEditable

	// Ctrl/Cmd + S: Save (always works, even in inputs)
	if ((event.ctrlKey || event.metaKey) && event.key === 's') {
		event.preventDefault()
		handleSave()
	}
	// Ctrl/Cmd + Z: Undo (only when not in input)
	else if ((event.ctrlKey || event.metaKey) && event.key === 'z' && !event.shiftKey && !isInputFocused) {
		event.preventDefault()
		builder.undo()
	}
	// Ctrl/Cmd + Y or Ctrl/Cmd + Shift + Z: Redo (only when not in input)
	else if ((event.ctrlKey || event.metaKey) && (event.key === 'y' || (event.key === 'z' && event.shiftKey)) && !isInputFocused) {
		event.preventDefault()
		builder.redo()
	}
	// Ctrl/Cmd + C: Copy (only when not in input - let native copy work in inputs)
	else if ((event.ctrlKey || event.metaKey) && event.key === 'c' && !isInputFocused) {
		event.preventDefault()
		builder.copyToClipboard()
	}
	// Ctrl/Cmd + X: Cut (only when not in input - let native cut work in inputs)
	else if ((event.ctrlKey || event.metaKey) && event.key === 'x' && !isInputFocused) {
		event.preventDefault()
		builder.cutToClipboard()
	}
	// Ctrl/Cmd + V: Paste (only when not in input - let native paste work in inputs)
	else if ((event.ctrlKey || event.metaKey) && event.key === 'v' && !isInputFocused) {
		event.preventDefault()
		builder.pasteFromClipboard()
	}
	// Ctrl/Cmd + D: Duplicate (only when not in input)
	else if ((event.ctrlKey || event.metaKey) && event.key === 'd' && !isInputFocused) {
		event.preventDefault()
		if (builder.canvasState.selectedLayerIds.length > 0) {
			builder.duplicateLayers(builder.canvasState.selectedLayerIds)
		}
	}
	// Delete/Backspace: Delete selected layers (only if not typing in input)
	else if ((event.key === 'Delete' || event.key === 'Backspace') && !isInputFocused) {
		if (builder.canvasState.selectedLayerIds.length > 0) {
			event.preventDefault()
			builder.removeLayers(builder.canvasState.selectedLayerIds)
		}
	}
	// Ctrl/Cmd + A: Select all (only when not in input - let native select all work in inputs)
	else if ((event.ctrlKey || event.metaKey) && event.key === 'a' && !isInputFocused) {
		event.preventDefault()
		builder.selectAll()
	}
}

onMounted(() => {
	window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
	window.removeEventListener('keydown', handleKeyDown)
})

// Computed values for toolbar
const hasSelection = computed(() => builder.canvasState.selectedLayerIds.length > 0)
const canAlign = computed(() => builder.canvasState.selectedLayerIds.length >= 1)
const canDistribute = computed(() => builder.canvasState.selectedLayerIds.length >= 3)
</script>

<template>
	<div class="w-full h-full flex flex-col bg-background overflow-hidden">
		<!-- Toolbar -->
		<BuilderToolbar
			:can-undo="builder.canUndo.value"
			:can-redo="builder.canRedo.value"
			:has-selection="hasSelection"
			:can-align="canAlign"
			:can-distribute="canDistribute"
			:zoom="builder.canvasState.zoom"
			:show-grid="builder.canvasState.showGrid"
			:snap-to-grid="builder.canvasState.snapToGrid"
			:overlay-id="initialProject?.id"
			:overlay-name="overlayName"
			@save="handleSave"
			@undo="builder.undo"
			@redo="builder.redo"
			@copy="builder.copyToClipboard"
			@cut="builder.cutToClipboard"
			@paste="builder.pasteFromClipboard"
			@delete="builder.removeLayers(builder.canvasState.selectedLayerIds)"
			@duplicate="builder.duplicateLayers(builder.canvasState.selectedLayerIds)"
			@align-left="builder.alignLayers('left')"
			@align-center="builder.alignLayers('center')"
			@align-right="builder.alignLayers('right')"
			@align-top="builder.alignLayers('top')"
			@align-middle="builder.alignLayers('middle')"
			@align-bottom="builder.alignLayers('bottom')"
			@distribute-horizontal="builder.distributeLayersHorizontally"
			@distribute-vertical="builder.distributeLayersVertically"
			@zoom-in="builder.zoomIn"
			@zoom-out="builder.zoomOut"
			@reset-zoom="builder.resetZoom"
			@toggle-grid="builder.canvasState.showGrid = !builder.canvasState.showGrid"
			@toggle-snap="builder.canvasState.snapToGrid = !builder.canvasState.snapToGrid"
		/>

		<!-- Main Content -->
		<div class="flex-1 flex overflow-hidden">
			<!-- Canvas -->
			<Canvas
				ref="canvasAreaRef"
				:layers="builder.project.layers"
				:selected-layer-ids="builder.canvasState.selectedLayerIds"
				:zoom="builder.canvasState.zoom"
				:pan-x="builder.canvasState.panX"
				:pan-y="builder.canvasState.panY"
				:canvas-width="builder.project.width"
				:canvas-height="builder.project.height"
				:show-grid="builder.canvasState.showGrid"
				:snap-to-grid="builder.canvasState.snapToGrid"
				:grid-size="builder.canvasState.gridSize"
				:alignment-guides="builder.alignmentGuides.value"
				:snap-to-guides-enabled="builder.canvasState.showGuides"
				@update-layer="handleUpdateLayer"
				@select-layer="handleSelectLayer"
				@deselect-all="handleDeselectAll"
				@find-guides="handleFindGuides"
				@clear-guides="handleClearGuides"
			>
				<template #layer-content="{ layer }">
					<!-- Custom layer content rendering can be added here -->
					<div class="w-full h-full flex items-center justify-center text-white/70 text-sm">
						{{ layer.name }}
					</div>
				</template>
			</Canvas>

			<!-- Right Sidebar -->
			<div class="w-80 flex flex-col border-l">
				<!-- Overlay Settings -->
				<div class="border-b bg-background p-2">
					<OverlaySettings
						v-model:overlay-name="overlayName"
						v-model:insta-save="instaSave"
					/>
				</div>

				<!-- Layers Panel -->
				<div class="flex-1 min-h-0 overflow-hidden p-2">
					<LayersPanel
						:layers="builder.project.layers"
						:selected-layer-ids="builder.canvasState.selectedLayerIds"
						@select="handleLayerSelect"
						@toggle-visibility="handleToggleVisibility"
						@toggle-lock="handleToggleLock"
						@duplicate="handleDuplicateLayer"
						@remove="handleRemoveLayer"
						@move-up="handleMoveLayerUp"
						@move-down="handleMoveLayerDown"
						@reorder="handleReorderLayers"
						@add-layer="handleAddLayer"
						@update-layer-properties="handleUpdateLayerProperties"
						@open-code-editor="handleOpenCodeEditor"
					/>
				</div>
			</div>
		</div>

		<!-- Add Layer Dialog -->
		<Dialog v-model:open="showAddLayerDialog">
			<DialogContent class="sm:max-w-md">
				<DialogHeader>
					<DialogTitle>Add New Layer</DialogTitle>
					<DialogDescription>
						Choose a layer type to add to your overlay
					</DialogDescription>
				</DialogHeader>
				<div class="grid gap-4 py-4">
					<Button variant="outline" class="h-auto p-4 flex flex-col items-start" @click="addHtmlLayer">
						<div class="flex items-center gap-2 mb-2">
							<span class="text-2xl">🌐</span>
							<span class="font-semibold">HTML Layer</span>
						</div>
						<p class="text-sm text-muted-foreground text-left">
							Create a custom layer with HTML, CSS, and JavaScript
						</p>
					</Button>

					<Button variant="outline" class="h-auto p-4 flex flex-col items-start" @click="addImageLayer">
						<div class="flex items-center gap-2 mb-2">
							<span class="text-2xl">🖼️</span>
							<span class="font-semibold">Image Layer</span>
						</div>
						<p class="text-sm text-muted-foreground text-left">
							Display an image from a URL
						</p>
					</Button>
				</div>
			</DialogContent>
		</Dialog>

		<!-- Code Editor Dialog -->
		<CodeEditorDialog
			v-model:open="showCodeEditor"
			:layer-id="editorLayer?.id"
			:layer-name="editorLayer?.name"
			:html="editorLayer?.settings?.htmlOverlayHtml"
			:css="editorLayer?.settings?.htmlOverlayCss"
			:js="editorLayer?.settings?.htmlOverlayJs"
			:refresh-interval="editorLayer?.settings?.htmlOverlayDataPollSecondsInterval"
			@save="handleSaveCode"
		/>
	</div>
</template>
