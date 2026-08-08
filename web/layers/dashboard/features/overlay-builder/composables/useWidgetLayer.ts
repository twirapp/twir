import type { AcceptableValue } from 'reka-ui'
import { type Ref, computed, ref, watch } from 'vue'

import { useProfile } from '~~/layers/dashboard/api/auth.js'
import { useChatOverlayApi } from '~~/layers/dashboard/api/overlays/chat.js'
import { useDudesOverlayApi } from '~~/layers/dashboard/api/overlays/dudes.js'
import { useNowPlayingOverlayApi } from '~~/layers/dashboard/api/overlays/now-playing.js'
import { defaultDudesSettings } from '~~/layers/dashboard/pages/dashboard/overlays/dudes/dudes-settings.js'
import { defaultSettings as defaultNowPlayingSettings } from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/use-now-playing-form.js'

import { overlayWidgetRegistry } from '../widgets-registry'
import type { Layer, LayerSettings } from '../types'
import { useWidgetEntityId } from './useWidgetEntityId'
import {
	buildWidgetUrl as buildRegisteredWidgetUrl,
	getWidgetUrlParams,
	resolveWidgetLayerUrl,
} from './widget-url'

type UpdateLayer = (updates: Partial<Layer>) => void
type IframeSource = 'custom' | 'twir'

