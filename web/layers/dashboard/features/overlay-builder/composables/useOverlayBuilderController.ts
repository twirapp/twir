import { type Ref, computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { toast } from 'vue-sonner'
import type { ChannelOverlayLayer, ChannelOverlayLayerType } from '~/gql/graphql.js'
import { useProfile } from '~~/layers/dashboard/api/auth.js'
import { useDudesOverlayApi } from '~~/layers/dashboard/api/overlays/dudes.js'
import { useNowPlayingOverlayApi } from '~~/layers/dashboard/api/overlays/now-playing.js'
import { defaultDudesSettings } from '~~/layers/dashboard/pages/dashboard/overlays/dudes/dudes-settings.js'
import { defaultSettings as defaultNowPlayingSettings } from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/use-now-playing-form.js'

import { useOverlayBuilder } from './useOverlayBuilder'
import { useChatOverlayPresetQuery } from './useChatOverlayPresets'
import { type OverlaySyncLayerPosition, type OverlaySyncSettingsUpdate, useOverlaySync } from './useOverlaySync'
import { createWidgetEntityResolver } from './useWidgetEntityId'
import { buildWidgetUrl as buildRegisteredWidgetUrl, getWidgetUrlParams, resolveWidgetLayerUrl } from './widget-url'
import type { Layer, OverlayProject } from '../types'
import { createLayerSettings } from '../types'
import { getLayerTypeMeta } from '../layer-type-meta'

export interface InitialProjectLayer {
	id: string
	type: ChannelOverlayLayer['type']
	name?: string
	posX: number
	posY: number
	width: number
	height: number
	rotation: number
	opacity?: number
	visible?: boolean
	locked?: boolean
	periodicallyRefetchData: boolean
	settings?: Partial<Layer['settings']>
}

export interface InitialOverlayProject {
	id: string
	name: string
	width: number
	height: number
	instaSave?: boolean
	layers: InitialProjectLayer[]
}

interface OverlayBuilderEmit {
	(event: 'save', project: OverlayProject): void
	(event: 'instantSave', project: OverlayProject): void
}

type CanvasRef = HTMLElement | { readonly $el?: Element }

export function useOverlayBuilderController(
	initialProject: Ref<InitialOverlayProject | undefined>,
	emit: OverlayBuilderEmit,
) {
	const builder = useOverlayBuilder()
	const { t } = useI18n()
	const overlayName = ref('')
	const instaSave = ref(false)
	const canvasAreaRef = ref<CanvasRef>()
	const loadedProjectId = ref('')
	const isLoadingProject = ref(false)
	const isApplyingRemote = ref(false)
	const addLayersHidden = ref(false)
	const showCodeEditor = ref(false)
	const showShortcuts = ref(false)
	const coveredLayerHintShown = ref(false)
	const editorLayer = ref<Layer | null>(null)
	const { data: profile } = useProfile()
	const requestUrl = useRequestURL()
	const { chatOverlaysData, fetchingOverlays: fetchingChatPresets } = useChatOverlayPresetQuery()
	const selectedDashboard = computed(() => profile.value?.availableDashboards.find(
		(dashboard) => dashboard.id === profile.value?.selectedDashboardId
	))
	const overlayApiKey = computed(() => selectedDashboard.value?.channelApiKey || profile.value?.channelApiKey || '')
	const chatPresetIds = computed(() => chatOverlaysData.value?.chatOverlays.map((preset) => preset.id) ?? [])
	const chatPresetsReady = computed(() => !fetchingChatPresets.value && chatOverlaysData.value !== undefined)

	function applyRemote(mutate: () => void) {
		isApplyingRemote.value = true
		try {
			mutate()
		} finally {
			void nextTick(() => {
				isApplyingRemote.value = false
			})
		}
	}

	function applyRemoteProject(project: OverlayProject) {
		applyRemote(() => {
			builder.applyRemoteProject(project)
			overlayName.value = project.name
			instaSave.value = project.instaSave
		})
	}

	const sync = useOverlaySync(loadedProjectId, {
		onLayerAdd: (layer) => applyRemote(() => builder.applyRemoteLayerAdd(layer)),
		onLayerRemove: (layerId) => applyRemote(() => builder.applyRemoteLayerRemove(layerId)),
		onLayerUpdate: (layer) => applyRemote(() => builder.applyRemoteLayerUpdate(layer)),
		onLayerPositions: (positions) => applyRemote(() => builder.applyRemoteLayerPositions(positions)),
		onLayersReorder: (layerIds) => applyRemote(() => builder.applyRemoteLayersReorder(layerIds)),
		onSettingsUpdate: (settings: OverlaySyncSettingsUpdate) =>
			applyRemote(() => {
				if (settings.name !== undefined) overlayName.value = settings.name
				if (settings.instaSave !== undefined) instaSave.value = settings.instaSave
				if (settings.width !== undefined) builder.project.width = settings.width
				if (settings.height !== undefined) builder.project.height = settings.height
			}),
		onProjectReplace: (project) => {
			if (isLoadingProject.value) {
				const stop = watch(isLoadingProject, (loading) => {
					if (loading) return
					stop()
					applyRemoteProject(project)
				})
				return
			}
			applyRemoteProject(project)
		},
		getSyncState: () => {
			if (!loadedProjectId.value || isLoadingProject.value) return null
			return projectSnapshot(overlayName.value, instaSave.value)
		},
	})

	function positionsOfLayers(layerIds: string[]): OverlaySyncLayerPosition[] {
		return builder.project.layers
			.filter((layer) => layerIds.includes(layer.id))
			.map((layer) => ({
				id: layer.id,
				posX: layer.posX,
				posY: layer.posY,
				rotation: layer.rotation ?? 0,
				width: layer.width,
				height: layer.height,
				visible: layer.visible ?? true,
				opacity: layer.opacity ?? 1.0,
			}))
	}

	const GEOMETRY_ONLY_KEYS = new Set(['posX', 'posY', 'width', 'height', 'rotation'])

	function broadcastLayerMutation(layerId: string, updates: Partial<Layer>) {
		const hasNonGeometryChange = Object.keys(updates).some((key) => !GEOMETRY_ONLY_KEYS.has(key))
		if (hasNonGeometryChange) {
			const layer = builder.project.layers.find((item) => item.id === layerId)
			if (layer) sync.sendLayerUpdate(layer)
			return
		}
		sync.sendLayerPositions(positionsOfLayers([layerId]))
	}

	function broadcastAddedLayers(layers: Layer[] | null | undefined) {
		layers?.forEach((layer) => sync.sendLayerAdd(layer))
	}

	function broadcastRemovedLayers(layerIds: string[]) {
		layerIds.forEach((layerId) => sync.sendLayerRemove(layerId))
	}

	function calculateFitZoom() {
		const canvasArea = canvasAreaRef.value instanceof HTMLElement ? canvasAreaRef.value : canvasAreaRef.value?.$el
		if (!canvasArea) return

		const availableWidth = canvasArea.clientWidth - 64
		const availableHeight = canvasArea.clientHeight - 64
		const scaleX = availableWidth / builder.project.width
		const scaleY = availableHeight / builder.project.height
		builder.setZoom(Math.min(scaleX, scaleY) * 0.8)
	}

	function normalizeWidgetUrls() {
		if (!overlayApiKey.value) return

		const context = {
			origin: requestUrl.origin,
			apiKey: overlayApiKey.value,
			chatPresetIds: chatPresetIds.value,
			chatPresetsReady: chatPresetsReady.value,
		}
		builder.project.layers.forEach((layer) => {
			if (!layer.settings.widgetKey) return
			const iframeUrl = resolveWidgetLayerUrl(layer, context)
			if (iframeUrl !== layer.settings.iframeUrl) {
				builder.updateLayer(layer.id, { settings: { ...layer.settings, iframeUrl } })
			}
		})
	}

	const hasDudesWidget = computed(() => builder.project.layers.some((layer) => layer.settings.widgetKey === 'dudes'))
	const hasNowPlayingWidget = computed(() => builder.project.layers.some((layer) => layer.settings.widgetKey === 'now-playing'))

	const dudesApi = useDudesOverlayApi()
	const dudesQuery = dudesApi.useDudesQuery(computed(() => !hasDudesWidget.value))
	const dudesCreator = dudesApi.useDudesCreate()
	const dudesResolver = createWidgetEntityResolver({
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
	})

	const nowPlayingApi = useNowPlayingOverlayApi()
	const nowPlayingQuery = nowPlayingApi.useNowPlayingQuery(computed(() => !hasNowPlayingWidget.value))
	const nowPlayingCreator = nowPlayingApi.useNowPlayingCreate()
	const nowPlayingResolver = createWidgetEntityResolver({
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
	})

	function healWidgetEntityIds() {
		if (!overlayApiKey.value) return

		const context = {
			origin: requestUrl.origin,
			apiKey: overlayApiKey.value,
			chatPresetIds: chatPresetIds.value,
			chatPresetsReady: chatPresetsReady.value,
		}

		builder.project.layers.forEach((layer) => {
			const widgetKey = layer.settings.widgetKey
			const resolver = widgetKey === 'dudes'
				? dudesResolver
				: widgetKey === 'now-playing'
					? nowPlayingResolver
					: null
			if (!resolver) return

			void (async () => {
				const currentId = getWidgetUrlParams(layer.settings.iframeUrl).id
				const id = await resolver.resolveId(currentId)
				if (!id || id === currentId) return

				const iframeUrl = buildRegisteredWidgetUrl(widgetKey, context, { id })
				if (iframeUrl && iframeUrl !== layer.settings.iframeUrl) {
					builder.updateLayer(layer.id, { settings: { ...layer.settings, iframeUrl } })
				}
			})()
		})
	}

	watch(
		[
			loadedProjectId,
			overlayApiKey,
			hasDudesWidget,
			hasNowPlayingWidget,
			() => dudesQuery.data.value?.dudesGetAll,
			() => nowPlayingQuery.data.value?.nowPlayingOverlays,
		],
		healWidgetEntityIds,
		{ immediate: true },
	)

	function loadInitialProject() {
		const project = initialProject.value
		if (!project || (loadedProjectId.value === project.id && project.id !== '')) return

		loadedProjectId.value = project.id
		overlayName.value = project.name || ''
		instaSave.value = project.instaSave || false
		const layers = project.layers.map((layer, index) => ({
			id: layer.id || `layer-${index}`,
			type: layer.type,
			name: layer.name || t('overlayBuilder.layerNames.fallback', { count: index + 1 }),
			posX: layer.posX,
			posY: layer.posY,
			width: layer.width,
			height: layer.height,
			rotation: Number(layer.rotation) || 0,
			opacity: layer.opacity || 1,
			visible: layer.visible !== undefined ? layer.visible : true,
			locked: layer.locked || false,
			zIndex: index,
			periodicallyRefetchData: layer.periodicallyRefetchData,
			settings: createLayerSettings(layer.settings, t('overlayBuilder.defaults.textContent')),
		}))

		isLoadingProject.value = true
		builder.loadProject({
			id: project.id,
			name: project.name,
			width: project.width,
			height: project.height,
			instaSave: project.instaSave || false,
			layers,
		})
		nextTick(() => { isLoadingProject.value = false; calculateFitZoom() })
	}

	watch(() => initialProject.value?.id, (newId) => {
		if (newId !== undefined) loadInitialProject()
	}, { immediate: true })

	watch(() => initialProject.value?.instaSave, (newInstaSave) => {
		if (loadedProjectId.value && newInstaSave !== undefined && newInstaSave !== instaSave.value) instaSave.value = newInstaSave
	})

	watch(
		[loadedProjectId, overlayApiKey, chatPresetsReady, () => chatPresetIds.value.join(',')],
		normalizeWidgetUrls,
		{ immediate: true },
	)

	onMounted(() => {
		loadInitialProject()
		window.addEventListener('resize', calculateFitZoom)
	})
	onUnmounted(() => window.removeEventListener('resize', calculateFitZoom))

	function addLayer(type: ChannelOverlayLayerType) {
		const layer = builder.addLayer(type, {
			name: t('overlayBuilder.layerNames.default', {
				type: t(getLayerTypeMeta(type).labelKey),
				count: builder.project.layers.length + 1,
			}),
			settings: { textContent: t('overlayBuilder.defaults.textContent') },
			visible: !addLayersHidden.value,
		})
		if (layer) sync.sendLayerAdd(layer)
	}

	function projectSnapshot(name: string, instantSave: boolean): OverlayProject {
		const project = builder.exportProject()
		project.name = name
		project.instaSave = instantSave
		return project
	}

	function handleSave() {
		emit('save', projectSnapshot(overlayName.value, instaSave.value))
	}

	async function handleLayerUpdate() {
		if (!instaSave.value) return
		await nextTick()
		emit('instantSave', projectSnapshot(overlayName.value, instaSave.value))
	}

	watch(instaSave, (newValue, oldValue) => {
		if (newValue === oldValue || isApplyingRemote.value || isLoadingProject.value) return
		if (loadedProjectId.value) emit('instantSave', projectSnapshot(overlayName.value, newValue))
		sync.sendSettingsUpdate({ instaSave: newValue })
	})
	watch([() => builder.project.width, () => builder.project.height], () => {
		if (isLoadingProject.value) return
		builder.constrainLayersToCanvas()
		void nextTick(calculateFitZoom)
		if (isApplyingRemote.value) return
		if (instaSave.value && loadedProjectId.value) emit('save', projectSnapshot(overlayName.value, instaSave.value))
		if (loadedProjectId.value) {
			sync.sendSettingsUpdate({ width: builder.project.width, height: builder.project.height })
		}
	})

	let nameSyncTimer: ReturnType<typeof setTimeout> | null = null
	watch(overlayName, (value) => {
		if (isApplyingRemote.value || isLoadingProject.value || !loadedProjectId.value) return
		if (nameSyncTimer) clearTimeout(nameSyncTimer)
		nameSyncTimer = setTimeout(() => sync.sendSettingsUpdate({ name: value }), 400)
	})

	function handleUpdateLayer(layerId: string, updates: Partial<Layer>) {
		builder.updateLayer(layerId, updates)
		broadcastLayerMutation(layerId, updates)
		if (updates.posX !== undefined || updates.posY !== undefined || updates.rotation !== undefined || updates.width !== undefined || updates.height !== undefined || updates.opacity !== undefined || updates.visible !== undefined) void handleLayerUpdate()
	}

	function handleSelectLayer(layerId: string, addToSelection: boolean) {
		builder.selectLayers([layerId], addToSelection)
	}

	function handleLayersPanelSelect(layerId: string, addToSelection: boolean) {
		builder.selectLayers([layerId], addToSelection)
		if (coveredLayerHintShown.value) return

		const layer = builder.project.layers.find((item) => item.id === layerId)
		if (!layer) return

		const isCovered = builder.project.layers.some((other) =>
			other.id !== layer.id &&
			other.visible &&
			!other.locked &&
			other.zIndex > layer.zIndex &&
			layer.posX < other.posX + other.width &&
			layer.posX + layer.width > other.posX &&
			layer.posY < other.posY + other.height &&
			layer.posY + layer.height > other.posY
		)
		if (!isCovered) return

		coveredLayerHintShown.value = true
		toast.info(t('overlayBuilder.canvas.coveredLayerHint'), { duration: 6000 })
	}

	function handleDeselectAll() {
		builder.deselectAll()
	}

	function handleFindGuides(layer: Layer) {
		builder.alignmentGuides.value = builder.findAlignmentGuides(layer)
	}

	function handleClearGuides() {
		builder.alignmentGuides.value = []
	}

	function handleToggleVisibility(layerId: string) {
		const layer = builder.project.layers.find((item) => item.id === layerId)
		if (!layer) return
		builder.updateLayer(layerId, { visible: !layer.visible })
		sync.sendLayerUpdate(layer)
		void handleLayerUpdate()
	}

	function handleToggleLock(layerId: string) {
		const layer = builder.project.layers.find((item) => item.id === layerId)
		if (!layer) return
		builder.updateLayer(layerId, { locked: !layer.locked })
		sync.sendLayerUpdate(layer)
	}

	function handleRemoveLayer(layerId: string) {
		builder.removeLayer(layerId)
		broadcastRemovedLayers([layerId])
	}

	function handleRemoveLayers(layerIds: string[]) {
		builder.removeLayers(layerIds)
		broadcastRemovedLayers(layerIds)
	}

	function handleDuplicateLayers(layerIds: string[]) {
		broadcastAddedLayers(builder.duplicateLayers(layerIds))
	}

	function handleCutSelection() {
		const layerIds = [...builder.canvasState.selectedLayerIds]
		builder.cutToClipboard()
		broadcastRemovedLayers(layerIds)
	}

	function handlePaste() {
		broadcastAddedLayers(builder.pasteFromClipboard())
	}

	function handleAlign(alignment: 'left' | 'center' | 'right' | 'top' | 'middle' | 'bottom') {
		builder.alignLayers(alignment)
		sync.sendLayerPositions(positionsOfLayers(builder.canvasState.selectedLayerIds))
	}

	function handleDistribute(direction: 'horizontal' | 'vertical') {
		if (direction === 'horizontal') {
			builder.distributeLayersHorizontally()
		} else {
			builder.distributeLayersVertically()
		}
		sync.sendLayerPositions(positionsOfLayers(builder.canvasState.selectedLayerIds))
	}

	function handleUndo() {
		builder.undo()
		sync.sendProjectReplace(projectSnapshot(overlayName.value, instaSave.value))
	}

	function handleRedo() {
		builder.redo()
		sync.sendProjectReplace(projectSnapshot(overlayName.value, instaSave.value))
	}

	function handleOpenLayerSettings(layerId: string) {
		builder.selectLayers([layerId])
		nextTick(() => {
			const card = document.getElementById('layer-properties-card')
			card?.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
			card?.focus({ preventScroll: true })
		})
	}

	function handleReorderLayers(layers: Layer[]) {
		builder.reorderLayers(layers)
		sync.sendLayersReorder(builder.project.layers.map((layer) => layer.id))
		void handleLayerUpdate()
	}

	function handleUpdateLayerProperties(layerId: string, updates: Partial<Layer>) {
		builder.updateLayer(layerId, updates)
		broadcastLayerMutation(layerId, updates)
		if (updates.posX !== undefined || updates.posY !== undefined || updates.rotation !== undefined || updates.width !== undefined || updates.height !== undefined || updates.opacity !== undefined || updates.visible !== undefined) void handleLayerUpdate()
	}

	function handleActiveLayerUpdate(updates: Partial<Layer>) {
		const layer = builder.activeLayer.value
		if (layer) handleUpdateLayerProperties(layer.id, updates)
	}

	function handleOpenCodeEditor() {
		if (builder.activeLayer.value) {
			editorLayer.value = builder.activeLayer.value
			showCodeEditor.value = true
		}
	}

	function handleSaveCode(data: { html: string; css: string; js: string; refreshInterval: number }) {
		if (!editorLayer.value) return
		builder.updateLayer(editorLayer.value.id, {
			settings: {
				...editorLayer.value.settings,
				htmlOverlayHtml: data.html,
				htmlOverlayCss: data.css,
				htmlOverlayJs: data.js,
				htmlOverlayDataPollSecondsInterval: data.refreshInterval,
			},
		})
		const layer = builder.project.layers.find((item) => item.id === editorLayer.value?.id)
		if (layer) sync.sendLayerUpdate(layer)
	}

	function handleKeyDown(event: KeyboardEvent) {
		const target = event.target
		const isInputFocused = target instanceof HTMLElement && (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable)
		if ((event.ctrlKey || event.metaKey) && event.key === 's') {
			event.preventDefault()
			handleSave()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'z' && !event.shiftKey && !isInputFocused) {
			event.preventDefault()
			handleUndo()
		} else if ((event.ctrlKey || event.metaKey) && (event.key === 'y' || (event.key === 'z' && event.shiftKey)) && !isInputFocused) {
			event.preventDefault()
			handleRedo()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'c' && !isInputFocused) {
			event.preventDefault()
			builder.copyToClipboard()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'x' && !isInputFocused) {
			event.preventDefault()
			handleCutSelection()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'v' && !isInputFocused) {
			event.preventDefault()
			handlePaste()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'd' && !isInputFocused) {
			event.preventDefault()
			if (builder.canvasState.selectedLayerIds.length > 0) handleDuplicateLayers(builder.canvasState.selectedLayerIds)
		} else if (event.key === 'Escape' && !isInputFocused && !showShortcuts.value) {
			builder.deselectAll()
		} else if (event.key === '?' && !isInputFocused) {
			event.preventDefault()
			showShortcuts.value = !showShortcuts.value
		} else if ((event.key === 'Delete' || event.key === 'Backspace') && !isInputFocused) {
			if (builder.canvasState.selectedLayerIds.length > 0) {
				event.preventDefault()
				handleRemoveLayers(builder.canvasState.selectedLayerIds)
			}
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'a' && !isInputFocused) {
			event.preventDefault()
			builder.selectAll()
		}
	}

	onMounted(() => window.addEventListener('keydown', handleKeyDown))
	onUnmounted(() => window.removeEventListener('keydown', handleKeyDown))

	return {
		builder,
		overlayName,
		instaSave,
		canvasAreaRef,
		addLayersHidden,
		showCodeEditor,
		showShortcuts,
		editorLayer,
		syncStatus: sync.status,
		hasSelection: computed(() => builder.canvasState.selectedLayerIds.length > 0),
		canAlign: computed(() => builder.canvasState.selectedLayerIds.length >= 1),
		canDistribute: computed(() => builder.canvasState.selectedLayerIds.length >= 3),
		addLayer,
		handleSave,
		handleUpdateLayer,
		handleSelectLayer,
		handleLayersPanelSelect,
		handleDeselectAll,
		handleFindGuides,
		handleClearGuides,
		handleToggleVisibility,
		handleToggleLock,
		handleRemoveLayer,
		handleRemoveLayers,
		handleDuplicateLayers,
		handleCutSelection,
		handlePaste,
		handleAlign,
		handleDistribute,
		handleUndo,
		handleRedo,
		handleOpenLayerSettings,
		handleReorderLayers,
		handleUpdateLayerProperties,
		handleActiveLayerUpdate,
		handleOpenCodeEditor,
		handleSaveCode,
	}
}
