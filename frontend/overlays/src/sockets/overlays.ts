import { useWebSocket } from '@vueuse/core';
import { ref, watch } from 'vue';

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
	periodically_refetch_data: boolean;

	htmlContent?: string;
}

export interface LayerSettings {
	htmlOverlayDataPollSecondsInterval: number
	htmlOverlayHtml: string
	htmlOverlayCss: string
	htmlOverlayJs: string
}

export const useOverlays = (apiKey: string, overlayId: string) => {
	const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
	const host = window.location.host;

	const layers = ref<Array<Layer>>([]);

	const { data, send } = useWebSocket(
		`${protocol}://${host}/socket/registry/overlays?apiKey=${apiKey}`,
		{
			immediate: true,
			autoReconnect: {
				delay: 500,
			},
			onConnected() {
				send(JSON.stringify({
					eventName: 'getLayers',
					data: {
						overlayId,
					},
				}));
			},
		},
	);

	const parsedLayersData = ref<Record<string, string>>({});

	watch(data, (d) => {
		const parsedData = JSON.parse(d);

		if (parsedData.eventName === 'layers') {
			const parsedLayers = parsedData.layers as Array<Layer>;

			layers.value = parsedLayers.map(l => ({
				...l,
				settings: {
					...l.settings,
					htmlOverlayCss: b64DecodeUnicode(l.settings.htmlOverlayCss),
				},
			}));
		}

		if (parsedData.eventName === 'parsedLayerVariables') {
			parsedLayersData.value[parsedData.layerId] = b64DecodeUnicode(parsedData.data);
		}

		if (parsedData.eventName === 'refreshOverlays') {
			window.location.reload();
		}
	});

	const requestLayerData = (layerId: string) => {
		send(JSON.stringify({
			eventName: 'parseLayerVariables',
			data: {
				layerId,
			},
		}));
	};

	return {
		layers,
		parsedLayersData,
		requestLayerData,
	};
};

function b64DecodeUnicode(str: string) {
	return decodeURIComponent(
		atob(str)
			.split('')
			.map(function (c) {
				return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
			})
			.join(''),
	);
}
