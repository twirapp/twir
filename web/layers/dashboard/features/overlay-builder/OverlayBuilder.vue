<script setup lang="ts">
import { toRef } from 'vue'

import BuilderToolbar from './components/BuilderToolbar.vue'
import LayersPanel from './components/LayersPanel.vue'
import LayerPropertiesCard from './components/LayerPropertiesCard.vue'
import Canvas from './components/Canvas.vue'
import CodeEditorDialog from './components/CodeEditorDialog.vue'
import OverlaySettings from './components/OverlaySettings.vue'
import type { OverlayProject } from './types'
import {
	type InitialOverlayProject,
	useOverlayBuilderController,
} from './composables/useOverlayBuilderController'

interface Props {
	initialProject?: InitialOverlayProject
}

const props = defineProps<Props>()

const emit = defineEmits<{
	save: [project: OverlayProject]
	instantSave: [project: OverlayProject]
}>()

const {
	builder,
	overlayName,
	instaSave,
	canvasAreaRef,
	addLayersHidden,
	showCodeEditor,
	editorLayer,
	hasSelection,
	canAlign,
	canDistribute,
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
} = useOverlayBuilderController(toRef(props, 'initialProject'), emit)
</script>

<template>
	<div class="w-full h-full flex flex-col bg-background overflow-hidden">
		<BuilderToolbar
			:can-undo="builder.canUndo.value"
			:can-redo="builder.canRedo.value"
			:has-selection="hasSelection"
			:can-align="canAlign"
			:can-distribute="canDistribute"
			:zoom="builder.canvasState.zoom"
			:show-grid="builder.canvasState.showGrid"
			:snap-to-grid="builder.canvasState.snapToGrid"
			:add-layers-hidden="addLayersHidden"
			:overlay-id="initialProject?.id"
			:overlay-name="overlayName"
			@save="handleSave"
			@undo="builder.undo"
			@redo="builder.redo"
			@copy="builder.copyToClipboard"
			@cut="builder.cutToClipboard"
			@paste="builder.pasteFromClipboard"
			@delete="builder.removeLayers(builder.canvasState.selectedLayerIds)"
			@duplicate="builder.duplicateLayers(builder.canvasState.selectedLayerIds)"
			@align-left="builder.alignLayers('left')"
			@align-center="builder.alignLayers('center')"
			@align-right="builder.alignLayers('right')"
			@align-top="builder.alignLayers('top')"
			@align-middle="builder.alignLayers('middle')"
			@align-bottom="builder.alignLayers('bottom')"
			@distribute-horizontal="builder.distributeLayersHorizontally"
			@distribute-vertical="builder.distributeLayersVertically"
			@zoom-in="builder.zoomIn"
			@zoom-out="builder.zoomOut"
			@reset-zoom="builder.resetZoom"
			@toggle-grid="builder.canvasState.showGrid = !builder.canvasState.showGrid"
			@toggle-snap="builder.canvasState.snapToGrid = !builder.canvasState.snapToGrid"
			@toggle-add-layers-hidden="addLayersHidden = !addLayersHidden"
		/>

		<div class="flex-1 flex overflow-hidden">
			<Canvas
				ref="canvasAreaRef"
				:layers="builder.project.layers"
				:selected-layer-ids="builder.canvasState.selectedLayerIds"
				:zoom="builder.canvasState.zoom"
				:pan-x="builder.canvasState.panX"
				:pan-y="builder.canvasState.panY"
				:canvas-width="builder.project.width"
				:canvas-height="builder.project.height"
				:show-grid="builder.canvasState.showGrid"
				:snap-to-grid="builder.canvasState.snapToGrid"
				:grid-size="builder.canvasState.gridSize"
				:alignment-guides="builder.alignmentGuides.value"
				:snap-to-guides-enabled="builder.canvasState.showGuides"
				@update-layer="handleUpdateLayer"
				@select-layer="handleSelectLayer"
				@deselect-all="handleDeselectAll"
				@find-guides="handleFindGuides"
				@clear-guides="handleClearGuides"
				@toggle-visibility="handleToggleVisibility"
				@toggle-lock="handleToggleLock"
				@remove-layer="handleRemoveLayer"
				@open-layer-settings="handleOpenLayerSettings"
			>
				<template #layer-content="{ layer }">
					<div class="w-full h-full flex items-center justify-center text-white/70 text-sm">{{ layer.name }}</div>
				</template>
			</Canvas>

			<div class="w-80 flex flex-col border-l">
				<div class="border-b bg-background p-2">
					<OverlaySettings v-model:overlay-name="overlayName" v-model:insta-save="instaSave" v-model:canvas-width="builder.project.width" v-model:canvas-height="builder.project.height" />
				</div>

				<div class="flex min-h-0 flex-1 flex-col">
					<div class="flex-1 min-h-0 overflow-hidden p-2">
						<LayersPanel
							:layers="builder.project.layers"
							:selected-layer-ids="builder.canvasState.selectedLayerIds"
							@select="handleSelectLayer"
							@toggle-visibility="handleToggleVisibility"
							@toggle-lock="handleToggleLock"
							@reorder="handleReorderLayers"
							@add-layer="addLayer"
							@update-layer-properties="handleUpdateLayerProperties"
						/>
					</div>

					<div v-if="builder.activeLayer.value" class="flex max-h-[60%] min-h-0 shrink-0 flex-col overflow-hidden p-2 pt-1">
						<LayerPropertiesCard :layer="builder.activeLayer.value" @update="handleActiveLayerUpdate" @open-code-editor="handleOpenCodeEditor" />
					</div>
				</div>
			</div>
		</div>

		<CodeEditorDialog
			v-model:open="showCodeEditor"
			:layer-id="editorLayer?.id"
			:layer-name="editorLayer?.name"
			:html="editorLayer?.settings?.htmlOverlayHtml"
			:css="editorLayer?.settings?.htmlOverlayCss"
			:js="editorLayer?.settings?.htmlOverlayJs"
			:refresh-interval="editorLayer?.settings?.htmlOverlayDataPollSecondsInterval"
			@save="handleSaveCode"
		/>
	</div>
</template>
