<script setup lang="ts">
import { GridItem, GridLayout } from "grid-layout-plus";
import { Layers, SquarePen } from "lucide-vue-next";
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";

import AuditLogs from "@/components/dashboard/audit-logs.vue";
import Chat from "@/components/dashboard/chat.vue";
import Events from "@/components/dashboard/events.vue";
import Stream from "@/components/dashboard/stream.vue";
import WidgetStackTabs from "@/components/dashboard/widget-stack-tabs.vue";
import { useWidgetStacks } from "@/components/dashboard/widget-stacks.ts";
import { useWidgets, type WidgetItem } from "@/components/dashboard/widgets.ts";
import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useIsMobile } from "@/composables/use-is-mobile";

const { isMobile } = useIsMobile();
const widgets = useWidgets();

const {
	gridItems,
	getStackByWidgetId,
	getStackWidgets,
	isActiveInStack,
	setActiveTab,
	createStack,
	unstackWidget,
	syncStackPositionAndSize,
} = useWidgetStacks(widgets);

// Layout items for GridLayout - only items that should be visible in grid
// This prevents grid-layout-plus from trying to separate stacked widgets
const layoutItems = ref<WidgetItem[]>([]);

// Flag to prevent infinite loop between watchers
let isSyncingFromGrid = false;

// Sync layoutItems with gridItems (gridItems → layoutItems)
// This is one-way: we copy gridItems to layoutItems for GridLayout to use
watch(
	gridItems,
	(newItems) => {
		isSyncingFromGrid = true;
		layoutItems.value = newItems.filter((v) => v.visible).map((item) => ({ ...item }));
		nextTick(() => {
			isSyncingFromGrid = false;
		});
	},
	{ immediate: true, deep: true },
);

// Sync changes from layoutItems back to widgets (layoutItems → widgets)
// This handles user interactions with GridLayout (drag, resize)
watch(
	layoutItems,
	(newLayout) => {
		// Skip if we just synced from gridItems (server data)
		if (isSyncingFromGrid) return;

		for (const layoutItem of newLayout) {
			const widget = widgets.value.find((w) => w.i === layoutItem.i);
			if (widget) {
				// Check if position/size changed
				if (
					widget.x !== layoutItem.x ||
					widget.y !== layoutItem.y ||
					widget.w !== layoutItem.w ||
					widget.h !== layoutItem.h
				) {
					// If widget is in a stack, sync all widgets in stack
					if (widget.stackId) {
						syncStackPositionAndSize(
							widget.i,
							layoutItem.x,
							layoutItem.y,
							layoutItem.w,
							layoutItem.h,
						);
					} else {
						widget.x = layoutItem.x;
						widget.y = layoutItem.y;
						widget.w = layoutItem.w;
						widget.h = layoutItem.h;
					}
				}
			}
		}
	},
	{ deep: true },
);

const invisibleWidgets = computed(() => widgets.value.filter((v) => !v.visible));

function addWidget(key: string | number) {
	const item = widgets.value.find((v) => v.i === key);
	if (!item) return;

	const widgetsLength = layoutItems.value.length;

	item.visible = true;
	item.x = (widgetsLength * 2) % 12;
	item.y = widgetsLength + 12;
	item.stackId = undefined;
	item.stackOrder = 0;
}

const showEmptyItem = ref(false);

// Drag-to-stack state
const draggedWidgetId = ref<string | number | null>(null);
const dropTargetWidgetId = ref<string | number | null>(null);

function onMouseUp() {
	showEmptyItem.value = false;
}

// Handle drag start
function onDragStart(widgetId: string | number) {
	draggedWidgetId.value = widgetId;
}

// Handle drag end - check if we should stack
function onDragEnd(widgetId: string | number) {
	draggedWidgetId.value = null;

	if (dropTargetWidgetId.value && dropTargetWidgetId.value !== widgetId) {
		// Create stack with the target widget
		createStack(dropTargetWidgetId.value, widgetId);
	}

	dropTargetWidgetId.value = null;
}

// Handle move event - detect overlap for stacking
function onMove(widgetId: string | number, x: number, y: number) {
	// Detect potential stack target during drag
	if (draggedWidgetId.value === widgetId) {
		const dragged = widgets.value.find((w) => w.i === widgetId);
		if (!dragged) return;

		// Find overlapping widget
		let foundTarget: string | number | null = null;
		for (const item of layoutItems.value) {
			if (item.i === widgetId) continue;
			// Skip if same stack
			if (dragged.stackId && item.stackId === dragged.stackId) continue;

			const overlap = !(
				x + dragged.w <= item.x ||
				x >= item.x + item.w ||
				y + dragged.h <= item.y ||
				y >= item.y + item.h
			);

			if (overlap) {
				foundTarget = item.i;
				break;
			}
		}
		dropTargetWidgetId.value = foundTarget;
	}
}

