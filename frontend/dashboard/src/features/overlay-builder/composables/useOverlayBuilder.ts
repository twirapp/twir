import { nanoid } from 'nanoid'
import { computed, reactive, ref, toRaw } from 'vue'

import { ChannelOverlayLayerType } from '@/gql/graphql'

import type { AlignmentGuide, CanvasState, HistoryState, Layer, OverlayProject } from '../types'

const MAX_HISTORY_SIZE = 50

export function useOverlayBuilder() {
	// Project state (canvas size fixed at 1920x1080)
	const project = reactive<OverlayProject>({
		id: '',
		name: '',
		width: 1920,
		height: 1080,
		layers: [],
	})

	// Canvas state
	const canvasState = reactive<CanvasState>({
		zoom: 1,
		panX: 0,
		panY: 0,
		selectedLayerIds: [],
		clipboardLayers: [],
		showGrid: true,
		snapToGrid: true,
		gridSize: 10,
		showRulers: false,
		showGuides: true,
	})

	// History management
	const history = reactive<HistoryState>({
		past: [],
		present: JSON.parse(JSON.stringify(toRaw(project))),
		future: [],
	})

	const canUndo = computed(() => history.past.length > 0)
	const canRedo = computed(() => history.future.length > 0)

	// Alignment guides
	const alignmentGuides = ref<AlignmentGuide[]>([])

	// Selected layers
	const selectedLayers = computed(() => {
		return project.layers.filter((layer) => canvasState.selectedLayerIds.includes(layer.id))
	})

	const activeLayer = computed(() => {
		if (canvasState.selectedLayerIds.length === 1) {
			return project.layers.find((layer) => layer.id === canvasState.selectedLayerIds[0])
		}
		return null
	})

	// Save current state to history
	function saveToHistory() {
		history.past.push(JSON.parse(JSON.stringify(toRaw(history.present))))
		if (history.past.length > MAX_HISTORY_SIZE) {
			history.past.shift()
		}
		history.present = JSON.parse(JSON.stringify(toRaw(project)))
		history.future = []
	}

	// Undo/Redo
	function undo() {
		if (!canUndo.value) return

		history.future.unshift(JSON.parse(JSON.stringify(toRaw(history.present))))
		const previous = history.past.pop()!
		history.present = previous

		Object.assign(project, JSON.parse(JSON.stringify(previous)))
	}

	function redo() {
		if (!canRedo.value) return

		history.past.push(JSON.parse(JSON.stringify(toRaw(history.present))))
		const next = history.future.shift()!
		history.present = next

		Object.assign(project, JSON.parse(JSON.stringify(next)))
	}

	// Layer operations
	function addLayer(type: ChannelOverlayLayerType, options?: Partial<Layer>) {
		const newLayer: Layer = {
			id: nanoid(),
			type,
			name: `${type} Layer ${project.layers.length + 1}`,
			posX: options?.posX ?? 100,
			posY: options?.posY ?? 100,
			width: options?.width ?? 200,
			height: options?.height ?? 200,
			rotation: 0,
			opacity: 1,
			visible: true,
			locked: false,
			zIndex: project.layers.length,
			periodicallyRefetchData: options?.periodicallyRefetchData ?? true,
			settings: options?.settings ?? {
				htmlOverlayHtml: '<span class="text">$(stream.uptime)</span>',
				htmlOverlayCss: '.text { color: #fff; font-size: 24px; }',
				htmlOverlayJs: 'function onDataUpdate() { console.log("updated") }',
				htmlOverlayDataPollSecondsInterval: 5,
			},
		}

		saveToHistory()
		project.layers.push(newLayer)
		selectLayers([newLayer.id])
	}

	function removeLayer(layerId: string) {
		saveToHistory()
		project.layers = project.layers.filter((layer) => layer.id !== layerId)
		canvasState.selectedLayerIds = canvasState.selectedLayerIds.filter((id) => id !== layerId)
		reorderLayers()
	}

	function removeLayers(layerIds: string[]) {
		saveToHistory()
		project.layers = project.layers.filter((layer) => !layerIds.includes(layer.id))
		canvasState.selectedLayerIds = canvasState.selectedLayerIds.filter(
			(id) => !layerIds.includes(id)
		)
		reorderLayers()
	}

	function updateLayer(layerId: string, updates: Partial<Layer>) {
		const layer = project.layers.find((l) => l.id === layerId)
		if (!layer) return

		saveToHistory()
		Object.assign(layer, updates)
	}

	function updateLayers(layerIds: string[], updates: Partial<Layer>) {
		saveToHistory()
		layerIds.forEach((id) => {
			const layer = project.layers.find((l) => l.id === id)
			if (layer) {
				Object.assign(layer, updates)
			}
		})
	}

	function duplicateLayer(layerId: string) {
		const layer = project.layers.find((l) => l.id === layerId)
		if (!layer) return

		saveToHistory()
		const duplicated: Layer = {
			...JSON.parse(JSON.stringify(toRaw(layer))),
			id: nanoid(),
			name: `${layer.name} (Copy)`,
			posX: layer.posX + 20,
			posY: layer.posY + 20,
			zIndex: project.layers.length,
		}

		project.layers.push(duplicated)
		selectLayers([duplicated.id])
	}

	function duplicateLayers(layerIds: string[]) {
		saveToHistory()
		const newIds: string[] = []

		layerIds.forEach((id) => {
			const layer = project.layers.find((l) => l.id === id)
			if (!layer) return

			const duplicated: Layer = {
				...JSON.parse(JSON.stringify(toRaw(layer))),
				id: nanoid(),
				name: `${layer.name} (Copy)`,
				posX: layer.posX + 20,
				posY: layer.posY + 20,
				zIndex: project.layers.length + newIds.length,
			}

			project.layers.push(duplicated)
			newIds.push(duplicated.id)
		})

		selectLayers(newIds)
	}

	// Layer ordering
	function moveLayerUp(layerId: string) {
		const index = project.layers.findIndex((l) => l.id === layerId)
		if (index === project.layers.length - 1) return

		saveToHistory()
		const temp = project.layers[index + 1]
		project.layers[index + 1] = project.layers[index]
		project.layers[index] = temp
		reorderLayers()
	}

	function moveLayerDown(layerId: string) {
		const index = project.layers.findIndex((l) => l.id === layerId)
		if (index === 0) return

		saveToHistory()
		const temp = project.layers[index - 1]
		project.layers[index - 1] = project.layers[index]
		project.layers[index] = temp
		reorderLayers()
	}

	function moveLayerToTop(layerId: string) {
		const layer = project.layers.find((l) => l.id === layerId)
		if (!layer) return

		saveToHistory()
		project.layers = project.layers.filter((l) => l.id !== layerId)
		project.layers.push(layer)
		reorderLayers()
	}

	function moveLayerToBottom(layerId: string) {
		const layer = project.layers.find((l) => l.id === layerId)
		if (!layer) return

		saveToHistory()
		project.layers = project.layers.filter((l) => l.id !== layerId)
		project.layers.unshift(layer)
		reorderLayers()
	}

	function reorderLayers(newOrder?: Layer[]) {
		if (newOrder) {
			saveToHistory()
			project.layers = newOrder
		}
		project.layers.forEach((layer, index) => {
			layer.zIndex = index
		})
	}

	// Selection
	function selectLayers(layerIds: string[], addToSelection = false) {
		if (addToSelection) {
			canvasState.selectedLayerIds = [...new Set([...canvasState.selectedLayerIds, ...layerIds])]
		} else {
			canvasState.selectedLayerIds = layerIds
		}
	}

	function deselectLayers(layerIds: string[]) {
		canvasState.selectedLayerIds = canvasState.selectedLayerIds.filter(
			(id) => !layerIds.includes(id)
		)
	}

	function selectAll() {
		canvasState.selectedLayerIds = project.layers
			.filter((l) => l.visible && !l.locked)
			.map((l) => l.id)
	}

	function deselectAll() {
		canvasState.selectedLayerIds = []
	}

	// Clipboard operations
	function copyToClipboard() {
		canvasState.clipboardLayers = selectedLayers.value.map((layer) =>
			JSON.parse(JSON.stringify(toRaw(layer)))
		)
	}

	function cutToClipboard() {
		copyToClipboard()
		removeLayers(canvasState.selectedLayerIds)
	}

	function pasteFromClipboard() {
		if (canvasState.clipboardLayers.length === 0) return

		saveToHistory()
		const newIds: string[] = []

		canvasState.clipboardLayers.forEach((layer) => {
			const pasted: Layer = {
				...JSON.parse(JSON.stringify(toRaw(layer))),
				id: nanoid(),
				name: `${layer.name} (Pasted)`,
				posX: layer.posX + 20,
				posY: layer.posY + 20,
				zIndex: project.layers.length + newIds.length,
			}

			project.layers.push(pasted)
			newIds.push(pasted.id)
		})

		selectLayers(newIds)
	}

	// Alignment
	function alignLayers(alignment: 'left' | 'center' | 'right' | 'top' | 'middle' | 'bottom') {
		if (selectedLayers.value.length < 2) return

		saveToHistory()
		const bounds = getSelectionBounds()

		selectedLayers.value.forEach((layer) => {
			switch (alignment) {
				case 'left':
					layer.posX = bounds.left
					break
				case 'center':
					layer.posX = bounds.left + (bounds.width - layer.width) / 2
					break
				case 'right':
					layer.posX = bounds.left + bounds.width - layer.width
					break
				case 'top':
					layer.posY = bounds.top
					break
				case 'middle':
					layer.posY = bounds.top + (bounds.height - layer.height) / 2
					break
				case 'bottom':
					layer.posY = bounds.top + bounds.height - layer.height
					break
			}
		})
	}

	function distributeLayersHorizontally() {
		if (selectedLayers.value.length < 3) return

		saveToHistory()
		const sorted = [...selectedLayers.value].sort((a, b) => a.posX - b.posX)
		const first = sorted[0]
		const last = sorted[sorted.length - 1]
		const totalWidth = last.posX + last.width - first.posX
		const totalLayerWidth = sorted.reduce((sum, layer) => sum + layer.width, 0)
		const spacing = (totalWidth - totalLayerWidth) / (sorted.length - 1)

		let currentX = first.posX + first.width
		for (let i = 1; i < sorted.length - 1; i++) {
			sorted[i].posX = currentX + spacing
			currentX = sorted[i].posX + sorted[i].width
		}
	}

	function distributeLayersVertically() {
		if (selectedLayers.value.length < 3) return

		saveToHistory()
		const sorted = [...selectedLayers.value].sort((a, b) => a.posY - b.posY)
		const first = sorted[0]
		const last = sorted[sorted.length - 1]
		const totalHeight = last.posY + last.height - first.posY
		const totalLayerHeight = sorted.reduce((sum, layer) => sum + layer.height, 0)
		const spacing = (totalHeight - totalLayerHeight) / (sorted.length - 1)

		let currentY = first.posY + first.height
		for (let i = 1; i < sorted.length - 1; i++) {
			sorted[i].posY = currentY + spacing
			currentY = sorted[i].posY + sorted[i].height
		}
	}

	function getSelectionBounds() {
		if (selectedLayers.value.length === 0) {
			return { left: 0, top: 0, width: 0, height: 0 }
		}

		const left = Math.min(...selectedLayers.value.map((l) => l.posX))
		const top = Math.min(...selectedLayers.value.map((l) => l.posY))
		const right = Math.max(...selectedLayers.value.map((l) => l.posX + l.width))
		const bottom = Math.max(...selectedLayers.value.map((l) => l.posY + l.height))

		return {
			left,
			top,
			width: right - left,
			height: bottom - top,
		}
	}

	// Canvas operations
	function setZoom(zoom: number) {
		canvasState.zoom = Math.max(0.1, Math.min(5, zoom))
	}

	function zoomIn() {
		setZoom(canvasState.zoom + 0.1)
	}

	function zoomOut() {
		setZoom(canvasState.zoom - 0.1)
	}

	function resetZoom() {
		canvasState.zoom = 1
	}

	function fitToScreen(containerWidth: number, containerHeight: number) {
		const scaleX = (containerWidth * 0.9) / project.width
		const scaleY = (containerHeight * 0.9) / project.height
		setZoom(Math.min(scaleX, scaleY))
	}

	// Snapping
	function snapToGrid(value: number): number {
		if (!canvasState.snapToGrid) return value
		return Math.round(value / canvasState.gridSize) * canvasState.gridSize
	}

	function findAlignmentGuides(layer: Layer): AlignmentGuide[] {
		if (!canvasState.showGuides) return []

		const guides: AlignmentGuide[] = []
		const threshold = 5

		project.layers.forEach((other) => {
			if (other.id === layer.id || !other.visible) return

			// Vertical alignment
			if (Math.abs(layer.posX - other.posX) < threshold) {
				guides.push({ type: 'vertical', position: other.posX, matchedLayers: [other.id] })
			}
			if (Math.abs(layer.posX + layer.width - (other.posX + other.width)) < threshold) {
				guides.push({
					type: 'vertical',
					position: other.posX + other.width,
					matchedLayers: [other.id],
				})
			}
			if (Math.abs(layer.posX + layer.width / 2 - (other.posX + other.width / 2)) < threshold) {
				guides.push({
					type: 'vertical',
					position: other.posX + other.width / 2,
					matchedLayers: [other.id],
				})
			}

			// Horizontal alignment
			if (Math.abs(layer.posY - other.posY) < threshold) {
				guides.push({ type: 'horizontal', position: other.posY, matchedLayers: [other.id] })
			}
			if (Math.abs(layer.posY + layer.height - (other.posY + other.height)) < threshold) {
				guides.push({
					type: 'horizontal',
					position: other.posY + other.height,
					matchedLayers: [other.id],
				})
			}
			if (Math.abs(layer.posY + layer.height / 2 - (other.posY + other.height / 2)) < threshold) {
				guides.push({
					type: 'horizontal',
					position: other.posY + other.height / 2,
					matchedLayers: [other.id],
				})
			}
		})

		return guides
	}

	// Constrain layers to canvas bounds
	function constrainLayersToCanvas() {
		project.layers.forEach((layer) => {
			// Ensure layer doesn't exceed canvas width
			if (layer.width > project.width) {
				layer.width = project.width
			}
			// Ensure layer doesn't exceed canvas height
			if (layer.height > project.height) {
				layer.height = project.height
			}

			// Constrain position to keep layer within canvas
			const maxX = project.width - layer.width
			const maxY = project.height - layer.height

			if (layer.posX > maxX) {
				layer.posX = Math.max(0, maxX)
			}
			if (layer.posY > maxY) {
				layer.posY = Math.max(0, maxY)
			}

			// Ensure position is not negative
			layer.posX = Math.max(0, layer.posX)
			layer.posY = Math.max(0, layer.posY)
		})
	}

	// Load project
	function loadProject(data: OverlayProject) {
		Object.assign(project, JSON.parse(JSON.stringify(data)))
		history.present = JSON.parse(JSON.stringify(toRaw(project)))
		history.past = []
		history.future = []
		deselectAll()
	}

	// Export project data
	function exportProject(): OverlayProject {
		const exported = JSON.parse(JSON.stringify(toRaw(project)))
		console.log(
			'[DEBUG] exportProject - layer rotations:',
			exported.layers.map((l: Layer) => ({ id: l.id, rotation: l.rotation }))
		)
		return exported
	}

	return {
		// State
		project,
		canvasState,
		selectedLayers,
		activeLayer,
		alignmentGuides,

		// History
		canUndo,
		canRedo,
		undo,
		redo,

		// Layer operations
		addLayer,
		removeLayer,
		removeLayers,
		updateLayer,
		updateLayers,
		duplicateLayer,
		duplicateLayers,

		// Layer ordering
		moveLayerUp,
		moveLayerDown,
		moveLayerToTop,
		moveLayerToBottom,
		reorderLayers,

		// Selection
		selectLayers,
		deselectLayers,
		selectAll,
		deselectAll,

		// Clipboard
		copyToClipboard,
		cutToClipboard,
		pasteFromClipboard,

		// Alignment
		alignLayers,
		distributeLayersHorizontally,
		distributeLayersVertically,
		getSelectionBounds,

		// Canvas
		setZoom,
		zoomIn,
		zoomOut,
		resetZoom,
		fitToScreen,

		// Snapping
		snapToGrid,
		findAlignmentGuides,

		// Project
		loadProject,
		exportProject,
		constrainLayersToCanvas,
	}
}
