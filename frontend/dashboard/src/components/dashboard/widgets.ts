import { computed, ref, watch } from 'vue';

import type { LayoutItem } from 'grid-layout-plus';

import {
	useDashboardWidgetsLayout,
	useDashboardWidgetsLayoutSubscription,
	useDashboardWidgetsLayoutUpdate,
} from '@/api/dashboard-widgets-layout.ts';

export type WidgetItem = LayoutItem & {
	visible: boolean;
	stackId?: string;
	stackOrder: number;
	customUrl?: string;
	displayName?: string;
};

const defaultWidgets: WidgetItem[] = [
	{
		x: 6,
		y: 0,
		w: 4,
		h: 8,
		i: 'chat',
		minW: 3,
		minH: 8,
		visible: true,
		stackOrder: 0,
	},
	{
		x: 3,
		y: 0,
		w: 3,
		h: 8,
		i: 'events',
		minW: 3,
		minH: 8,
		visible: true,
		stackOrder: 0,
	},
	{
		x: 0,
		y: 0,
		w: 3,
		h: 8,
		i: 'stream',
		minW: 3,
		minH: 8,
		visible: true,
		stackOrder: 0,
	},
	{
		x: 0,
		y: 9,
		w: 3,
		h: 8,
		i: 'audit-logs',
		minW: 3,
		minH: 8,
		visible: true,
		stackOrder: 0,
	},
];

export function useWidgets() {
	const widgets = ref<WidgetItem[]>([...defaultWidgets]);
	const { layout, fetching } = useDashboardWidgetsLayout();
	const subscription = useDashboardWidgetsLayoutSubscription();
	const updateMutation = useDashboardWidgetsLayoutUpdate();

	// Flag to prevent sending updates before initial load from server
	let isInitialized = false;
	// Flag to prevent recursion when updating from server
	let isUpdatingFromServer = false;

	// Load from server on mount
	watch(
		layout,
		(serverLayout) => {
			if (!fetching.value && serverLayout.length > 0) {
				isUpdatingFromServer = true;

				// Map server layout to widgets - all data is now in the layout including custom fields
				widgets.value = serverLayout.map((item) => ({
					x: item.x,
					y: item.y,
					w: item.w,
					h: item.h,
					i: item.widgetId,
					minW: item.minW,
					minH: item.minH,
					visible: item.visible,
					stackId: item.stackId ?? undefined,
					stackOrder: item.stackOrder,
					customUrl: item.customUrl ?? undefined,
					displayName: item.customName ?? undefined,
				}));

				// Mark as initialized after first server load
				isInitialized = true;
				// Reset flag on next tick to allow user changes
				setTimeout(() => {
					isUpdatingFromServer = false;
				}, 0);
			}
		},
		{ immediate: true },
	);

	// Listen to subscription updates
	watch(
		subscription.layout,
		(serverLayout) => {
			if (serverLayout.length > 0) {
				isUpdatingFromServer = true;

				widgets.value = serverLayout.map((item) => ({
					x: item.x,
					y: item.y,
					w: item.w,
					h: item.h,
					i: item.widgetId,
					minW: item.minW,
					minH: item.minH,
					visible: item.visible,
					stackId: item.stackId ?? undefined,
					stackOrder: item.stackOrder,
					customUrl: item.customUrl ?? undefined,
					displayName: item.customName ?? undefined,
				}));
				// Reset flag on next tick to allow user changes
				setTimeout(() => {
					isUpdatingFromServer = false;
				}, 0);
			}
		},
		{ immediate: true },
	);

	// Watch for user changes and sync to server
	let saveTimeout: ReturnType<typeof setTimeout> | null = null;
	watch(
		widgets,
		(newWidgets) => {
			// Skip if not yet initialized (waiting for server data)
			if (!isInitialized) {
				return;
			}

			// Skip if this update came from server to prevent recursion
			if (isUpdatingFromServer) {
				return;
			}

			// Debounce server updates
			if (saveTimeout) {
				clearTimeout(saveTimeout);
			}

			saveTimeout = setTimeout(() => {
				updateMutation.executeMutation({
					input: newWidgets.map((w) => ({
						widgetId: String(w.i),
						x: w.x,
						y: w.y,
						w: w.w,
						h: w.h,
						minW: w.minW ?? 3,
						minH: w.minH ?? 8,
						visible: w.visible,
						stackId: w.stackId ?? null,
						stackOrder: w.stackOrder,
					})),
				});
			}, 1000);
		},
		{ deep: true },
	);

	return computed({
		get: () => widgets.value,
		set: (val) => {
			widgets.value = val;
		},
	});
}
