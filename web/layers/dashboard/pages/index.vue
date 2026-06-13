<script setup lang="ts">
import { GridItem, GridLayout } from "grid-layout-plus";
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";

import AuditLogs from "~~/layers/dashboard/components/dashboard/audit-logs.vue";
import Chat from "~~/layers/dashboard/components/dashboard/chat.vue";
import CustomWidget from "~~/layers/dashboard/components/dashboard/custom-widget.vue";
import Events from "~~/layers/dashboard/components/dashboard/events.vue";
import Stream from "~~/layers/dashboard/components/dashboard/stream.vue";
import WidgetStackTabs from "~~/layers/dashboard/components/dashboard/widget-stack-tabs.vue";
import { useWidgetStacks } from "~~/layers/dashboard/components/dashboard/widget-stacks.js";
import { type WidgetItem, useWidgets } from "~~/layers/dashboard/components/dashboard/widgets.js";
import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { useIsMobile } from "~~/layers/dashboard/composables/use-is-mobile";
import { useDashboardWidgetsCreateCustom } from "~~/layers/dashboard/api/dashboard-widgets-layout.js";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { z } from "zod";
import { toast } from "vue-sonner";

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

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

const layoutItems = ref<WidgetItem[]>([]);

let isSyncingFromGrid = false;

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

