<script setup lang="ts">
import {
	AlignCenter,
	AlignCenterVertical,
	AlignEndVertical,
	AlignHorizontalDistributeCenter,
	AlignLeft,
	AlignRight,
	AlignStartVertical,
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
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useProfile } from '@/api/auth.js'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from '@/components/ui/tooltip'

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
const { data: profile } = useProfile()

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
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" @click="goBack">
						<ArrowLeft class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('sharedButtons.back') || 'Back to Overlays' }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator orientation="vertical" class="h-6" />

		<!-- Undo/Redo -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!canUndo" @click="emit('undo')">
						<Undo class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Undo (Ctrl+Z)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!canRedo" @click="emit('redo')">
						<Redo class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Redo (Ctrl+Y)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator orientation="vertical" class="h-6" />

		<!-- Clipboard -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!hasSelection" @click="emit('copy')">
						<Copy class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Copy (Ctrl+C)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!hasSelection" @click="emit('cut')">
						<Scissors class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Cut (Ctrl+X)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!hasSelection" @click="emit('duplicate')">
						<CopyPlus class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Duplicate (Ctrl+D)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!hasSelection" @click="emit('delete')">
						<Trash2 class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Delete (Del)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator orientation="vertical" class="h-6" />

		<!-- Alignment -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!canAlign" @click="emit('alignLeft')">
						<AlignLeft class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Align Left</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!canAlign" @click="emit('alignCenter')">
						<AlignCenter class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Align Center</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!canAlign" @click="emit('alignRight')">
						<AlignRight class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Align Right</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator orientation="vertical" class="h-6" />

		<!-- Vertical Alignment -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!canAlign" @click="emit('alignTop')">
						<AlignStartVertical class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Align Top</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!canAlign" @click="emit('alignMiddle')">
						<AlignCenterVertical class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Align Middle</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" :disabled="!canAlign" @click="emit('alignBottom')">
						<AlignEndVertical class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Align Bottom</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator orientation="vertical" class="h-6" />

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
						<AlignHorizontalDistributeCenter class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Distribute Horizontally</p>
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
						<AlignVerticalDistributeCenter class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Distribute Vertically</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator orientation="vertical" class="h-6" />

		<!-- Zoom -->
		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" @click="emit('zoomOut')">
						<Minus class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Zoom Out</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Button variant="ghost" size="sm" class="min-w-16" @click="emit('resetZoom')">
			{{ formatZoom(zoom) }}
		</Button>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" @click="emit('zoomIn')">
						<Plus class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Zoom In</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<Separator orientation="vertical" class="h-6" />

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
						<Grid3x3 class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Toggle Grid</p>
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
						<Layers class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Snap to Grid</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<div class="flex-1" />

		<!-- Right Side Actions -->
		<TooltipProvider v-if="overlayId">
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="ghost" size="icon" @click="copyOverlayLink">
						<ExternalLink class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('overlaysRegistry.copyLink') || 'Copy Overlay Link' }}</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>

		<TooltipProvider>
			<Tooltip>
				<TooltipTrigger as-child>
					<Button variant="default" size="icon" @click="emit('save')">
						<Save class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>{{ t('sharedButtons.save') || 'Save' }} (Ctrl+S)</p>
				</TooltipContent>
			</Tooltip>
		</TooltipProvider>
	</div>
</template>
