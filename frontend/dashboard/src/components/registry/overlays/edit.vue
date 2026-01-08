<script setup lang="ts">
import { computed, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import OverlayBuilder from '@/features/overlay-builder/OverlayBuilder.vue'
import {
	useChannelOverlayByIdQuery,
	useChannelOverlaysQuery,
} from '@/api/overlays/custom'
import type { OverlayProject } from '@/features/overlay-builder/types'
import { useOverlaySave } from '@/features/overlay-builder/composables/useOverlaySave'
import { useOverlayInstantSave } from '@/features/overlay-builder/composables/useOverlayInstantSave'

const route = useRoute()
const router = useRouter()

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

// Fetch all overlays to generate default name
const { data: allOverlaysData } = useChannelOverlaysQuery()
const overlayCount = computed(() => allOverlaysData.value?.channelOverlays?.length ?? 0)

// Initialize save and instant save composables
const {
	saveOverlay,
	instantSavePositions,
	instaSaveEnabled,
} = useOverlaySave(overlayId)

const { close: closeWebSocket } = useOverlayInstantSave(overlayId)

// Check if we're creating a new overlay or editing existing one
const isNewOverlay = computed(() => overlayId.value === '')

// Convert existing overlay data to builder format
const projectData = computed(() => {
	// For new overlay, return empty project immediately
	if (isNewOverlay.value) {
		return {
			id: '',
			name: `Overlay #${overlayCount.value + 1}`,
			width: 1920,
			height: 1080,
			instaSave: false,
			layers: [],
		}
	}

	// For existing overlay, wait for data to load
	if (!overlay.value) {
		return null
	}

	// Sync instaSave state with composable
	instaSaveEnabled.value = overlay.value.instaSave || false

	// Convert existing overlay data to builder format (canvas size fixed at 1920x1080)
	const converted = {
		id: overlay.value.id,
		name: overlay.value.name,
		width: 1920,
		height: 1080,
		instaSave: overlay.value.instaSave || false,
		layers: overlay.value.layers.map((layer, index) => {
			return {
				id: layer.id, // Use real layer ID from backend
				type: layer.type,
				name: `${layer.type} Layer ${index + 1}`,
				posX: layer.posX,
				posY: layer.posY,
				width: layer.width,
				height: layer.height,
				rotation: Number(layer.rotation) || 0,
				opacity: layer.opacity ?? 1.0,
				visible: layer.visible ?? true,
				locked: layer.locked ?? false,
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

	return converted
})

// Handle save from builder
async function handleSave(project: OverlayProject) {
	const newId = await saveOverlay(project)

	// If created new overlay, redirect to edit page
	if (!project.id && newId) {
		await router.replace({
			name: 'RegistryOverlayEdit',
			params: { id: newId },
		})
	}
}
// Handle instant save (position/rotation changes OR instaSave setting toggle)
async function handleInstantSave(project: OverlayProject) {
	if (!project.id) {
		console.log('[OverlayEdit] Cannot instant save: no project ID')
		return
	}

	instantSavePositions(project)
}

onUnmounted(() => {
	closeWebSocket()
})
</script>

<template>
	<div class="fixed inset-0 w-full h-full overflow-hidden">
		<OverlayBuilder
			v-if="projectData"
			:initial-project="projectData"
			@save="handleSave"
			@instant-save="handleInstantSave"
		/>
		<div v-else class="flex items-center justify-center w-full h-full">
			<p class="text-muted-foreground">Loading overlay...</p>
		</div>
	</div>
</template>