export function useWidgetLayer(layer: Ref<Layer>, updateLayer: UpdateLayer) {
	const { data: profile } = useProfile()
	const requestUrl = useRequestURL()

	const selectedDashboard = computed(() => {
		return profile.value?.availableDashboards.find(
			(dashboard) => dashboard.id === profile.value?.selectedDashboardId
		)
	})

	const overlayApiKey = computed(() => {
		return selectedDashboard.value?.channelApiKey || profile.value?.channelApiKey || ''
	})

	const selectedWidget = computed(() => {
		return overlayWidgetRegistry.find((widget) => widget.key === layer.value.settings.widgetKey)
	})

	const { data: chatOverlaysData, fetching: fetchingChatPresets } = useChatOverlayApi().useOverlaysQuery()
	const widgetUrlContext = computed(() => ({
		origin: requestUrl.origin,
		apiKey: overlayApiKey.value,
		chatPresetIds: chatOverlaysData.value?.chatOverlays.map((preset) => preset.id) ?? [],
		chatPresetsReady: !fetchingChatPresets.value && chatOverlaysData.value !== undefined,
	}))
	const firstChatPresetId = computed(() => {
		if (layer.value.settings.widgetKey !== 'chat') return undefined
		return chatOverlaysData.value?.chatOverlays[0]?.id ?? undefined
	})

	const dudesApi = useDudesOverlayApi()
	const dudesQuery = dudesApi.useDudesQuery(computed(() => layer.value.settings.widgetKey !== 'dudes'))
	const dudesCreator = dudesApi.useDudesCreate()

	useWidgetEntityId({
		layer,
		widgetKey: 'dudes',
		entities: computed(() => dudesQuery.data.value?.dudesGetAll),
		refetch: async () => {
			const result = await dudesQuery.executeQuery({ requestPolicy: 'network-only' })
			return result.data.value?.dudesGetAll ?? []
		},
		create: () => {
			const { id: _id, ...input } = defaultDudesSettings
			return dudesCreator.executeMutation({ input })
		},
		buildUrl: (id) => buildWidgetUrl({ id }),
		updateSettings,
	})

	const nowPlayingApi = useNowPlayingOverlayApi()
	const nowPlayingQuery = nowPlayingApi.useNowPlayingQuery(
		computed(() => layer.value.settings.widgetKey !== 'now-playing')
	)
	const nowPlayingCreator = nowPlayingApi.useNowPlayingCreate()

	useWidgetEntityId({
		layer,
		widgetKey: 'now-playing',
		entities: computed(() => nowPlayingQuery.data.value?.nowPlayingOverlays),
		refetch: async () => {
			const result = await nowPlayingQuery.executeQuery({ requestPolicy: 'network-only' })
			return result.data.value?.nowPlayingOverlays ?? []
		},
		create: () => {
			const { id: _id, channelId: _channelId, ...input } = defaultNowPlayingSettings
			return nowPlayingCreator.executeMutation({ input })
		},
		buildUrl: (id) => buildWidgetUrl({ id }),
		updateSettings,
	})

	function updateSettings(updates: Partial<LayerSettings>) {
		updateLayer({ settings: { ...layer.value.settings, ...updates } })
	}

	function buildWidgetUrl(params?: Record<string, string>) {
		const widget = selectedWidget.value
		if (!widget) return ''

		return buildRegisteredWidgetUrl(widget.key, widgetUrlContext.value, params)
	}

	watch(widgetUrlContext, (context) => {
		const iframeUrl = resolveWidgetLayerUrl(layer.value, context)
		if (iframeUrl !== layer.value.settings.iframeUrl) updateSettings({ iframeUrl })
	})

	function handleWidgetPresetSelect(presetId: string) {
		const iframeUrl = buildWidgetUrl({ id: presetId })
		if (iframeUrl !== layer.value.settings.iframeUrl) updateSettings({ iframeUrl })
	}

	function handleWidgetParamsUpdate(params: Record<string, string>) {
		const iframeUrl = buildWidgetUrl(params)
		if (iframeUrl !== layer.value.settings.iframeUrl) updateSettings({ iframeUrl })
	}

	const selectedWidgetSettings = computed(() => selectedWidget.value?.settingsComponent)
	const widgetParams = computed(() => getWidgetUrlParams(layer.value.settings.iframeUrl))
	const selectedWidgetParams = computed(() => {
		if (selectedWidget.value?.key !== 'faceit-stats' && selectedWidget.value?.key !== 'valorant-stats') return {}
		return { params: widgetParams.value }
	})
	const widgetSettingsOpen = ref(false)
	const iframeSource = ref<IframeSource>(layer.value.settings.widgetKey ? 'twir' : 'custom')

	watch(
		() => layer.value.settings.widgetKey,
		(widgetKey) => {
			iframeSource.value = widgetKey ? 'twir' : 'custom'
			if (!widgetKey) widgetSettingsOpen.value = false
		}
	)

	function handleIframeSourceChange(value: AcceptableValue) {
		if (typeof value !== 'string') return

		if (value === 'custom') {
			iframeSource.value = 'custom'
			updateSettings({ widgetKey: '' })
			return
		}

		if (value === 'twir') iframeSource.value = 'twir'
	}

	function handleWidgetChange(key: AcceptableValue) {
		if (typeof key !== 'string') return

		const widget = overlayWidgetRegistry.find((entry) => entry.key === key)
		if (!widget) return

		iframeSource.value = 'twir'
		updateSettings({
			widgetKey: widget.key,
			iframeUrl: widget.buildUrl({
				origin: requestUrl.origin,
				apiKey: overlayApiKey.value,
				params: widget.key === 'chat' && firstChatPresetId.value ? { id: firstChatPresetId.value } : undefined,
			}),
		})
	}

	const iframeUrl = computed({
		get: () => layer.value.settings.iframeUrl,
		set: (value: string) => {
			iframeSource.value = 'custom'
			updateSettings({ iframeUrl: value, widgetKey: '' })
		},
	})

	const iframeScale = computed({
		get: () => layer.value.settings.iframeScale,
		set: (value: string | number) => {
			const parsed = Number(value)
			if (Number.isFinite(parsed)) updateSettings({ iframeScale: Math.min(4, Math.max(0.1, parsed)) })
		},
	})

	const fieldId = (name: string) => `layer-${layer.value.id}-${name}`

	return {
		fieldId,
		overlayApiKey,
		selectedWidget,
		selectedWidgetSettings,
		selectedWidgetParams,
		widgetSettingsOpen,
		iframeSource,
		iframeUrl,
		iframeScale,
		handleIframeSourceChange,
		handleWidgetChange,
		handleWidgetPresetSelect,
		handleWidgetParamsUpdate,
		firstChatPresetId,
		buildWidgetUrl,
		overlayWidgetRegistry,
	}
}