watch(
	layoutItems,
	(newLayout) => {
		if (isSyncingFromGrid) return;

		for (const layoutItem of newLayout) {
			const widget = widgets.value.find((w) => w.i === layoutItem.i);
			if (widget) {
				if (
					widget.x !== layoutItem.x ||
					widget.y !== layoutItem.y ||
					widget.w !== layoutItem.w ||
					widget.h !== layoutItem.h
				) {
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

const isCreateWidgetDialogOpen = ref(false);
const createMutation = useDashboardWidgetsCreateCustom();

const formSchema = z.object({
	name: z.string().min(2, "Name must be at least 2 characters."),
	url: z.string().url("Must be a valid URL"),
});

const { handleSubmit, resetForm } = useForm({
	validationSchema: toTypedSchema(formSchema),
});

const onSubmitWidget = handleSubmit(async (values) => {
	const visibleWidgets = widgets.value.filter((w) => w.visible);

	let newY = 0;
	let newX = 0;
	const newW = 4;
	const newH = 8;

	let found = false;
	for (let y = 0; y < 100 && !found; y++) {
		for (let x = 0; x <= 12 - newW && !found; x++) {
			const overlaps = visibleWidgets.some((widget) => {
				return !(
					x + newW <= widget.x ||
					x >= widget.x + widget.w ||
					y + newH <= widget.y ||
					y >= widget.y + widget.h
				);
			});

			if (!overlaps) {
				newX = x;
				newY = y;
				found = true;
			}
		}
	}

	const result = await createMutation.executeMutation({
		input: {
			name: values.name,
			url: values.url,
			x: newX,
			y: newY,
			w: newW,
			h: newH,
		},
	});

	if (result.error) {
		toast.error("Failed to create widget", {
			description: result.error.message,
		});
	} else {
		toast.success("Widget added to dashboard");
		resetForm();
		isCreateWidgetDialogOpen.value = false;
	}
});

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

const draggedWidgetId = ref<string | number | null>(null);
const dropTargetWidgetId = ref<string | number | null>(null);

function onMouseUp() {
	showEmptyItem.value = false;
}

function onDragStart(widgetId: string | number) {
	draggedWidgetId.value = widgetId;
}

function onDragEnd(widgetId: string | number) {
	draggedWidgetId.value = null;

	if (dropTargetWidgetId.value && dropTargetWidgetId.value !== widgetId) {
		createStack(dropTargetWidgetId.value, widgetId);
	}

	dropTargetWidgetId.value = null;
}

function onResized(_widgetId: string | number) {
}

function onMove(widgetId: string | number, x: number, y: number) {
	if (draggedWidgetId.value === widgetId) {
		const dragged = widgets.value.find((w) => w.i === widgetId);
		if (!dragged) return;

		let foundTarget: string | number | null = null;
		for (const item of layoutItems.value) {
			if (item.i === widgetId) continue;
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

function getWidgetStack(widgetId: string | number) {
	return getStackByWidgetId(widgetId);
}

function getActiveWidgetInStack(stackId: string): string | number {
	const stack = getStackWidgets(stackId);
	if (!stack.length) return "";

	for (const widget of stack) {
		if (isActiveInStack(widget.i)) {
			return widget.i;
		}
	}
	return stack[0]?.i ?? "";
}

function handleTabChange(stackId: string, widgetId: string | number) {
	setActiveTab(stackId, String(widgetId));
}

function handleUnstack(widgetId: string | number) {
	unstackWidget(widgetId);
}

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
				@resized="(i) => onResized(i)"
				@mousedown="onDragStart(item.i)"
			>
				<div
					v-if="isDropTarget(item.i)"
					class="absolute inset-0 z-40 border-2 border-dashed border-primary bg-primary/10 rounded-lg flex items-center justify-center pointer-events-none"
				>
					<div class="flex items-center gap-2 text-primary font-medium">
						<Icon name="lucide:layers" class="size-5" />
						<span>Stack widgets</span>
					</div>
				</div>

				<div v-if="showEmptyItem" class="w-full h-full absolute z-50"></div>

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
						<CustomWidget
							v-if="String(stackWidget.i).startsWith('custom-') && stackWidget.customUrl"
							v-show="isActiveInStack(stackWidget.i)"
							:item="stackWidget"
							:url="stackWidget.customUrl"
							class="h-full"
						/>
					</template>
				</template>

				<template v-else>
					<Chat v-if="item.i === 'chat'" :item="item" class="h-full" />
					<Stream v-if="item.i === 'stream'" :item="item" class="h-full" />
					<Events v-if="item.i === 'events'" :item="item" class="h-full" />
					<AuditLogs v-if="item.i === 'audit-logs'" :item="item" class="h-full" />
					<CustomWidget
						v-if="String(item.i).startsWith('custom-') && item.customUrl"
						:item="item"
						:url="item.customUrl"
						class="h-full"
					/>
				</template>

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

		<div class="fixed right-8 bottom-8 z-50" :class="[{ 'right-24!': isMobile }]">
			<DropdownMenu>
				<DropdownMenuTrigger as-child>
					<Button variant="secondary" class="h-14 w-14" size="icon">
						<Icon name="lucide:square-pen" class="size-8" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem @click="isCreateWidgetDialogOpen = true">
						<Icon name="lucide:plus" class="h-4 w-4 mr-2" />
						Create Custom Widget
					</DropdownMenuItem>
					<DropdownMenuSeparator v-if="invisibleWidgets.length" />
					<DropdownMenuItem
						v-for="widget in invisibleWidgets"
						:key="widget.i"
						@click="addWidget(widget.i)"
					>
						{{ widget.displayName || String(widget.i) }}
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>

		<Dialog v-model:open="isCreateWidgetDialogOpen">
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Create Custom Widget</DialogTitle>
					<DialogDescription>
						Add a new custom widget to your dashboard by providing a name and URL.
					</DialogDescription>
				</DialogHeader>

				<form @submit="onSubmitWidget" class="space-y-4">
					<FormField v-slot="{ componentField }" name="name">
						<FormItem>
							<FormLabel>Widget Name</FormLabel>
							<FormControl>
								<Input v-bind="componentField" placeholder="My Custom Widget" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="url">
						<FormItem>
							<FormLabel>Website URL</FormLabel>
							<FormControl>
								<Input v-bind="componentField" placeholder="https://example.com" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<div class="flex justify-end gap-2">
						<Button type="button" variant="outline" @click="isCreateWidgetDialogOpen = false">
							Cancel
						</Button>
						<Button type="submit" :disabled="createMutation.fetching.value">
							{{ createMutation.fetching.value ? "Creating..." : "Create Widget" }}
						</Button>
					</div>
				</form>
			</DialogContent>
		</Dialog>
	</div>
</template>

<style scoped>
@reference '~/assets/css/tailwind.css';

.vgl-layout {
	@apply w-full;
}

:deep(.vgl-item__resizer) {
	z-index: 51;
}
</style>
