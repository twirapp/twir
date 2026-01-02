import { createGlobalState, useWebSocket } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

import type { ChannelOverlayLayerType } from '@/gql/graphql'

import { base64DecodeUnicode, generateSocketUrlWithParams } from '@/helpers.js'

import { useCustomOverlayById } from './use-custom-overlay.js'

export interface Layer {
	id: string
	type: ChannelOverlayLayerType
	settings: LayerSettings
	overlayId: string
	posX: number
	posY: number
	width: number
	height: number
	createdAt: string
	updatedAt: string
	periodicallyRefetchData: boolean
}

export interface LayerSettings {
	htmlOverlayDataPollSecondsInterval: number
	htmlOverlayHtml: string
	htmlOverlayCss: string
	htmlOverlayJs: string
}

export const useOverlays = createGlobalState(() => {
	const overlayUrl = ref('')
	const overlayId = ref('')
	const apiKey = ref('')

	// Use GraphQL query to fetch overlay data
	const { data: overlayData } = useCustomOverlayById(overlayId)

	// Transform GraphQL data to the expected Layer format
	const layers = computed<Layer[]>(() => {
		if (!overlayData.value?.channelOverlayById?.layers) {
			return []
		}

		return overlayData.value.channelOverlayById.layers.map((layer) => ({
			id: layer.id,
			type: layer.type,
			settings: {
				htmlOverlayDataPollSecondsInterval: layer.settings.htmlOverlayDataPollSecondsInterval,
				htmlOverlayHtml: layer.settings.htmlOverlayHtml,
				htmlOverlayCss: layer.settings.htmlOverlayCss,
				htmlOverlayJs: layer.settings.htmlOverlayJs,
			},
			overlayId: layer.overlayId,
			posX: layer.posX,
			posY: layer.posY,
			width: layer.width,
			height: layer.height,
			createdAt: layer.createdAt,
			updatedAt: layer.updatedAt,
			periodicallyRefetchData: layer.periodicallyRefetchData,
		}))
	})

	// Keep WebSocket for real-time variable parsing
	const { data, status, send, open } = useWebSocket(overlayUrl, {
		immediate: false,
		autoReconnect: {
			delay: 500,
		},
		onConnected() {
			// No longer need to fetch layers via WebSocket
			// Layers are now fetched via GraphQL
		},
	})

	const parsedLayersData = ref<Record<string, string>>({})

	watch(data, (d) => {
		if (!d) return

		const parsedData = JSON.parse(d)

		if (parsedData.eventName === 'parsedLayerVariables') {
			parsedLayersData.value[parsedData.layerId] = parsedData.data
				? base64DecodeUnicode(parsedData.data)
				: ''
		}

		if (parsedData.eventName === 'refreshOverlays') {
			window.location.reload()
		}
	})

	function requestLayerData(layerId: string): void {
		send(
			JSON.stringify({
				eventName: 'parseLayerVariables',
				data: {
					layerId,
				},
			})
		)
	}

	function connectToOverlays(_apiKey: string, _overlayId: string): void {
		const url = generateSocketUrlWithParams('/overlays/registry/overlays', {
			apiKey: _apiKey,
		})

		overlayUrl.value = url
		overlayId.value = _overlayId
		apiKey.value = _apiKey

		if (status.value !== 'OPEN') {
			open()
		}
	}

	return {
		layers,
		parsedLayersData,
		requestLayerData,
		connectToOverlays,
	}
})
