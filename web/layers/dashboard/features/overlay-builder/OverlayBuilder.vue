<script setup lang="ts">
import { toRef } from 'vue'

import BuilderToolbar from './components/BuilderToolbar.vue'
import LayersPanel from './components/LayersPanel.vue'
import LayerPropertiesCard from './components/LayerPropertiesCard.vue'
import Canvas from './components/Canvas.vue'
import CodeEditorDialog from './components/CodeEditorDialog.vue'
import ShortcutsDialog from './components/ShortcutsDialog.vue'
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
	save: [project: OverlayProject, options?: { silent?: boolean }]
	instantSave: [project: OverlayProject]
}>()

const {
	builder,
	overlayName,
	instaSave,
	canvasAreaRef,
	addLayersHidden,
	showCodeEditor,
	showShortcuts,
	editorLayer,
	syncStatus,
	hasSelection,
	canAlign,
	canDistribute,
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
			:sync-status="syncStatus"
			@save="handleSave"
			@undo="handleUndo"
			@redo="handleRedo"
			@copy="builder.copyToClipboard"
			@cut="handleCutSelection"
			@paste="handlePaste"
			@delete="handleRemoveLayers(builder.canvasState.selectedLayerIds)"
			@duplicate="handleDuplicateLayers(builder.canvasState.selectedLayerIds)"
			@align-left="handleAlign('left')"
			@align-center="handleAlign('center')"
			@align-right="handleAlign('right')"
			@align-top="handleAlign('top')"
			@align-middle="handleAlign('middle')"
			@align-bottom="handleAlign('bottom')"
			@distribute-horizontal="handleDistribute('horizontal')"
			@distribute-vertical="handleDistribute('vertical')"
			@zoom-in="builder.zoomIn"
			@zoom-out="builder.zoomOut"
			@reset-zoom="builder.resetZoom"
			@toggle-grid="builder.canvasState.showGrid = !builder.canvasState.showGrid"
			@toggle-snap="builder.canvasState.snapToGrid = !builder.canvasState.snapToGrid"
			@toggle-add-layers-hidden="addLayersHidden = !addLayersHidden"
			@open-shortcuts="showShortcuts = true"
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
						@select="handleLayersPanelSelect"
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

		<ShortcutsDialog v-model:open="showShortcuts" />
	</div>
</template>
