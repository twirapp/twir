import { type MaybeRefOrGetter, computed, onScopeDispose, ref, toValue, watch } from 'vue'

import type { Layer, OverlayProject } from '../types'

import { useOverlayInstantSaveGlobal } from './useOverlayInstantSave'

export interface OverlaySyncLayerPosition {
	id: string
	posX: number
	posY: number
	rotation: number
	width: number
	height: number
	visible: boolean
	opacity: number
}

export interface OverlaySyncSettingsUpdate {
	name?: string
	width?: number
	height?: number
	instaSave?: boolean
}

export interface OverlaySyncHandlers {
	onLayerAdd: (layer: Layer) => void
	onLayerRemove: (layerId: string) => void
	onLayerUpdate: (layer: Layer) => void
	onLayerPositions: (positions: OverlaySyncLayerPosition[]) => void
	onLayersReorder: (layerIds: string[]) => void
	onSettingsUpdate: (settings: OverlaySyncSettingsUpdate) => void
	onProjectReplace: (project: OverlayProject) => void
	getSyncState: () => OverlayProject | null
}

const POSITIONS_THROTTLE_MS = 60
const SYNC_STATE_TIMEOUT_MS = 2000

function isLayer(value: unknown): value is Layer {
	return (
		!!value &&
		typeof value === 'object' &&
		typeof (value as Layer).id === 'string' &&
		typeof (value as Layer).type === 'string'
	)
}

function isLayerPosition(value: unknown): value is OverlaySyncLayerPosition {
	return (
		!!value &&
		typeof value === 'object' &&
		typeof (value as OverlaySyncLayerPosition).id === 'string'
	)
}

function isProject(value: unknown): value is OverlayProject {
	return (
		!!value &&
		typeof value === 'object' &&
		Array.isArray((value as OverlayProject).layers)
	)
}

export function useOverlaySync(
	overlayId: MaybeRefOrGetter<string>,
	handlers: OverlaySyncHandlers
) {
	const { apiKey, status, open, close, sendEvent, addEditorEventHandler } = useOverlayInstantSaveGlobal()

	const currentOverlayId = computed(() => toValue(overlayId))
	const awaitingSyncState = ref(false)

	let syncStateTimer: ReturnType<typeof setTimeout> | null = null
	let positionsTrailingTimer: ReturnType<typeof setTimeout> | null = null
	let pendingPositions: OverlaySyncLayerPosition[] | null = null
	let lastPositionsSentAt = 0

	function send(eventName: string, data: Record<string, unknown>): boolean {
		const overlayIdValue = currentOverlayId.value
		if (!overlayIdValue) return false
		return sendEvent(eventName, { ...data, overlayId: overlayIdValue })
	}

	function sendLayerAdd(layer: Layer) {
		send('overlayEditorLayerAdd', { layer: JSON.parse(JSON.stringify(layer)) })
	}

	function sendLayerRemove(layerId: string) {
		send('overlayEditorLayerRemove', { layerId })
	}

	function sendLayerUpdate(layer: Layer) {
		send('overlayEditorLayerUpdate', { layer: JSON.parse(JSON.stringify(layer)) })
	}

	function sendLayersReorder(layerIds: string[]) {
		send('overlayEditorLayersReorder', { layerIds })
	}

	function sendSettingsUpdate(settings: OverlaySyncSettingsUpdate) {
		send('overlayEditorSettingsUpdate', settings)
	}

	function sendProjectReplace(project: OverlayProject) {
		send('overlayEditorProjectReplace', { project: JSON.parse(JSON.stringify(project)) })
	}

	function flushPositions() {
		lastPositionsSentAt = Date.now()
		positionsTrailingTimer = null
		if (!pendingPositions) return
		send('overlayEditorLayerPositions', { layers: pendingPositions })
		pendingPositions = null
	}

	function sendLayerPositions(positions: OverlaySyncLayerPosition[]) {
		pendingPositions = positions
		const elapsed = Date.now() - lastPositionsSentAt
		if (elapsed >= POSITIONS_THROTTLE_MS) {
			flushPositions()
			return
		}
		if (!positionsTrailingTimer) {
			positionsTrailingTimer = setTimeout(flushPositions, POSITIONS_THROTTLE_MS - elapsed)
		}
	}

	function sendSyncRequest() {
		if (!send('overlayEditorSyncRequest', {})) return

		awaitingSyncState.value = true
		if (syncStateTimer) clearTimeout(syncStateTimer)
		syncStateTimer = setTimeout(() => {
			awaitingSyncState.value = false
		}, SYNC_STATE_TIMEOUT_MS)
	}

	const removeHandler = addEditorEventHandler((eventName, data) => {
		if (data.overlayId !== currentOverlayId.value) return

		switch (eventName) {
			case 'overlayEditorLayerAdd':
				if (isLayer(data.layer)) handlers.onLayerAdd(data.layer)
				break
			case 'overlayEditorLayerUpdate':
				if (isLayer(data.layer)) handlers.onLayerUpdate(data.layer)
				break
			case 'overlayEditorLayerRemove':
				if (typeof data.layerId === 'string') handlers.onLayerRemove(data.layerId)
				break
			case 'overlayEditorLayerPositions':
				if (Array.isArray(data.layers)) {
					handlers.onLayerPositions(data.layers.filter(isLayerPosition))
				}
				break
			case 'overlayEditorLayersReorder':
				if (Array.isArray(data.layerIds)) {
					handlers.onLayersReorder(data.layerIds.filter((id): id is string => typeof id === 'string'))
				}
				break
			case 'overlayEditorSettingsUpdate':
				handlers.onSettingsUpdate({
					name: typeof data.name === 'string' ? data.name : undefined,
					width: typeof data.width === 'number' ? data.width : undefined,
					height: typeof data.height === 'number' ? data.height : undefined,
					instaSave: typeof data.instaSave === 'boolean' ? data.instaSave : undefined,
				})
				break
			case 'overlayEditorProjectReplace':
				if (isProject(data.project)) handlers.onProjectReplace(data.project)
				break
			case 'overlayEditorSyncRequest': {
				const state = handlers.getSyncState()
				if (state) send('overlayEditorSyncState', { project: JSON.parse(JSON.stringify(state)) })
				break
			}
			case 'overlayEditorSyncState':
				if (!awaitingSyncState.value || !isProject(data.project)) break
				awaitingSyncState.value = false
				if (syncStateTimer) clearTimeout(syncStateTimer)
				handlers.onProjectReplace(data.project)
				break
		}
	})

	watch(
		[status, currentOverlayId, apiKey],
		([socketStatus, id, key]) => {
			if (!id || !key) return
			if (socketStatus === 'CLOSED') open()
			if (socketStatus === 'OPEN') sendSyncRequest()
		},
		{ immediate: true }
	)

	onScopeDispose(() => {
		removeHandler()
		if (syncStateTimer) clearTimeout(syncStateTimer)
		if (positionsTrailingTimer) clearTimeout(positionsTrailingTimer)
	})

	return {
		status,
		close,
		sendLayerAdd,
		sendLayerRemove,
		sendLayerUpdate,
		sendLayerPositions,
		sendLayersReorder,
		sendSettingsUpdate,
		sendProjectReplace,
		sendSyncRequest,
	}
}
