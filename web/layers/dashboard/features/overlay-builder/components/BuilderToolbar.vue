<script setup lang="ts">
import {
	AlignCenter,
	AlignHorizontalDistributeCenter,
	AlignLeft,
	AlignRight,
	AlignVerticalDistributeCenter,
	ArrowLeft,
	Copy,
	CopyPlus,
	ExternalLink,
	Grid3x3,
	Layers,
	Minus,
	Plus,
	Redo,
	Save,
	Scissors,
	Trash2,
	Undo,
} from 'lucide-vue-next'
import { computed } from 'vue'

import { useRouter } from 'vue-router'

// ...existing code...

interface Props {
	canUndo: boolean
	canRedo: boolean
	hasSelection: boolean
	canAlign: boolean
	canDistribute: boolean
	zoom: number
	showGrid: boolean
	snapToGrid: boolean
	overlayId?: string
	overlayName?: string
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
}>()

const { t } = useI18n()
const router = useRouter()
const { user: profile } = storeToRefs(useDashboardAuth())

const selectedDashboardUser = computed(() => {
	return profile.value?.availableDashboards.find(
		(dashboard) => dashboard.id === profile.value?.selectedDashboardId
	)
})

const formatZoom = computed(() => (zoom: number) => `${Math.round(zoom * 100)}%`)

function goBack() {
	router.push('/dashboard/overlays')
}

function copyOverlayLink() {
	if (!props.overlayId || !selectedDashboardUser.value?.apiKey) return

	const baseUrl = window.location.origin
	const overlayUrl = `${baseUrl}/overlays/${selectedDashboardUser.value.apiKey}/registry/overlays/${props.overlayId}`

	navigator.clipboard.writeText(overlayUrl).then(() => {
		// Use vue-sonner toast if available, fallback to message
		const toastModule = import('vue-sonner')
		toastModule.then(({ toast }) => {
			toast.success(t('sharedTexts.copied') || 'Link copied to clipboard!')
		}).catch(() => {
			console.log('Link copied to clipboard!')
		})
	}).catch(() => {
		console.error('Failed to copy link')
	})
}
</script>

<template>
	<div class="flex items-center gap-2 bg-background border-b px-4 py-2 h-14">
		<!-- Back Button -->
		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" @click="goBack">
						<ArrowLeft class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>{{ t('sharedButtons.back') || 'Back to Overlays' }}</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiSeparator orientation="vertical" class="h-6" />

		<!-- Undo/Redo -->
		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" :disabled="!canUndo" @click="emit('undo')">
						<Undo class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Undo (Ctrl+Z)</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" :disabled="!canRedo" @click="emit('redo')">
						<Redo class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Redo (Ctrl+Y)</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiSeparator orientation="vertical" class="h-6" />

		<!-- Clipboard -->
		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" :disabled="!hasSelection" @click="emit('copy')">
						<Copy class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Copy (Ctrl+C)</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" :disabled="!hasSelection" @click="emit('cut')">
						<Scissors class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Cut (Ctrl+X)</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" :disabled="!hasSelection" @click="emit('duplicate')">
						<CopyPlus class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Duplicate (Ctrl+D)</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" :disabled="!hasSelection" @click="emit('delete')">
						<Trash2 class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Delete (Del)</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiSeparator orientation="vertical" class="h-6" />

		<!-- Alignment -->
		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" :disabled="!canAlign" @click="emit('alignLeft')">
						<AlignLeft class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Align Left</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" :disabled="!canAlign" @click="emit('alignCenter')">
						<AlignCenter class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Align Center</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" :disabled="!canAlign" @click="emit('alignRight')">
						<AlignRight class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Align Right</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiSeparator orientation="vertical" class="h-6" />

		<!-- Distribution -->
		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton
						variant="ghost"
						size="icon"
						:disabled="!canDistribute"
						@click="emit('distributeHorizontal')"
					>
						<AlignHorizontalDistributeCenter class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Distribute Horizontally</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton
						variant="ghost"
						size="icon"
						:disabled="!canDistribute"
						@click="emit('distributeVertical')"
					>
						<AlignVerticalDistributeCenter class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Distribute Vertically</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiSeparator orientation="vertical" class="h-6" />

		<!-- Zoom -->
		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" @click="emit('zoomOut')">
						<Minus class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Zoom Out</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiButton variant="ghost" size="sm" class="min-w-16" @click="emit('resetZoom')">
			{{ formatZoom(zoom) }}
		</UiButton>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" @click="emit('zoomIn')">
						<Plus class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Zoom In</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiSeparator orientation="vertical" class="h-6" />

		<!-- Grid -->
		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton
						variant="ghost"
						size="icon"
						:class="{ 'bg-accent': showGrid }"
						@click="emit('toggleGrid')"
					>
						<Grid3x3 class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Toggle Grid</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton
						variant="ghost"
						size="icon"
						:class="{ 'bg-accent': snapToGrid }"
						@click="emit('toggleSnap')"
					>
						<Layers class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>Snap to Grid</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<div class="flex-1" />

		<!-- Right Side Actions -->
		<UiTooltipProvider v-if="overlayId">
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="ghost" size="icon" @click="copyOverlayLink">
						<ExternalLink class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>{{ t('overlaysRegistry.copyLink') || 'Copy Overlay Link' }}</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>

		<UiTooltipProvider>
			<UiTooltip>
				<UiTooltipTrigger as-child>
					<UiButton variant="default" size="icon" @click="emit('save')">
						<Save class="h-4 w-4" />
					</UiButton>
				</UiTooltipTrigger>
				<UiTooltipContent>
					<p>{{ t('sharedButtons.save') || 'Save' }} (Ctrl+S)</p>
				</UiTooltipContent>
			</UiTooltip>
		</UiTooltipProvider>
	</div>
</template>
