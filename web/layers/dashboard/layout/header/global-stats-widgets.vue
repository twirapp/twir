<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import { computed, ref } from 'vue'

import Popover from '@/components/ui/popover/Popover.vue'
import PopoverContent from '@/components/ui/popover/PopoverContent.vue'
import PopoverTrigger from '@/components/ui/popover/PopoverTrigger.vue'

import StatWidget from './ui/stat-widget.vue'

import type { DashboardStats } from './composables/use-platform-stats.js'

interface Props {
	stats: DashboardStats | undefined
	uptime: string
}

const props = defineProps<Props>()

const { t } = useI18n()

// Widget management
type WidgetType =
	| 'uptime'
	| 'viewers'
	| 'followers'
	| 'messages'
	| 'subs'
	| 'usedEmotes'
	| 'requestedSongs'

interface WidgetConfig {
	id: WidgetType
	enabled: boolean
	order: number
}

const defaultWidgets: WidgetConfig[] = [
	{ id: 'uptime', enabled: true, order: 0 },
	{ id: 'viewers', enabled: true, order: 1 },
	{ id: 'followers', enabled: true, order: 2 },
	{ id: 'messages', enabled: true, order: 3 },
	{ id: 'subs', enabled: true, order: 4 },
	{ id: 'usedEmotes', enabled: true, order: 5 },
	{ id: 'requestedSongs', enabled: true, order: 6 },
]

const widgetsConfig = useLocalStorage<WidgetConfig[]>('twirHeaderStatsWidgetsv1', defaultWidgets)
const isEditMode = ref(false)

const visibleWidgets = computed(() => {
	return widgetsConfig.value.filter((w) => w.enabled).sort((a, b) => a.order - b.order)
})

const hiddenWidgets = computed(() => {
	return widgetsConfig.value.filter((w) => !w.enabled)
})

function toggleEditMode() {
	isEditMode.value = !isEditMode.value
}

function removeWidget(widgetId: WidgetType) {
	const widget = widgetsConfig.value.find((w) => w.id === widgetId)
	if (widget) {
		widget.enabled = false
	}
}

function addWidget(widgetId: WidgetType) {
	const widget = widgetsConfig.value.find((w) => w.id === widgetId)
	if (widget) {
		widget.enabled = true
		// Set order to be last
		const maxOrder = Math.max(...widgetsConfig.value.map((w) => w.order), -1)
		widget.order = maxOrder + 1
	}
}

function getWidgetValue(widgetId: WidgetType): string | number {
	switch (widgetId) {
		case 'uptime':
			return props.uptime
		case 'viewers':
			return props.stats?.viewers ?? 0
		case 'followers':
			return props.stats?.followers ?? 0
		case 'messages':
			return props.stats?.chatMessages ?? 0
		case 'subs':
			return props.stats?.subs ?? 0
		case 'usedEmotes':
			return props.stats?.usedEmotes ?? 0
		case 'requestedSongs':
			return props.stats?.requestedSongs ?? 0
		default:
			return 0
	}
}

// Drag & Drop functionality
const draggedWidgetId = ref<WidgetType | null>(null)
const dragOverWidgetId = ref<WidgetType | null>(null)

function onDragStart(widgetId: WidgetType, event: DragEvent) {
	draggedWidgetId.value = widgetId
	if (event.dataTransfer) {
		event.dataTransfer.effectAllowed = 'move'
		event.dataTransfer.setData('text/plain', widgetId)
	}
}

function onDragOver(widgetId: WidgetType, event: DragEvent) {
	event.preventDefault()
	if (event.dataTransfer) {
		event.dataTransfer.dropEffect = 'move'
	}
	dragOverWidgetId.value = widgetId
}

function onDragLeave() {
	dragOverWidgetId.value = null
}

