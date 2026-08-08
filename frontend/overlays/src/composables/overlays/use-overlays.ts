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
	visible: boolean
	zIndex: number
}

export interface LayerSettings {
	htmlOverlayDataPollSecondsInterval: number
	htmlOverlayHtml: string
	htmlOverlayCss: string
	htmlOverlayJs: string
	imageUrl: string
	textContent: string
	textFontFamily: string
	textFontSize: number
	textFontWeight: number
	textColor: string
	textAlign: 'left' | 'center' | 'right'
	textFontStyle: string
	textAlignVertical: string
	textStrokeWidth: number
	textStrokeColor: string
	textShadowColor: string
	textShadowBlur: number
	textShadowOffsetX: number
	textShadowOffsetY: number
	textLineHeight: number
	textLetterSpacing: number
	textTransform: string
	videoUrl: string
	videoLoop: boolean
	videoMuted: boolean
	iframeUrl: string
	iframeScale: number
	widgetKey: string
	youtubeVideoId: string
	youtubeAutoplay: boolean
	youtubeLoop: boolean
	youtubeMuted: boolean
	emoteUrl: string
	emoteName: string
	emoteProvider: string
}

function normalizeTextAlign(value: string): LayerSettings['textAlign'] {
	switch (value) {
		case 'center':
			return 'center'
		case 'right':
			return 'right'
		default:
			return 'left'
	}
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
					instaSave
					layers {
						id
						type
						settings {
							htmlOverlayHtml
							htmlOverlayCss
							htmlOverlayJs
							htmlOverlayDataPollSecondsInterval
							imageUrl
							textContent
							textFontFamily
							textFontSize
							textFontWeight
							textColor
							textAlign
							textFontStyle
							textAlignVertical
							textStrokeWidth
							textStrokeColor
							textShadowColor
							textShadowBlur
							textShadowOffsetX
							textShadowOffsetY
							textLineHeight
							textLetterSpacing
							textTransform
							videoUrl
							videoLoop
							videoMuted
							iframeUrl
							iframeScale
							widgetKey
							youtubeVideoId
							youtubeAutoplay
							youtubeLoop
							youtubeMuted
							emoteUrl
							emoteName
							emoteProvider
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
						visible
						zIndex
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

	// Transform GraphQL data to the expected Layer format, filtering out hidden layers
	const layers = computed<Layer[]>(() => {
		if (!overlayData.value?.customOverlaySettings?.layers) {
			return []
		}

		return overlayData.value.customOverlaySettings.layers
			.filter((layer) => layer.visible)
			.map((layer) => ({
				id: layer.id,
				type: layer.type,
				settings: {
					htmlOverlayDataPollSecondsInterval: layer.settings.htmlOverlayDataPollSecondsInterval,
					htmlOverlayHtml: layer.settings.htmlOverlayHtml,
					htmlOverlayCss: layer.settings.htmlOverlayCss,
					htmlOverlayJs: layer.settings.htmlOverlayJs,
					imageUrl: layer.settings.imageUrl || '',
					textContent: layer.settings.textContent,
					textFontFamily: layer.settings.textFontFamily,
					textFontSize: layer.settings.textFontSize,
					textFontWeight: layer.settings.textFontWeight,
					textColor: layer.settings.textColor,
					textAlign: normalizeTextAlign(layer.settings.textAlign),
					textFontStyle: layer.settings.textFontStyle,
					textAlignVertical: layer.settings.textAlignVertical,
					textStrokeWidth: layer.settings.textStrokeWidth,
					textStrokeColor: layer.settings.textStrokeColor,
					textShadowColor: layer.settings.textShadowColor,
					textShadowBlur: layer.settings.textShadowBlur,
					textShadowOffsetX: layer.settings.textShadowOffsetX,
					textShadowOffsetY: layer.settings.textShadowOffsetY,
					textLineHeight: layer.settings.textLineHeight,
					textLetterSpacing: layer.settings.textLetterSpacing,
					textTransform: layer.settings.textTransform,
					videoUrl: layer.settings.videoUrl,
					videoLoop: layer.settings.videoLoop,
					videoMuted: layer.settings.videoMuted,
					iframeUrl: layer.settings.iframeUrl,
					iframeScale: layer.settings.iframeScale,
					widgetKey: layer.settings.widgetKey,
					youtubeVideoId: layer.settings.youtubeVideoId,
					youtubeAutoplay: layer.settings.youtubeAutoplay,
					youtubeLoop: layer.settings.youtubeLoop,
					youtubeMuted: layer.settings.youtubeMuted,
					emoteUrl: layer.settings.emoteUrl,
					emoteName: layer.settings.emoteName,
					emoteProvider: layer.settings.emoteProvider,
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
				visible: layer.visible,
				zIndex: layer.zIndex,
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
