<script setup lang="ts">
import { IconEdit } from "@tabler/icons-vue";
import { useIntervalFn, useLocalStorage } from "@vueuse/core";
import { intervalToDuration } from "date-fns";
import { Edit3, GripVertical, Plus, X } from "lucide-vue-next";
import { computed, onBeforeUnmount, ref } from "vue";
import { useI18n } from "vue-i18n";

import StreamInfoEditor from "../stream-info-editor.vue";

import { useRealtimeDashboardStats } from "@/api/dashboard";
import CommandMenu from "@/components/command-menu/CommandMenu.vue";
import Popover from "@/components/ui/popover/Popover.vue";
import PopoverContent from "@/components/ui/popover/PopoverContent.vue";
import PopoverTrigger from "@/components/ui/popover/PopoverTrigger.vue";
import { padTo2Digits } from "@/helpers/convertMillisToTime";
import HeaderProfile from "@/layout/header/header-profile.vue";
import HeaderBotStatus from "@/layout/header/header-bot-status.vue";

const { stats } = useRealtimeDashboardStats();

const currentTime = ref(new Date());
const { pause: pauseUptimeInterval } = useIntervalFn(() => {
	currentTime.value = new Date();
}, 1000);

const uptime = computed(() => {
	if (!stats.value?.startedAt) return "00:00:00";

	const duration = intervalToDuration({
		start: new Date(stats.value.startedAt),
		end: currentTime.value,
	});

	const mappedDuration = [duration.hours ?? 0, duration.minutes ?? 0, duration.seconds ?? 0];
	if (duration.days !== undefined && duration.days !== 0) mappedDuration.unshift(duration.days);

	return mappedDuration
		.map((v) => padTo2Digits(v!))
		.filter((v) => typeof v !== "undefined")
		.join(":");
});

onBeforeUnmount(() => {
	pauseUptimeInterval();
});

const { t } = useI18n();

const streamInfoEditorOpen = ref(false);

function openInfoEditor() {
	streamInfoEditorOpen.value = true;
}

// Widget management
type WidgetType =
	| "uptime"
	| "viewers"
	| "followers"
	| "messages"
	| "subs"
	| "usedEmotes"
	| "requestedSongs";

interface WidgetConfig {
	id: WidgetType;
	enabled: boolean;
	order: number;
}

const defaultWidgets: WidgetConfig[] = [
	{ id: "uptime", enabled: true, order: 0 },
	{ id: "viewers", enabled: true, order: 1 },
	{ id: "followers", enabled: true, order: 2 },
	{ id: "messages", enabled: true, order: 3 },
	{ id: "subs", enabled: true, order: 4 },
	{ id: "usedEmotes", enabled: true, order: 5 },
	{ id: "requestedSongs", enabled: true, order: 6 },
];

const widgetsConfig = useLocalStorage<WidgetConfig[]>("twirHeaderStatsWidgetsv1", defaultWidgets);
const isEditMode = ref(false);

const visibleWidgets = computed(() => {
	return widgetsConfig.value.filter((w) => w.enabled).sort((a, b) => a.order - b.order);
});

const hiddenWidgets = computed(() => {
	return widgetsConfig.value.filter((w) => !w.enabled);
});

function toggleEditMode() {
	isEditMode.value = !isEditMode.value;
}

function removeWidget(widgetId: WidgetType) {
	const widget = widgetsConfig.value.find((w) => w.id === widgetId);
	if (widget) {
		widget.enabled = false;
	}
}

function addWidget(widgetId: WidgetType) {
	const widget = widgetsConfig.value.find((w) => w.id === widgetId);
	if (widget) {
		widget.enabled = true;
		// Set order to be last
		const maxOrder = Math.max(...widgetsConfig.value.map((w) => w.order), -1);
		widget.order = maxOrder + 1;
	}
}

function getWidgetValue(widgetId: WidgetType): string | number {
	switch (widgetId) {
		case "uptime":
			return uptime.value;
		case "viewers":
			return stats.value?.viewers ?? 0;
		case "followers":
			return stats.value?.followers ?? 0;
		case "messages":
			return stats.value?.chatMessages ?? 0;
		case "subs":
			return stats.value?.subs ?? 0;
		case "usedEmotes":
			return stats.value?.usedEmotes ?? 0;
		case "requestedSongs":
			return stats.value?.requestedSongs ?? 0;
		default:
			return 0;
	}
}

// Drag & Drop functionality
const draggedWidgetId = ref<WidgetType | null>(null);
const dragOverWidgetId = ref<WidgetType | null>(null);

function onDragStart(widgetId: WidgetType, event: DragEvent) {
	draggedWidgetId.value = widgetId;
	if (event.dataTransfer) {
		event.dataTransfer.effectAllowed = "move";
		event.dataTransfer.setData("text/plain", widgetId);
	}
}

