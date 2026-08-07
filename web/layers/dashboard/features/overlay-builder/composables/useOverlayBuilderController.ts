import { type Ref, computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import type { ChannelOverlayLayer, ChannelOverlayLayerType } from '~/gql/graphql.js'
import { useProfile } from '~~/layers/dashboard/api/auth.js'

import { useOverlayBuilder } from './useOverlayBuilder'
import { useChatOverlayPresetQuery } from './useChatOverlayPresets'
import { resolveWidgetLayerUrl } from './widget-url'
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
	const addLayersHidden = ref(false)
	const showCodeEditor = ref(false)
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
		builder.addLayer(type, {
			name: t('overlayBuilder.layerNames.default', {
				type: t(getLayerTypeMeta(type).labelKey),
				count: builder.project.layers.length + 1,
			}),
			settings: { textContent: t('overlayBuilder.defaults.textContent') },
			visible: !addLayersHidden.value,
		})
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
		if (newValue !== oldValue && loadedProjectId.value) emit('instantSave', projectSnapshot(overlayName.value, newValue))
	})
	watch([() => builder.project.width, () => builder.project.height], () => {
		if (isLoadingProject.value) return
		builder.constrainLayersToCanvas()
		void nextTick(calculateFitZoom)
		if (instaSave.value && loadedProjectId.value) emit('save', projectSnapshot(overlayName.value, instaSave.value))
	})

	function handleUpdateLayer(layerId: string, updates: Partial<Layer>) {
		builder.updateLayer(layerId, updates)
		if (updates.posX !== undefined || updates.posY !== undefined || updates.rotation !== undefined || updates.width !== undefined || updates.height !== undefined || updates.opacity !== undefined || updates.visible !== undefined) void handleLayerUpdate()
	}

	function handleSelectLayer(layerId: string, addToSelection: boolean) {
		builder.selectLayers([layerId], addToSelection)
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
		void handleLayerUpdate()
	}

	function handleToggleLock(layerId: string) {
		const layer = builder.project.layers.find((item) => item.id === layerId)
		if (layer) builder.updateLayer(layerId, { locked: !layer.locked })
	}

	function handleRemoveLayer(layerId: string) {
		builder.removeLayer(layerId)
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
	}

	function handleUpdateLayerProperties(layerId: string, updates: Partial<Layer>) {
		builder.updateLayer(layerId, updates)
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
	}

	function handleKeyDown(event: KeyboardEvent) {
		const target = event.target
		const isInputFocused = target instanceof HTMLElement && (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable)
		if ((event.ctrlKey || event.metaKey) && event.key === 's') {
			event.preventDefault()
			handleSave()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'z' && !event.shiftKey && !isInputFocused) {
			event.preventDefault()
			builder.undo()
		} else if ((event.ctrlKey || event.metaKey) && (event.key === 'y' || (event.key === 'z' && event.shiftKey)) && !isInputFocused) {
			event.preventDefault()
			builder.redo()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'c' && !isInputFocused) {
			event.preventDefault()
			builder.copyToClipboard()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'x' && !isInputFocused) {
			event.preventDefault()
			builder.cutToClipboard()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'v' && !isInputFocused) {
			event.preventDefault()
			builder.pasteFromClipboard()
		} else if ((event.ctrlKey || event.metaKey) && event.key === 'd' && !isInputFocused) {
			event.preventDefault()
			if (builder.canvasState.selectedLayerIds.length > 0) builder.duplicateLayers(builder.canvasState.selectedLayerIds)
		} else if ((event.key === 'Delete' || event.key === 'Backspace') && !isInputFocused) {
			if (builder.canvasState.selectedLayerIds.length > 0) {
				event.preventDefault()
				builder.removeLayers(builder.canvasState.selectedLayerIds)
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
		editorLayer,
		hasSelection: computed(() => builder.canvasState.selectedLayerIds.length > 0),
		canAlign: computed(() => builder.canvasState.selectedLayerIds.length >= 1),
		canDistribute: computed(() => builder.canvasState.selectedLayerIds.length >= 3),
		addLayer,
		handleSave,
		handleUpdateLayer,
		handleSelectLayer,
		handleDeselectAll,
		handleFindGuides,
		handleClearGuides,
		handleToggleVisibility,
		handleToggleLock,
		handleRemoveLayer,
		handleOpenLayerSettings,
		handleReorderLayers,
		handleUpdateLayerProperties,
		handleActiveLayerUpdate,
		handleOpenCodeEditor,
		handleSaveCode,
	}
}
