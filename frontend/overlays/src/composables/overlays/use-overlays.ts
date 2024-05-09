import { createGlobalState, useWebSocket } from '@vueuse/core'
import { ref, watch } from 'vue'

import { base64DecodeUnicode, generateSocketUrlWithParams } from '@/helpers.js'

export interface Layer {
	id: string
	type: 'HTML'
	settings: LayerSettings
	overlay_id: string
	pos_x: number
	pos_y: number
	width: number
	height: number
	createdAt: string
	updatedAt: string
	overlay: any
	periodically_refetch_data: boolean
	htmlContent?: string
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
	const layers = ref<Array<Layer>>([])

	const { data, status, send, open } = useWebSocket(
		overlayUrl,
		{
			immediate: false,
			autoReconnect: {
				delay: 500,
			},
			onConnected() {
				send(
					JSON.stringify({
						eventName: 'getLayers',
						data: {
							overlayId: overlayId.value,
						},
					}),
				)
			},
		},
	)

	const parsedLayersData = ref<Record<string, string>>({})

	watch(data, (d) => {
		const parsedData = JSON.parse(d)

		if (parsedData.eventName === 'layers') {
			const parsedLayers = parsedData.layers as Array<Layer>

			layers.value = parsedLayers.map((l) => ({
				...l,
				settings: {
					...l.settings,
					htmlOverlayCss: l.settings.htmlOverlayCss
						? base64DecodeUnicode(l.settings.htmlOverlayCss)
						: '',
					htmlOverlayJs: l.settings.htmlOverlayJs ? base64DecodeUnicode(l.settings.htmlOverlayJs) : '',
				},
			}))
		}

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
			}),
		)
	}

	function connectToOverlays(apiKey: string, _overlayId: string): void {
		if (status.value === 'OPEN') return

		const url = generateSocketUrlWithParams('/overlays/registry/overlays', {
			apiKey,
		})

		overlayUrl.value = url
		overlayId.value = _overlayId

		open()
	}

	return {
		layers,
		parsedLayersData,
		requestLayerData,
		connectToOverlays,
	}
})