function onDrop(targetWidgetId: WidgetType, event: DragEvent) {
	event.preventDefault()

	if (!draggedWidgetId.value || draggedWidgetId.value === targetWidgetId) {
		draggedWidgetId.value = null
		dragOverWidgetId.value = null
		return
	}

	const draggedWidget = widgetsConfig.value.find((w) => w.id === draggedWidgetId.value)
	const targetWidget = widgetsConfig.value.find((w) => w.id === targetWidgetId)

	if (draggedWidget && targetWidget) {
		// Swap orders
		const tempOrder = draggedWidget.order
		draggedWidget.order = targetWidget.order
		targetWidget.order = tempOrder
	}

	draggedWidgetId.value = null
	dragOverWidgetId.value = null
}

function onDragEnd() {
	draggedWidgetId.value = null
	dragOverWidgetId.value = null
}
</script>

<template>
	<!-- Stats widgets -->
	<template v-for="widget in visibleWidgets" :key="widget.id">
		<StatWidget
			:class="{
				'pl-9': isEditMode,
				'opacity-50': draggedWidgetId === widget.id,
				'ring-primary ring-2': dragOverWidgetId === widget.id && draggedWidgetId !== widget.id,
			}"
			:draggable="isEditMode"
			@dragstart="onDragStart(widget.id, $event)"
			@dragover="onDragOver(widget.id, $event)"
			@dragleave="onDragLeave"
			@drop="onDrop(widget.id, $event)"
			@dragend="onDragEnd"
		>
			<!-- Edit mode: Grip icon -->
			<div
				v-if="isEditMode"
				class="text-muted-foreground/70 absolute top-1/2 left-2 -translate-y-1/2 cursor-grab active:cursor-grabbing"
				@mousedown.stop
			>
				<Icon
					name="lucide:grip-vertical"
					:size="16"
					:stroke-width="1.5"
				/>
			</div>

			<!-- Edit mode: Remove button -->
			<button
				v-if="isEditMode"
				class="hover-show text-muted-foreground/50 absolute top-1.5 right-1.5 opacity-0 transition-colors hover:text-red-400"
				@click.stop="removeWidget(widget.id)"
			>
				<Icon
					name="lucide:x"
					:size="14"
				/>
			</button>

			<!-- Widget content; layouts can override per-widget rendering -->
			<slot :name="widget.id" :widget="widget">
				<div class="header-widget-content">
					<p class="header-widget-value">
						{{ getWidgetValue(widget.id) }}
					</p>
					<p class="header-widget-label">
						{{ t(`dashboard.statsWidgets.${widget.id}`) }}
					</p>
				</div>
			</slot>
		</StatWidget>
	</template>

	<!-- Edit button -->
	<button
		class="text-muted-foreground hover:text-foreground flex items-center gap-1.5 rounded-lg border border-transparent px-3 py-2 text-xs transition-all hover:border-white/10 hover:bg-white/5"
		:class="{ 'text-foreground border-white/10 bg-white/5': isEditMode }"
		@click="toggleEditMode"
	>
		<Icon
			name="lucide:edit-3"
			:size="14"
		/>
		<span>{{ isEditMode ? t('sharedButtons.close') : t('sharedButtons.edit') }}</span>
	</button>

	<!-- Add widget button -->
	<Popover v-if="isEditMode && hiddenWidgets.length > 0">
		<PopoverTrigger as-child>
			<button
				class="text-muted-foreground hover:text-foreground flex items-center gap-1.5 rounded-lg border border-transparent px-3 py-2 text-xs transition-all hover:border-white/10 hover:bg-white/5"
			>
				<Icon
					name="lucide:plus"
					:size="14"
				/>
				<span>{{ t('sharedButtons.add') }}</span>
			</button>
		</PopoverTrigger>
		<PopoverContent
			class="w-56 p-2"
			align="start"
		>
			<div class="space-y-1">
				<p class="text-muted-foreground px-2 py-1 text-xs font-semibold">
					{{ t('dashboard.statsWidgets.addWidget') }}
				</p>
				<button
					v-for="widget in hiddenWidgets"
					:key="widget.id"
					class="text-foreground flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-xs transition-colors hover:bg-white/10"
					@click="addWidget(widget.id)"
				>
					<Icon
						name="lucide:plus"
						:size="12"
					/>
					<span>{{ t(`dashboard.statsWidgets.${widget.id}`) }}</span>
				</button>
			</div>
		</PopoverContent>
	</Popover>
</template>