function onDragOver(widgetId: WidgetType, event: DragEvent) {
	event.preventDefault();
	if (event.dataTransfer) {
		event.dataTransfer.dropEffect = "move";
	}
	dragOverWidgetId.value = widgetId;
}

function onDragLeave() {
	dragOverWidgetId.value = null;
}

function onDrop(targetWidgetId: WidgetType, event: DragEvent) {
	event.preventDefault();

	if (!draggedWidgetId.value || draggedWidgetId.value === targetWidgetId) {
		draggedWidgetId.value = null;
		dragOverWidgetId.value = null;
		return;
	}

	const draggedWidget = widgetsConfig.value.find((w) => w.id === draggedWidgetId.value);
	const targetWidget = widgetsConfig.value.find((w) => w.id === targetWidgetId);

	if (draggedWidget && targetWidget) {
		// Swap orders
		const tempOrder = draggedWidget.order;
		draggedWidget.order = targetWidget.order;
		targetWidget.order = tempOrder;
	}

	draggedWidgetId.value = null;
	dragOverWidgetId.value = null;
}

function onDragEnd() {
	draggedWidgetId.value = null;
	dragOverWidgetId.value = null;
}
</script>

<template>
	<div
		class="flex flex-wrap justify-between bg-card w-full px-2 py-1 gap-2 border-b border-b-border"
	>
		<div class="flex flex-wrap md:flex-row flex-col gap-2 py-1">
			<!-- Stream info widget -->
			<div class="header-widget cursor-pointer" @click="openInfoEditor">
				<div class="header-widget-content">
					<div class="flex items-center gap-2">
						<div class="flex flex-col flex-1">
							<p class="header-widget-value">
								{{ stats?.title ?? "No title" }}
							</p>
							<p class="header-widget-label">
								{{ stats?.categoryName ?? "No category" }}
							</p>
						</div>
						<IconEdit class="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
					</div>
				</div>
			</div>

			<!-- Stats widgets -->
			<div
				v-for="widget in visibleWidgets"
				:key="widget.id"
				class="header-widget"
				:class="{
					'pl-9': isEditMode,
					'opacity-50': draggedWidgetId === widget.id,
					'ring-2 ring-primary': dragOverWidgetId === widget.id && draggedWidgetId !== widget.id,
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
					class="absolute left-0 top-1/3 -translate-y-1/2 text-muted-foreground/70 cursor-grab active:cursor-grabbing"
					@mousedown.stop
				>
					<GripVertical :size="16" :stroke-width="1.5" />
				</div>

				<!-- Edit mode: Remove button -->
				<button
					v-if="isEditMode"
					class="hover-show absolute right-1.5 top-1.5 text-muted-foreground/50 hover:text-red-400 transition-colors opacity-0"
					@click.stop="removeWidget(widget.id)"
				>
					<X :size="14" />
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

			<!-- Edit button -->
			<button
				class="flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all border border-transparent hover:border-white/10"
				:class="{ 'bg-white/5 border-white/10 text-foreground': isEditMode }"
				@click="toggleEditMode"
			>
				<Edit3 :size="14" />
				<span>{{ isEditMode ? t("sharedButtons.close") : t("sharedButtons.edit") }}</span>
			</button>

			<!-- Add widget button -->
			<Popover v-if="isEditMode && hiddenWidgets.length > 0">
				<PopoverTrigger as-child>
					<button
						class="flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs text-muted-foreground hover:text-foreground hover:bg-white/5 transition-all border border-transparent hover:border-white/10"
					>
						<Plus :size="14" />
						<span>{{ t("sharedButtons.add") }}</span>
					</button>
				</PopoverTrigger>
				<PopoverContent class="w-56 p-2" align="start">
					<div class="space-y-1">
						<p class="text-xs font-semibold text-muted-foreground px-2 py-1">
							{{ t("dashboard.statsWidgets.addWidget") }}
						</p>
						<button
							v-for="widget in hiddenWidgets"
							:key="widget.id"
							class="w-full flex items-center gap-2 px-2 py-1.5 text-xs text-foreground hover:bg-white/10 rounded transition-colors text-left"
							@click="addWidget(widget.id)"
						>
							<Plus :size="12" />
							<span>{{ t(`dashboard.statsWidgets.${widget.id}`) }}</span>
						</button>
					</div>
				</PopoverContent>
			</Popover>
		</div>

		<div class="ml-auto flex flex-wrap justify-end gap-2 flex-end items-center">
			<CommandMenu />
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
@import "@/assets/index.css";

.header-widget {
	position: relative;
	display: block;
	background-color: oklch(1 0 0 / 5%);
	border-radius: var(--radius-lg);
	padding: 0.5rem 0.75rem;
	transition: all 150ms;
	border: 1px solid oklch(1 0 0 / 10%);
}

.header-widget[draggable="true"] {
	cursor: grab;
}

.header-widget[draggable="true"]:active {
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
