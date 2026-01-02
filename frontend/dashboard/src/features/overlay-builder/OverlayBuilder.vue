<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'

import BuilderToolbar from './components/BuilderToolbar.vue'
import LayersPanel from './components/LayersPanel.vue'
import PropertiesPanel from './components/PropertiesPanel.vue'
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
import { ChannelOverlayLayerType } from '@/gql/graphql'
import type { Layer } from './types'

interface Props {
	initialProject?: {
		id: string
		name: string
		width: number
		height: number
		layers: any[]
	}
}

const props = defineProps<Props>()

const emit = defineEmits<{
	save: [project: any]
}>()

// Initialize builder
const builder = useOverlayBuilder()

// Overlay name state
const overlayName = ref('')

// Load initial project if provided
onMounted(() => {
	if (props.initialProject) {
		overlayName.value = props.initialProject.name || ''
		const layers = props.initialProject.layers.map((layer, index) => ({
			id: layer.id || `layer-${index}`,
			type: layer.type,
			name: layer.name || `Layer ${index + 1}`,
			posX: layer.posX,
			posY: layer.posY,
			width: layer.width,
			height: layer.height,
			rotation: layer.rotation || 0,
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
			},
		}))

		builder.loadProject({
			id: props.initialProject.id,
			name: props.initialProject.name,
			width: props.initialProject.width,
			height: props.initialProject.height,
			layers,
		})
	}
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

// Toolbar handlers
function handleSave() {
	const project = {
		...builder.exportProject(),
		name: overlayName.value,
	}
	emit('save', project)
}

// Canvas handlers
function handleUpdateLayer(layerId: string, updates: Partial<Layer>) {
	builder.updateLayer(layerId, updates)
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
function handleUpdateLayerProperties(updates: Partial<Layer>) {
	if (builder.activeLayer.value) {
		builder.updateLayer(builder.activeLayer.value.id, updates)
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
	// Ctrl/Cmd + S: Save
	if ((event.ctrlKey || event.metaKey) && event.key === 's') {
		event.preventDefault()
		handleSave()
	}
	// Ctrl/Cmd + Z: Undo
	else if ((event.ctrlKey || event.metaKey) && event.key === 'z' && !event.shiftKey) {
		event.preventDefault()
		builder.undo()
	}
	// Ctrl/Cmd + Y or Ctrl/Cmd + Shift + Z: Redo
	else if ((event.ctrlKey || event.metaKey) && (event.key === 'y' || (event.key === 'z' && event.shiftKey))) {
		event.preventDefault()
		builder.redo()
	}
	// Ctrl/Cmd + C: Copy
	else if ((event.ctrlKey || event.metaKey) && event.key === 'c') {
		event.preventDefault()
		builder.copyToClipboard()
	}
	// Ctrl/Cmd + X: Cut
	else if ((event.ctrlKey || event.metaKey) && event.key === 'x') {
		event.preventDefault()
		builder.cutToClipboard()
	}
	// Ctrl/Cmd + V: Paste
	else if ((event.ctrlKey || event.metaKey) && event.key === 'v') {
		event.preventDefault()
		builder.pasteFromClipboard()
	}
	// Ctrl/Cmd + D: Duplicate
	else if ((event.ctrlKey || event.metaKey) && event.key === 'd') {
		event.preventDefault()
		if (builder.canvasState.selectedLayerIds.length > 0) {
			builder.duplicateLayers(builder.canvasState.selectedLayerIds)
		}
	}
	// Delete/Backspace: Delete selected layers
	else if (event.key === 'Delete' || event.key === 'Backspace') {
		if (builder.canvasState.selectedLayerIds.length > 0) {
			event.preventDefault()
			builder.removeLayers(builder.canvasState.selectedLayerIds)
		}
	}
	// Ctrl/Cmd + A: Select all
	else if ((event.ctrlKey || event.metaKey) && event.key === 'a') {
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
const canAlign = computed(() => builder.canvasState.selectedLayerIds.length >= 2)
const canDistribute = computed(() => builder.canvasState.selectedLayerIds.length >= 3)
const multipleSelected = computed(() => builder.canvasState.selectedLayerIds.length > 1)
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
				<div class="border-b bg-background">
					<OverlaySettings
						v-model:overlay-name="overlayName"
						:overlay-id="initialProject?.id"
						@save="handleSave"
						@add-layer="handleAddLayer"
					/>
				</div>

				<!-- Layers Panel -->
				<div class="flex-1 min-h-0 overflow-hidde p-2">
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
					/>
				</div>

				<!-- Properties Panel -->
				<div class="flex-1 min-h-0 overflow-hidden p-2">
					<PropertiesPanel
						:layer="builder.activeLayer.value ?? null"
						:multiple-selected="multipleSelected"
						@update="handleUpdateLayerProperties"
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
