import { createGlobalState, useWebSocket } from '@vueuse/core'
import { type MaybeRefOrGetter, computed, ref, toValue, watch } from 'vue'

import { useProfile } from '~~/layers/dashboard/api/auth'

import type { OverlayProject } from '../types'

interface LayerPosition {
	id: string
	posX: number
	posY: number
	rotation: number
	width: number
	height: number
	visible: boolean
	opacity: number
	zIndex: number
}

export type OverlayEditorEventHandler = (
	eventName: string,
	data: Record<string, unknown>,
) => void

export const useOverlayInstantSaveGlobal = createGlobalState(() => {
	const { data: profile } = useProfile()
	const requestUrl = useRequestURL()

	const selectedDashboard = computed(() => {
		return profile.value?.availableDashboards.find(
			(dashboard) => dashboard.id === profile.value?.selectedDashboardId
		)
	})

	const apiKey = computed(() => selectedDashboard.value?.channelApiKey || profile.value?.channelApiKey || '')

	// Random per-tab id; echoed back in broadcasts so we can ignore our own events.
	const clientId = crypto.randomUUID()

	const wsUrl = computed(() => {
		if (!apiKey.value) return null

		const wsProtocol = requestUrl.protocol === 'https:' ? 'wss:' : 'ws:'
		return `${wsProtocol}//${requestUrl.host}/socket/overlays/registry/overlays?apiKey=${apiKey.value}`
	})

	const {
		status,
		data: wsData,
		send,
		open,
		close,
	} = useWebSocket(wsUrl as any, {
		autoReconnect: {
			retries: 3,
			delay: 1000,
			onFailed() {
				console.error('[InstantSave] Failed to connect after retries')
			},
		},
		heartbeat: {
			message: JSON.stringify({ eventName: 'ping' }),
			interval: 30000,
		},
		immediate: false,
	})

	const editorEventHandlers = new Set<OverlayEditorEventHandler>()

	watch(wsData, (message) => {
		if (!message) return

		try {
			const parsed = JSON.parse(message)
			if (typeof parsed?.eventName !== 'string') return
			if (!parsed.eventName.startsWith('overlayEditor')) return

			const data = parsed.data
			if (!data || typeof data !== 'object') return
			if (data.clientId === clientId) return

			for (const handler of editorEventHandlers) {
				handler(parsed.eventName, data as Record<string, unknown>)
			}
		} catch (error) {
			console.error('[InstantSave] Failed to parse WebSocket message:', error)
		}
	})

	function addEditorEventHandler(handler: OverlayEditorEventHandler) {
		editorEventHandlers.add(handler)
		return () => {
			editorEventHandlers.delete(handler)
		}
	}

	function sendEvent(eventName: string, data: Record<string, unknown>): boolean {
		if (status.value !== 'OPEN') return false

		try {
			send(JSON.stringify({ eventName, data: { ...data, clientId } }))
			return true
		} catch (error) {
			console.error('[InstantSave] Failed to send message:', error)
			return false
		}
	}

	return {
		apiKey,
		clientId,
		status,
		open,
		close,
		sendEvent,
		addEditorEventHandler,
	}
})

export function useOverlayInstantSave(overlayId: MaybeRefOrGetter<string>) {
	const { apiKey, status, open, close, sendEvent } = useOverlayInstantSaveGlobal()

	const currentOverlayId = computed(() => toValue(overlayId))
	const isEnabled = ref(false)

	watch(
		[currentOverlayId, apiKey],
		([id, key]) => {
			if (id && key && status.value === 'CLOSED') {
				open()
			}
		},
		{ immediate: true }
	)

	function sendLayerPositions(project: OverlayProject) {
		const overlayIdValue = currentOverlayId.value

		if (!overlayIdValue) {
			console.warn('[InstantSave] No overlay ID')
			return false
		}

		const layersData: LayerPosition[] = project.layers.map((layer, index) => {
			return {
				id: layer.id,
				posX: layer.posX,
				posY: layer.posY,
				width: layer.width,
				height: layer.height,
				rotation: layer.rotation ?? 0,
				visible: layer.visible ?? true,
				opacity: layer.opacity ?? 1.0,
				zIndex: index,
			}
		})

		return sendEvent('instantSaveLayerPositions', {
			overlayId: overlayIdValue,
			layers: layersData,
		})
	}

	return {
		status,
		isEnabled,
		sendLayerPositions,
		close,
	}
}
