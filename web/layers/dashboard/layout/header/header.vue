<script setup lang="ts">
import { useIntervalFn, useLocalStorage, useMediaQuery } from '@vueuse/core'
import { intervalToDuration } from 'date-fns'
import { computed, onBeforeUnmount, ref } from 'vue'
import { useRealtimeDashboardStats } from '~~/layers/dashboard/api/dashboard'
import CommandMenu from '~~/layers/dashboard/components/command-menu/CommandMenu.vue'
import { padTo2Digits } from '~~/layers/dashboard/helpers/convertMillisToTime'

import Popover from '@/components/ui/popover/Popover.vue'
import PopoverContent from '@/components/ui/popover/PopoverContent.vue'
import PopoverTrigger from '@/components/ui/popover/PopoverTrigger.vue'

import StreamInfoEditor from '../stream-info-editor.vue'
import HeaderBotStatus from './header-bot-status.vue'
import HeaderProfile from './header-profile.vue'

const { stats } = useRealtimeDashboardStats()

// Mobile detection
const isDesktop = useMediaQuery('(min-width: 768px)')

const currentTime = ref(new Date())
const { pause: pauseUptimeInterval } = useIntervalFn(() => {
	currentTime.value = new Date()
}, 1000)

const uptime = computed(() => {
	if (!stats.value?.startedAt) return '00:00:00'

	const duration = intervalToDuration({
		start: new Date(stats.value.startedAt),
		end: currentTime.value,
	})

	const mappedDuration = [duration.hours ?? 0, duration.minutes ?? 0, duration.seconds ?? 0]
	if (duration.days !== undefined && duration.days !== 0) mappedDuration.unshift(duration.days)

	return mappedDuration
		.map((v) => padTo2Digits(v!))
		.filter((v) => typeof v !== 'undefined')
		.join(':')
})

onBeforeUnmount(() => {
	pauseUptimeInterval()
})

const { t } = useI18n()

const streamInfoEditorOpen = ref(false)

function openInfoEditor() {
	streamInfoEditorOpen.value = true
}

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
			return uptime.value
		case 'viewers':
			return stats.value?.viewers ?? 0
		case 'followers':
			return stats.value?.followers ?? 0
		case 'messages':
			return stats.value?.chatMessages ?? 0
		case 'subs':
			return stats.value?.subs ?? 0
		case 'usedEmotes':
			return stats.value?.usedEmotes ?? 0
		case 'requestedSongs':
			return stats.value?.requestedSongs ?? 0
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
	<div
		class="bg-card border-b-border flex w-full flex-wrap justify-between gap-2 border-b px-2 py-1"
	>
		<div class="flex flex-col flex-wrap gap-2 py-1 md:flex-row">
			<!-- Mobile search icon -->
			<CommandMenu
				v-if="!isDesktop"
				:icon-only="true"
			/>

			<!-- Stream info widget -->
			<div
				v-if="isDesktop"
				class="header-widget header-widget-stream cursor-pointer"
				@click="openInfoEditor"
			>
				<div class="header-widget-content">
					<div class="flex items-center gap-2">
						<div class="flex min-w-0 flex-1 flex-col">
							<p class="header-widget-value truncate">
								{{ stats?.title ?? 'No title' }}
							</p>
							<p class="header-widget-label truncate">
								{{ stats?.categoryName ?? 'No category' }}
							</p>
						</div>
						<Icon
							name="tabler:edit"
							class="text-muted-foreground h-3.5 w-3.5 flex-shrink-0"
						/>
					</div>
				</div>
			</div>

			<!-- Stats widgets -->
			<template v-if="isDesktop">
				<div
					v-for="widget in visibleWidgets"
					:key="widget.id"
					class="header-widget"
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

					<!-- Widget content -->
					<div class="header-widget-content">
						<p class="header-widget-value">
							{{ getWidgetValue(widget.id) }}
						</p>
						<p class="header-widget-label">
							{{ t(`dashboard.statsWidgets.${widget.id}`) }}
						</p>
					</div>
				</div>
			</template>

			<!-- Edit button -->
			<button
				v-if="isDesktop"
				class="text-muted-foreground hover:text-foreground flex items-center gap-1.5 rounded-lg border border-transparent px-3 py-2 text-xs transition-all hover:border-white/10 hover:bg-white/5"
				:class="{ 'text-foreground border-white/10 bg-white/5': isEditMode }"
				@click="toggleEditMode"
			>
				<Icon
					name="lucide:edit3"
					:size="14"
				/>
				<span>{{ isEditMode ? t('sharedButtons.close') : t('sharedButtons.edit') }}</span>
			</button>

			<!-- Add widget button -->
			<Popover v-if="isDesktop && isEditMode && hiddenWidgets.length > 0">
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
		</div>

		<div class="flex-end ml-auto flex flex-wrap items-center justify-end gap-2">
			<CommandMenu v-if="isDesktop" />
			<HeaderBotStatus />
			<HeaderProfile />
		</div>
	</div>

	<StreamInfoEditor
		v-model:open="streamInfoEditorOpen"
		:title="stats?.title"
		:category-id="stats?.categoryId"
		:category-name="stats?.categoryName"
	/>
</template>

<style scoped>
.header-widget {
	position: relative;
	display: block;
	background-color: oklch(1 0 0 / 5%);
	border-radius: var(--radius-lg);
	padding: 0.5rem 0.75rem;
	transition: all 150ms;
	border: 1px solid oklch(1 0 0 / 10%);
}

.header-widget-stream {
	max-width: 300px;
}

.header-widget[draggable='true'] {
	cursor: grab;
}

.header-widget[draggable='true']:active {
	cursor: grabbing;
}

.header-widget:hover {
	background-color: oklch(1 0 0 / 10%);
}

.header-widget:hover .hover-show {
	opacity: 1;
}

.header-widget-content {
	display: flex;
	flex-direction: column;
	gap: 0.125rem;
}

.header-widget-value {
	font-size: 1rem;
	line-height: 1.25rem;
	font-weight: 600;
	color: var(--color-foreground);
	line-height: 1.2;
}

.header-widget-label {
	font-size: 0.6875rem;
	line-height: 1rem;
	color: var(--color-muted-foreground);
}
</style>
