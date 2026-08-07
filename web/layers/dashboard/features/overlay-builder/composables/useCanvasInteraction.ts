import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import Moveable from 'vue3-moveable'
import type { OnDrag, OnResize, OnRotate } from 'vue3-moveable'

import type { AlignmentGuide, Layer } from '../types'

export interface CanvasInteractionProps {
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

export interface CanvasInteractionEmit {
	(event: 'updateLayer', layerId: string, updates: Partial<Layer>): void
	(event: 'selectLayer', layerId: string, addToSelection: boolean): void
	(event: 'deselectAll'): void
	(event: 'findGuides', layer: Layer): void
	(event: 'clearGuides'): void
}

export function useCanvasInteraction(
	props: CanvasInteractionProps,
	emit: CanvasInteractionEmit,
) {
	const canvasElement = ref<HTMLElement>()
	const moveableRef = ref<InstanceType<typeof Moveable>>()

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

	const canvasTransform = computed(() => `translate(${props.panX}px, ${props.panY}px)`)
	const canvasStyle = computed(() => {
		const baseStyle = {
			width: `${props.canvasWidth}px`,
			height: `${props.canvasHeight}px`,
			transform: canvasTransform.value,
			transformOrigin: 'center',
			zIndex: 1,
		}

		if (!props.showGrid) return baseStyle

		return {
			...baseStyle,
			backgroundImage: `
				linear-gradient(to right, rgba(255, 255, 255, 0.1) 1px, transparent 1px),
				linear-gradient(to bottom, rgba(255, 255, 255, 0.1) 1px, transparent 1px)
			`,
			backgroundSize: `${props.gridSize}px ${props.gridSize}px`,
		}
	})
	const canvasWrapperStyle = computed(() => ({
		transform: `scale(${props.zoom})`,
		transformOrigin: 'center',
	}))

	const selectedLayers = computed(() => props.layers.filter((layer) => props.selectedLayerIds.includes(layer.id)))
	const moveableTargets = computed(() => props.selectedLayerIds.map((id) => `#layer-${id}`))
	const isDragging = ref(false)

	let canvasMouseDownTarget: EventTarget | null = null

	function handleCanvasMouseDown(event: MouseEvent) {
		canvasMouseDownTarget = event.target
	}

	function handleCanvasClick(event: MouseEvent) {
		if (event.target === canvasElement.value && canvasMouseDownTarget === canvasElement.value) emit('deselectAll')
	}

	function snapToGuides(layer: Layer, posX: number, posY: number): { x: number; y: number } {
		if (!props.snapToGuidesEnabled) return { x: posX, y: posY }

		const snapThreshold = 8
		const releaseThreshold = 15
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
		const isSameLayer = snappedState.value.layerId === layer.id
		const wasSnappedX = isSameLayer && snappedState.value.snappedX
		const wasSnappedY = isSameLayer && snappedState.value.snappedY

		const checkVerticalSnap = (targetPos: number, snapTo: number) => {
			const distance = Math.abs(targetPos - snapTo)
			if (wasSnappedX && snappedState.value.snapPositionX === snapTo) return distance < releaseThreshold
			return distance < snapThreshold
		}
		const checkHorizontalSnap = (targetPos: number, snapTo: number) => {
			const distance = Math.abs(targetPos - snapTo)
			if (wasSnappedY && snappedState.value.snapPositionY === snapTo) return distance < releaseThreshold
			return distance < snapThreshold
		}

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

		props.layers.forEach((other) => {
			if (other.id === layer.id || !other.visible) return

			const otherCenterX = other.posX + other.width / 2
			const otherCenterY = other.posY + other.height / 2
			const otherRight = other.posX + other.width
			const otherBottom = other.posY + other.height

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
		const layer = props.layers.find((item) => item.id === layerId)
		if (layer?.locked) return

		const addToSelection = event.ctrlKey || event.metaKey
		if (!addToSelection && props.selectedLayerIds.length === 1 && props.selectedLayerIds[0] === layerId) return
		emit('selectLayer', layerId, addToSelection)
	}

	function handleLayerMouseDown(layerId: string, event: MouseEvent) {
		if (event.button !== 0) return
		const layer = props.layers.find((item) => item.id === layerId)
		if (!layer || layer.locked || props.selectedLayerIds.includes(layerId)) return

		emit('selectLayer', layerId, event.ctrlKey || event.metaKey)
		nextTick(() => moveableRef.value?.dragStart(event))
	}

	function onDrag(event: OnDrag) {
		isDragging.value = true
		if (!(event.target instanceof HTMLElement)) return
		const target = event.target
		const layerId = target.id.replace('layer-', '')
		const layer = props.layers.find((item) => item.id === layerId)
		if (!layer || layer.locked) return

		const maxX = props.canvasWidth - layer.width
		const maxY = props.canvasHeight - layer.height
		let newPosX = Math.max(0, Math.min(Math.round(event.translate[0] ?? 0), maxX))
		let newPosY = Math.max(0, Math.min(Math.round(event.translate[1] ?? 0), maxY))
		const snappedPos = snapToGuides(layer, newPosX, newPosY)
		newPosX = snappedPos.x
		newPosY = snappedPos.y

		target.style.transform = `translate(${newPosX}px, ${newPosY}px) rotate(${layer.rotation}deg)`
		emit('updateLayer', layerId, { posX: newPosX, posY: newPosY })
		emit('findGuides', { ...layer, posX: newPosX, posY: newPosY })
	}

	function onDragEnd() {
		isDragging.value = false
		emit('clearGuides')
		snappedState.value = {
			layerId: null,
			snappedX: false,
			snappedY: false,
			snapPositionX: null,
			snapPositionY: null,
		}
	}

	function onResize(event: OnResize) {
		if (!(event.target instanceof HTMLElement)) return
		const target = event.target
		const layerId = target.id.replace('layer-', '')
		const layer = props.layers.find((item) => item.id === layerId)
		if (!layer || layer.locked) return

		const minSize = 10
		const width = Math.min(props.canvasWidth, Math.max(minSize, Math.round(event.width)))
		const height = Math.min(props.canvasHeight, Math.max(minSize, Math.round(event.height)))
		const maxX = props.canvasWidth - width
		const maxY = props.canvasHeight - height
		const posX = Math.max(0, Math.min(Math.round(event.drag.translate[0] ?? 0), maxX))
		const posY = Math.max(0, Math.min(Math.round(event.drag.translate[1] ?? 0), maxY))

		target.style.width = `${width}px`
		target.style.height = `${height}px`
		target.style.transform = `translate(${posX}px, ${posY}px) rotate(${layer.rotation}deg)`
		emit('updateLayer', layerId, { width, height, posX, posY })
	}

	function onRotate(event: OnRotate) {
		if (!(event.target instanceof HTMLElement)) return
		const target = event.target
		const layerId = target.id.replace('layer-', '')
		const layer = props.layers.find((item) => item.id === layerId)
		if (!layer || layer.locked) return

		const rotation = Math.round(event.rotate)
		const maxX = props.canvasWidth - layer.width
		const maxY = props.canvasHeight - layer.height
		const posX = Math.max(0, Math.min(Math.round(event.drag.translate[0] ?? 0), maxX))
		const posY = Math.max(0, Math.min(Math.round(event.drag.translate[1] ?? 0), maxY))

		target.style.transform = `translate(${posX}px, ${posY}px) rotate(${rotation}deg)`
		emit('updateLayer', layerId, { rotation, posX, posY })
	}

	function getLayerStyle(layer: Layer) {
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

	const isLayerSelected = (layerId: string) => props.selectedLayerIds.includes(layerId)
	const quickActionsLayer = computed(() => selectedLayers.value.length === 1 ? selectedLayers.value[0] ?? null : null)
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
		return {
			left: `${Math.max(0, Math.min(layer.posX, props.canvasWidth - barWidth))}px`,
			top: `${quickActionsFlipAbove.value ? layer.posY : layer.posY + layer.height}px`,
			zIndex: 10000,
		}
	})
	const quickActionsInnerStyle = computed(() => ({
		transform: `scale(${1 / props.zoom})`,
		transformOrigin: '0 0',
	}))
	const quickActionsOffsetStyle = computed(() => ({
		transform: quickActionsFlipAbove.value ? 'translateY(calc(-100% - 8px))' : 'translateY(8px)',
	}))

	function handleKeyDown(event: KeyboardEvent) {
		if ((event.key === 'Delete' || event.key === 'Backspace') && props.selectedLayerIds.length > 0) event.preventDefault()
	}

	watch(
		() => props.layers.map((layer) => ({ id: layer.id, posX: layer.posX, posY: layer.posY, width: layer.width, height: layer.height, rotation: layer.rotation })),
		() => {
			if (!isDragging.value && moveableRef.value) nextTick(() => moveableRef.value?.updateRect())
		},
		{ deep: true }
	)

	onMounted(() => window.addEventListener('keydown', handleKeyDown))
	onUnmounted(() => window.removeEventListener('keydown', handleKeyDown))

	return {
		canvasElement,
		moveableRef,
		canvasTransform,
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
	}
}
