import { useSubscription } from '@urql/vue'
import { createGlobalState, useWebSocket } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

import type { ChannelOverlayLayerType } from '@/gql/graphql'

import { graphql } from '@/gql'
import { base64DecodeUnicode, generateSocketUrlWithParams } from '@/helpers.js'

export interface Layer {
	id: string
	type: ChannelOverlayLayerType
	settings: LayerSettings
	overlayId: string
	posX: number
	posY: number
	width: number
	height: number
	rotation: number
	createdAt: string
	updatedAt: string
	periodicallyRefetchData: boolean
}

export interface LayerSettings {
	htmlOverlayDataPollSecondsInterval: number
	htmlOverlayHtml: string
	htmlOverlayCss: string
	htmlOverlayJs: string
	imageUrl: string
}

export const useOverlays = createGlobalState(() => {
	const overlayUrl = ref('')
	const overlayId = ref('')
	const apiKey = ref('')

	const pauseGqlSub = computed(() => {
		return !overlayId.value || !apiKey.value
	})

	// Use GraphQL subscription to get real-time overlay updates
	const { data: overlayData } = useSubscription({
		query: graphql(`
			subscription CustomOverlaySettings($id: UUID!, $apiKey: String!) {
				customOverlaySettings(id: $id, apiKey: $apiKey) {
					id
					channelId
					name
					createdAt
					updatedAt
					width
					height
					layers {
						id
						type
						settings {
							htmlOverlayHtml
							htmlOverlayCss
							htmlOverlayJs
							htmlOverlayDataPollSecondsInterval
							imageUrl
						}
						overlayId
						posX
						posY
						width
						height
						rotation
						createdAt
						updatedAt
						periodicallyRefetchData
					}
				}
			}
		`),
		pause: pauseGqlSub,
		get variables() {
			return {
				id: overlayId.value,
				apiKey: apiKey.value,
			}
		},
		context: {},
	})

	// Transform GraphQL data to the expected Layer format
	const layers = computed<Layer[]>(() => {
		if (!overlayData.value?.customOverlaySettings?.layers) {
			return []
		}

		return overlayData.value.customOverlaySettings.layers.map((layer) => ({
			id: layer.id,
			type: layer.type,
			settings: {
				htmlOverlayDataPollSecondsInterval: layer.settings.htmlOverlayDataPollSecondsInterval,
				htmlOverlayHtml: layer.settings.htmlOverlayHtml,
				htmlOverlayCss: layer.settings.htmlOverlayCss,
				htmlOverlayJs: layer.settings.htmlOverlayJs,
				imageUrl: layer.settings.imageUrl || '',
			},
			overlayId: layer.overlayId,
			posX: layer.posX,
			posY: layer.posY,
			width: layer.width,
			height: layer.height,
			rotation: layer.rotation || 0,
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
			// Layers are now fetched via GraphQL subscription
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
