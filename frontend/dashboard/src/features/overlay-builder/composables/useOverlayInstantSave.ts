import { createGlobalState, useWebSocket } from '@vueuse/core'
import { type MaybeRefOrGetter, computed, ref, toValue, watch } from 'vue'

import { useProfile } from '@/api/auth'

import type { OverlayProject } from '../types'

interface LayerPosition {
	id: string
	posX: number
	posY: number
	rotation: number
	width: number
	height: number
}

interface InstantSaveMessage {
	eventName: string
	data: {
		overlayId: string
		layers: LayerPosition[]
	}
}

export const useOverlayInstantSaveGlobal = createGlobalState(() => {
	const { data: profile } = useProfile()

	const selectedDashboard = computed(() => {
		return profile.value?.availableDashboards.find(
			(dashboard) => dashboard.id === profile.value?.selectedDashboardId
		)
	})

	const apiKey = computed(() => selectedDashboard.value?.apiKey ?? '')

	return {
		apiKey,
	}
})

export function useOverlayInstantSave(overlayId: MaybeRefOrGetter<string>) {
	const { apiKey } = useOverlayInstantSaveGlobal()

	const currentOverlayId = computed(() => toValue(overlayId))
	const isEnabled = ref(false)

	// Build WebSocket URL
	const wsUrl = computed(() => {
		if (!currentOverlayId.value || !apiKey.value) return null

		const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
		return `${wsProtocol}//${window.location.host}/socket/overlays/registry/overlays?apiKey=${apiKey.value}`
	})

	// Use VueUse WebSocket with auto-reconnect
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
		immediate: false, // Don't connect immediately, wait for overlay ID
	})

	// Watch for WebSocket messages
	watch(wsData, (message) => {
		if (!message) return

		try {
			const parsed = JSON.parse(message)

			if (parsed.eventName === 'instantSaveAck') {
				// Acknowledgment received
			}
		} catch (error) {
			console.error('[InstantSave] Failed to parse WebSocket message:', error)
		}
	})

	// Auto-connect when overlay ID is available
	watch(
		[currentOverlayId, apiKey],
		([id, key]) => {
			if (id && key && status.value === 'CLOSED') {
				open()
			}
		},
		{ immediate: true }
	)

	// Send layer positions via WebSocket
	function sendLayerPositions(project: OverlayProject) {
		const overlayIdValue = currentOverlayId.value

		if (!overlayIdValue) {
			console.warn('[InstantSave] No overlay ID')
			return false
		}

		if (status.value !== 'OPEN') {
			console.warn('[InstantSave] WebSocket not connected, status:', status.value)
			return false
		}

		const layersData: LayerPosition[] = project.layers.map((layer) => {
			return {
				id: layer.id,
				posX: layer.posX,
				posY: layer.posY,
				width: layer.width,
				height: layer.height,
				rotation: layer.rotation ?? 0,
			}
		})

		const message: InstantSaveMessage = {
			eventName: 'instantSaveLayerPositions',
			data: {
				overlayId: overlayIdValue,
				layers: layersData,
			},
		}

		try {
			send(JSON.stringify(message))
			return true
		} catch (error) {
			console.error('[InstantSave] Failed to send message:', error)
			return false
		}
	}

	return {
		status,
		isEnabled,
		sendLayerPositions,
		close,
	}
}
