<script setup lang="ts">
import { toRef } from 'vue'

import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { useBuilderToolbar } from '../composables/useBuilderToolbar'
import { useWidgetPreviewMode } from '../composables/useWidgetPreviewMode'

interface Props {
	canUndo: boolean
	canRedo: boolean
	hasSelection: boolean
	canAlign: boolean
	canDistribute: boolean
	zoom: number
	showGrid: boolean
	snapToGrid: boolean
	addLayersHidden: boolean
	overlayId?: string
	overlayName?: string
	syncStatus?: 'OPEN' | 'CONNECTING' | 'CLOSED'
}

const props = defineProps<Props>()

const emit = defineEmits<{
	save: []
	undo: []
	redo: []
	copy: []
	cut: []
	paste: []
	delete: []
	duplicate: []
	alignLeft: []
	alignCenter: []
	alignRight: []
	alignTop: []
	alignMiddle: []
	alignBottom: []
	distributeHorizontal: []
	distributeVertical: []
	zoomIn: []
	zoomOut: []
	resetZoom: []
	toggleGrid: []
	toggleSnap: []
	toggleAddLayersHidden: []
	openShortcuts: []
}>()

const { t } = useI18n()
const { formatZoom, goBack, copyOverlayLink } = useBuilderToolbar(toRef(props, 'overlayId'))
const { enabled: widgetsPreview } = useWidgetPreviewMode()

const syncStatusMeta = computed(() => {
	switch (props.syncStatus) {
		case 'OPEN':
			return { class: 'bg-green-500', label: t('overlayBuilder.sync.connected') }
		case 'CONNECTING':
			return { class: 'bg-yellow-500 animate-pulse', label: t('overlayBuilder.sync.connecting') }
		default:
			return { class: 'bg-muted-foreground/40', label: t('overlayBuilder.sync.disconnected') }
	}
})
</script>

<template>
	<div class="bg-background flex h-14 items-center gap-2 border-b px-4 py-2">
		<!-- Back Button -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						@click="goBack"
					>
						<Icon
							name="lucide:arrow-left"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.back') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator
			orientation="vertical"
			class="h-6"
		/>

		<!-- Undo/Redo -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canUndo"
						@click="emit('undo')"
					>
						<Icon
							name="lucide:undo"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.undo') }} (Ctrl+Z)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canRedo"
						@click="emit('redo')"
					>
						<Icon
							name="lucide:redo"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.redo') }} (Ctrl+Y)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator
			orientation="vertical"
			class="h-6"
		/>

		<!-- Clipboard -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!hasSelection"
						@click="emit('copy')"
					>
						<Icon
							name="lucide:copy"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.copy') }} (Ctrl+C)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!hasSelection"
						@click="emit('cut')"
					>
						<Icon
							name="lucide:scissors"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.cut') }} (Ctrl+X)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!hasSelection"
						@click="emit('duplicate')"
					>
						<Icon
							name="lucide:copy-plus"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.duplicate') }} (Ctrl+D)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!hasSelection"
						@click="emit('delete')"
					>
						<Icon
							name="lucide:trash"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.delete') }} (Del)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator
			orientation="vertical"
			class="h-6"
		/>

		<!-- Alignment -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canAlign"
						@click="emit('alignLeft')"
					>
						<Icon
							name="lucide:align-left"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.alignLeft') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canAlign"
						@click="emit('alignCenter')"
					>
						<Icon
							name="lucide:align-center"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.alignCenter') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canAlign"
						@click="emit('alignRight')"
					>
						<Icon
							name="lucide:align-right"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.alignRight') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator
			orientation="vertical"
			class="h-6"
		/>

		<!-- Vertical Alignment -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canAlign"
						@click="emit('alignTop')"
					>
						<Icon
							name="lucide:align-start-vertical"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.alignTop') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canAlign"
						@click="emit('alignMiddle')"
					>
						<Icon
							name="lucide:align-center-vertical"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.alignMiddle') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canAlign"
						@click="emit('alignBottom')"
					>
						<Icon
							name="lucide:align-end-vertical"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.alignBottom') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator
			orientation="vertical"
			class="h-6"
		/>

		<!-- Distribution -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canDistribute"
						@click="emit('distributeHorizontal')"
					>
						<Icon
							name="lucide:align-horizontal-distribute-center"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.distributeHorizontal') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:disabled="!canDistribute"
						@click="emit('distributeVertical')"
					>
						<Icon
							name="lucide:align-vertical-distribute-center"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.distributeVertical') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator
			orientation="vertical"
			class="h-6"
		/>

		<!-- Zoom -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						@click="emit('zoomOut')"
					>
						<Icon
							name="lucide:minus"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.zoomOut') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Button
			variant="ghost"
			size="sm"
			class="min-w-16"
			@click="emit('resetZoom')"
		>
			{{ formatZoom(zoom) }}
		</Button>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						@click="emit('zoomIn')"
					>
						<Icon
							name="lucide:plus"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.zoomIn') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator
			orientation="vertical"
			class="h-6"
		/>

		<!-- Grid -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:class="{ 'bg-accent': showGrid }"
						@click="emit('toggleGrid')"
					>
						<Icon
							name="lucide:grid-3x3"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.toggleGrid') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:class="{ 'bg-accent': snapToGrid }"
						@click="emit('toggleSnap')"
					>
						<Icon
							name="lucide:layers"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.snapToGrid') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						:class="{ 'bg-accent': widgetsPreview }"
						@click="widgetsPreview = !widgetsPreview"
					>
						<Icon
							name="lucide:play"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.widgetsPreview') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<div class="flex items-center gap-2 px-2 text-xs text-muted-foreground">
			<Switch
				id="add-layers-hidden"
				:model-value="addLayersHidden"
				@update:model-value="emit('toggleAddLayersHidden')"
			/>
			<Label for="add-layers-hidden" class="cursor-pointer whitespace-nowrap">{{ t('overlayBuilder.toolbar.addLayersHidden') }}</Label>
		</div>

		<div class="flex-1" />

		<!-- Right Side Actions -->
		<TooltipProvider v-if="overlayId && syncStatus">
			<Tooltip>
				<TooltipTrigger as-child>
					<div class="flex items-center px-1">
						<span class="h-2 w-2 rounded-full" :class="syncStatusMeta.class" />
					</div>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ syncStatusMeta.label }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						@click="emit('openShortcuts')"
					>
						<Icon
							name="lucide:keyboard"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.shortcuts') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider v-if="overlayId">
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="ghost"
						size="icon"
						@click="copyOverlayLink"
					>
						<Icon
							name="lucide:external-link"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlays.copyOverlayLink') }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button
						variant="default"
						size="icon"
						@click="emit('save')"
					>
						<Icon
							name="lucide:save"
							class="h-4 w-4"
						/>
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlayBuilder.toolbar.save') }} (Ctrl+S)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>
	</div>
</template>