// Get stack info for a widget
function getWidgetStack(widgetId: string | number) {
	return getStackByWidgetId(widgetId);
}

// Get active widget id in a stack
function getActiveWidgetInStack(stackId: string): string | number {
	const stack = getStackWidgets(stackId);
	if (!stack.length) return "";

	// Find which widget is active
	for (const widget of stack) {
		if (isActiveInStack(widget.i)) {
			return widget.i;
		}
	}
	return stack[0]?.i ?? "";
}

// Handle tab change in stack
function handleTabChange(stackId: string, widgetId: string | number) {
	setActiveTab(stackId, String(widgetId));
}

// Handle unstack
function handleUnstack(widgetId: string | number) {
	unstackWidget(widgetId);
}

// Check if widget is a drop target for stacking
function isDropTarget(widgetId: string | number): boolean {
	return dropTargetWidgetId.value === widgetId;
}

onMounted(async () => {
	await nextTick();

	document.querySelectorAll(".vgl-item__resizer").forEach((el) => {
		el.addEventListener("mousedown", () => {
			showEmptyItem.value = true;
		});

		window.addEventListener("mouseup", onMouseUp);
	});
});

onBeforeUnmount(() => {
	window.removeEventListener("mouseup", onMouseUp);
});
</script>

<template>
	<div class="w-full h-full pl-1">
		<GridLayout v-model:layout="layoutItems" :row-height="30" :use-css-transforms="false">
			<GridItem
				v-for="item in layoutItems"
				:key="item.i"
				:x="item.x"
				:y="item.y"
				:w="item.w"
				:h="item.h"
				:i="item.i"
				:min-w="item.minW"
				:min-h="item.minH"
				drag-allow-from=".widgets-draggable-handle"
				@move="(i, x, y) => onMove(i, x, y)"
				@moved="(i) => onDragEnd(i)"
				@mousedown="onDragStart(item.i)"
			>
				<!-- Drop target indicator -->
				<div
					v-if="isDropTarget(item.i)"
					class="absolute inset-0 z-40 border-2 border-dashed border-primary bg-primary/10 rounded-lg flex items-center justify-center pointer-events-none"
				>
					<div class="flex items-center gap-2 text-primary font-medium">
						<Layers class="size-5" />
						<span>Stack widgets</span>
					</div>
				</div>

				<div v-if="showEmptyItem" class="w-full h-full absolute z-50"></div>

				<!-- Stacked widgets - render all widgets in stack, show only active -->
				<template v-if="item.stackId">
					<template v-for="stackWidget in getStackWidgets(item.stackId)" :key="stackWidget.i">
						<Chat
							v-if="stackWidget.i === 'chat'"
							v-show="isActiveInStack(stackWidget.i)"
							:item="stackWidget"
							class="h-full"
						/>
						<Stream
							v-if="stackWidget.i === 'stream'"
							v-show="isActiveInStack(stackWidget.i)"
							:item="stackWidget"
							class="h-full"
						/>
						<Events
							v-if="stackWidget.i === 'events'"
							v-show="isActiveInStack(stackWidget.i)"
							:item="stackWidget"
							class="h-full"
						/>
						<AuditLogs
							v-if="stackWidget.i === 'audit-logs'"
							v-show="isActiveInStack(stackWidget.i)"
							:item="stackWidget"
							class="h-full"
						/>
					</template>
				</template>

				<!-- Standalone widget content -->
				<template v-else>
					<Chat v-if="item.i === 'chat'" :item="item" class="h-full" />
					<Stream v-if="item.i === 'stream'" :item="item" class="h-full" />
					<Events v-if="item.i === 'events'" :item="item" class="h-full" />
					<AuditLogs v-if="item.i === 'audit-logs'" :item="item" class="h-full" />
				</template>

				<!-- Stack tabs -->
				<WidgetStackTabs
					v-if="item.stackId && getWidgetStack(item.i)"
					:stack="getWidgetStack(item.i)!"
					:active-widget-id="getActiveWidgetInStack(item.stackId)"
					:stack-id="item.stackId"
					@tab-change="(widgetId) => handleTabChange(item.stackId!, widgetId)"
					@unstack="handleUnstack"
				/>
			</GridItem>
		</GridLayout>

		<div
			v-if="invisibleWidgets.length"
			class="fixed right-8 bottom-8 z-50"
			:class="[{ 'right-24!': isMobile }]"
		>
			<DropdownMenu>
				<DropdownMenuTrigger as-child>
					<Button variant="secondary" class="h-14 w-14" size="icon">
						<SquarePen class="size-8" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem
						v-for="widget in invisibleWidgets"
						:key="widget.i"
						@click="addWidget(widget.i)"
					>
						{{ String(widget.i) }}
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>
	</div>
</template>

<style scoped>
@reference '@/assets/index.css';

.vgl-layout {
	@apply w-full;
}

:deep(.vgl-item__resizer) {
	z-index: 51;
}
</style>
