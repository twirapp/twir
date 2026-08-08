import { type MaybeRefOrGetter, computed, ref, toValue } from 'vue'
import { toast } from 'vue-sonner'

import type { ChannelOverlayLayerInput } from '~/gql/graphql.js'

import {
	useChannelOverlayCreate,
	useChannelOverlayUpdate,
	useChannelOverlaysQuery,
} from '~~/layers/dashboard/api/overlays/custom'

import type { OverlayProject } from '../types'

import { useOverlayInstantSave } from './useOverlayInstantSave'

export function useOverlaySave(overlayId: MaybeRefOrGetter<string>) {
	const createOverlayMutation = useChannelOverlayCreate()
	const updateOverlayMutation = useChannelOverlayUpdate()
	const { executeQuery: refetchOverlays } = useChannelOverlaysQuery()

	const currentOverlayId = computed(() => toValue(overlayId))
	const { sendLayerPositions, isEnabled: instaSaveEnabled } = useOverlayInstantSave(overlayId)

	const isSaving = ref(false)

	// Convert project to GraphQL input
	function projectToLayersInput(project: OverlayProject): ChannelOverlayLayerInput[] {
		return project.layers.map((layer, index) => {
			const rotation = Number(layer.rotation ?? 0)
			return {
				id: layer.id ?? undefined, // Include layer ID for updates
				type: layer.type,
				name: layer.name,
				posX: layer.posX,
				posY: layer.posY,
				width: layer.width,
				height: layer.height,
				rotation: rotation,
				periodicallyRefetchData: layer.periodicallyRefetchData,
				locked: layer.locked ?? false,
				visible: layer.visible ?? true,
				opacity: layer.opacity ?? 1.0,
				// index is canonical: layer.zIndex mirrors it only after reindexLayers() runs
				zIndex: index,
				settings: {
					htmlOverlayHtml: layer.settings.htmlOverlayHtml,
					htmlOverlayCss: layer.settings.htmlOverlayCss,
					htmlOverlayJs: layer.settings.htmlOverlayJs,
					htmlOverlayDataPollSecondsInterval:
						layer.settings.htmlOverlayDataPollSecondsInterval,
					imageUrl: layer.settings.imageUrl,
					textContent: layer.settings.textContent,
					textFontFamily: layer.settings.textFontFamily,
					textFontSize: layer.settings.textFontSize,
					textFontWeight: layer.settings.textFontWeight,
					textColor: layer.settings.textColor,
					textAlign: layer.settings.textAlign,
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
			}
		})
	}

	// Save full overlay (creates new or updates existing)
	async function saveOverlay(project: OverlayProject): Promise<string | null> {
		// Validate project data
		if (!project.name || project.name.length > 30) {
			toast.error('Overlay name is required and must be less than 30 characters')
			return null
		}

		if (!project.layers.length || project.layers.length > 15) {
			toast.error('Overlay must have between 1 and 15 layers')
			return null
		}

		isSaving.value = true

		try {
			const layersInput = projectToLayersInput(project)

			if (project.id) {
				// Update existing overlay
				const result = await updateOverlayMutation.executeMutation({
					id: project.id,
					input: {
						name: project.name,
						width: project.width,
						height: project.height,
						instaSave: project.instaSave || false,
						layers: layersInput,
					},
				})

				if (result.error) {
					toast.error(result.error.message)
					return null
				}

				toast.success('Overlay updated successfully!')
				refetchOverlays({ requestPolicy: 'network-only' })
				return project.id
			} else {
				// Create new overlay
				const result = await createOverlayMutation.executeMutation({
					input: {
						name: project.name,
						width: project.width,
						height: project.height,
						instaSave: project.instaSave || false,
						layers: layersInput,
					},
				})

				if (result.error) {
					toast.error(result.error.message)
					return null
				}

				toast.success('Overlay created successfully!')
				refetchOverlays({ requestPolicy: 'network-only' })
				return result.data?.channelOverlayCreate?.id ?? null
			}
		} catch (error) {
			console.error('[OverlaySave] Error saving overlay:', error)
			toast.error('Failed to save overlay')
			return null
		} finally {
			isSaving.value = false
		}
	}

	// Save only instaSave setting (when toggling the checkbox)
	async function saveInstaSaveSetting(project: OverlayProject): Promise<boolean> {
		if (!project.id) {
			console.error('[OverlaySave] Cannot save instaSave: no project ID')
			return false
		}

		try {
			const layersInput = projectToLayersInput(project)

			const result = await updateOverlayMutation.executeMutation({
				id: project.id,
				input: {
					name: project.name,
					width: project.width,
					height: project.height,
					instaSave: project.instaSave || false,
					layers: layersInput,
				},
			})

			if (result.error) {
				console.error('[OverlaySave] Error saving instaSave:', result.error)
				return false
			}

			instaSaveEnabled.value = project.instaSave || false
			refetchOverlays({ requestPolicy: 'network-only' })
			return true
		} catch (error) {
			console.error('[OverlaySave] Error saving instaSave:', error)
			return false
		}
	}

	// Debounced instant save for position/rotation changes
	const instantSavePositions = (project: OverlayProject) => {
		if (!project.id) {
			console.log('[OverlaySave] Cannot instant save: no project ID')
			return
		}

		if (!instaSaveEnabled.value) {
			console.log('[OverlaySave] Instant save is disabled')
			return
		}

		// Verify overlayId matches project.id before sending
		if (currentOverlayId.value !== project.id) {
			console.log('[OverlaySave] OverlayId mismatch, skipping instant save')
			return
		}

		sendLayerPositions(project)
	}

	return {
		isSaving,
		saveOverlay,
		saveInstaSaveSetting,
		instantSavePositions,
		instaSaveEnabled,
	}
}
