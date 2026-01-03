<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'

import OverlayBuilder from '@/features/overlay-builder/OverlayBuilder.vue'
import {
	useChannelOverlayByIdQuery,
	useChannelOverlayCreate,
	useChannelOverlayUpdate,
	useChannelOverlaysQuery,
} from '@/api/index.js'
import type { ChannelOverlayLayerInput } from '@/gql/graphql'

const route = useRoute()
const messages = useMessage()

// Get overlay ID from route
const overlayId = computed(() => {
	const id = route.params.id
	if (typeof id !== 'string' || id === 'new') {
		return ''
	}
	return id
})

// Fetch existing overlay data
const { data: overlayData } = useChannelOverlayByIdQuery(overlayId)
const overlay = computed(() => overlayData.value?.channelOverlayById)

// Mutations for creating/updating overlays
const { executeQuery: refetchOverlays } = useChannelOverlaysQuery()
const createOverlayMutation = useChannelOverlayCreate()
const updateOverlayMutation = useChannelOverlayUpdate()

// Check if we're creating a new overlay or editing existing one
const isNewOverlay = computed(() => overlayId.value === '')

// Convert existing overlay data to builder format
const projectData = computed(() => {
	// For new overlay, return empty project immediately
	if (isNewOverlay.value) {
		return {
			id: '',
			name: '',
			width: 1920,
			height: 1080,
			layers: [],
		}
	}

	// For existing overlay, wait for data to load
	if (!overlay.value) {
		return null
	}

	console.log('[edit.vue] Loading overlay data from API:', overlay.value)

	// Convert existing overlay data to builder format (canvas size fixed at 1920x1080)
	const converted = {
		id: overlay.value.id,
		name: overlay.value.name,
		width: 1920,
		height: 1080,
		layers: overlay.value.layers.map((layer, index) => {
			console.log('[edit.vue] Converting layer from API:', {
				type: layer.type,
				imageUrl: layer.settings.imageUrl,
				fullSettings: layer.settings,
			})
			return {
				id: `layer-${layer.id || index}`,
				type: layer.type,
				name: `${layer.type} Layer ${index + 1}`,
				posX: layer.posX,
				posY: layer.posY,
				width: layer.width,
				height: layer.height,
				rotation: Number(layer.rotation) || 0,
				opacity: 1,
				visible: true,
				locked: false,
				zIndex: index,
				periodicallyRefetchData: layer.periodicallyRefetchData,
				settings: {
					htmlOverlayHtml: layer.settings.htmlOverlayHtml || '',
					htmlOverlayCss: layer.settings.htmlOverlayCss || '',
					htmlOverlayJs: layer.settings.htmlOverlayJs || '',
					htmlOverlayDataPollSecondsInterval: layer.settings.htmlOverlayDataPollSecondsInterval || 5,
					imageUrl: layer.settings.imageUrl || '',
				},
			}
		}),
	}

	console.log('[edit.vue] Converted project data:', converted)
	return converted
})

// Handle save from builder
async function handleSave(project: any) {
	console.log('[edit.vue] Saving overlay project:', project)

	// Validate project data
	if (!project.name || project.name.length > 30) {
		messages.error('Overlay name is required and must be less than 30 characters')
		return
	}

	if (!project.layers.length || project.layers.length > 15) {
		messages.error('Overlay must have between 1 and 15 layers')
		return
	}

	// Convert builder format back to API format
	const layersInput: ChannelOverlayLayerInput[] = project.layers.map((layer: any) => {
		const rotation = Number(layer.rotation ?? 0)
		console.log('[edit.vue] Converting layer to API format:', {
			type: layer.type,
			imageUrl: layer.settings?.imageUrl,
			fullSettings: layer.settings,
		})
		return {
			type: layer.type,
			posX: layer.posX,
			posY: layer.posY,
			width: layer.width,
			height: layer.height,
			rotation: rotation,
			periodicallyRefetchData: layer.periodicallyRefetchData,
			settings: {
				htmlOverlayHtml: layer.settings?.htmlOverlayHtml ?? '',
				htmlOverlayCss: layer.settings?.htmlOverlayCss ?? '',
				htmlOverlayJs: layer.settings?.htmlOverlayJs ?? '',
				htmlOverlayDataPollSecondsInterval: layer.settings?.htmlOverlayDataPollSecondsInterval ?? 5,
				imageUrl: layer.settings?.imageUrl ?? '',
			},
		}
	})

	console.log('[edit.vue] Layers to be saved:', layersInput)

	try {
		if (project.id) {
			// Update existing overlay
			const mutationInput = {
				id: project.id,
				input: {
					name: project.name,
					width: 1920,
					height: 1080,
					layers: layersInput,
				},
			}

			const result = await updateOverlayMutation.executeMutation(mutationInput)

			if (result.error) {
				messages.error(result.error.message)
				return
			}

			messages.success('Overlay updated successfully!')
		} else {
			// Create new overlay
			const mutationInput = {
				input: {
					name: project.name,
					width: 1920,
					height: 1080,
					layers: layersInput,
				},
			}

			const result = await createOverlayMutation.executeMutation(mutationInput)

			if (result.error) {
				messages.error(result.error.message)
				return
			}

			messages.success('Overlay created successfully!')
		}

		// Refresh the overlays list
		refetchOverlays({ requestPolicy: 'network-only' })
	} catch (error) {
		console.error('Error saving overlay:', error)
		messages.error('Failed to save overlay')
	}
}
</script>

<template>
	<div class="fixed inset-0 w-full h-full overflow-hidden">
		<OverlayBuilder
			v-if="projectData"
			:initial-project="projectData"
			@save="handleSave"
		/>
		<div v-else class="flex items-center justify-center w-full h-full">
			<p class="text-muted-foreground">Loading overlay...</p>
		</div>
	</div>
</template>
