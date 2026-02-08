import { computed, ref, type ComputedRef, type Ref, type WritableComputedRef } from 'vue';
import type { WidgetItem } from './widgets.ts';

export function useWidgetStacks(
	widgets: Ref<WidgetItem[]> | ComputedRef<WidgetItem[]> | WritableComputedRef<WidgetItem[]>,
) {
	const activeStackTabs = ref<Record<string, string>>({});

	// Group widgets by stackId
	const widgetStacks = computed(() => {
		const stacks = new Map<string, WidgetItem[]>();
		const standalone: WidgetItem[] = [];

		// Only process visible widgets
		widgets.value
			.filter((w) => w.visible)
			.forEach((widget) => {
				if (widget.stackId) {
					const stack = stacks.get(widget.stackId) || [];
					stack.push(widget);
					stacks.set(widget.stackId, stack);
				} else {
					standalone.push(widget);
				}
			});

		// Sort each stack by stackOrder
		stacks.forEach((stack) => {
			stack.sort((a, b) => a.stackOrder - b.stackOrder);
		});

		return { stacks, standalone };
	});

	// Get visible item for grid (either standalone or representative of stack)
	// Returns only one item per stack - the first widget as representative
	const gridItems = computed(() => {
		const items: WidgetItem[] = [...widgetStacks.value.standalone];

		widgetStacks.value.stacks.forEach((stack, stackId) => {
			// Get the first widget in stack as reference/representative
			const referenceWidget = stack[0];
			if (!referenceWidget) return;

			// Get active tab widget id for reference
			const activeTabId = activeStackTabs.value[stackId] || stack[0]?.i;

			// Use the first widget as representative with consistent position
			const representativeWidget = {
				...referenceWidget,
				// Mark which widget is active in the stack
				_activeTabId: activeTabId,
			};
			items.push(representativeWidget);
		});

		return items;
	});

	// Get all widgets in a stack by stackId
	function getStackWidgets(stackId: string): WidgetItem[] {
		return widgetStacks.value.stacks.get(stackId) || [];
	}

	// Check if a widget is the active tab in its stack
	function isActiveInStack(widgetId: string | number): boolean {
		const widget = widgets.value.find((w) => w.i === widgetId);
		if (!widget?.stackId) return true; // Standalone widgets are always "active"

		const activeTabId = activeStackTabs.value[widget.stackId];
		// If no active tab set, first widget (stackOrder 0) is active
		if (!activeTabId) {
			const stack = widgetStacks.value.stacks.get(widget.stackId);
			return stack?.[0]?.i === widgetId;
		}
		return activeTabId === String(widgetId);
	}

	// Get stack by widgetId
	function getStackByWidgetId(widgetId: string | number): WidgetItem[] | null {
		const widget = widgets.value.find((w) => w.i === widgetId);
		if (!widget?.stackId) return null;
		return widgetStacks.value.stacks.get(widget.stackId) || null;
	}

	// Set active tab in stack
	function setActiveTab(stackId: string, widgetId: string) {
		activeStackTabs.value[stackId] = widgetId;
	}

	// Create stack from two widgets or add widget to existing stack
	function createStack(targetWidgetId: string | number, draggedWidgetId: string | number) {
		const targetWidget = widgets.value.find((w) => w.i === targetWidgetId);
		const draggedWidget = widgets.value.find((w) => w.i === draggedWidgetId);

		if (!targetWidget || !draggedWidget) return;

		// If target already has a stack, add dragged widget to it
		if (targetWidget.stackId) {
			const existingStack = widgetStacks.value.stacks.get(targetWidget.stackId);
			if (existingStack) {
				// Get max stackOrder in existing stack
				const maxOrder = Math.max(...existingStack.map((w) => w.stackOrder));

				// Add dragged widget to existing stack
				draggedWidget.x = targetWidget.x;
				draggedWidget.y = targetWidget.y;
				draggedWidget.w = targetWidget.w;
				draggedWidget.h = targetWidget.h;
				draggedWidget.stackId = targetWidget.stackId;
				draggedWidget.stackOrder = maxOrder + 1;
				return;
			}
		}

		// Create new stack
		const stackId = crypto.randomUUID();

		// Both widgets take position and size of target widget
		const targetX = targetWidget.x;
		const targetY = targetWidget.y;
		const targetW = targetWidget.w;
		const targetH = targetWidget.h;

		targetWidget.stackId = stackId;
		targetWidget.stackOrder = 0;

		draggedWidget.x = targetX;
		draggedWidget.y = targetY;
		draggedWidget.w = targetW;
		draggedWidget.h = targetH;
		draggedWidget.stackId = stackId;
		draggedWidget.stackOrder = 1;

		// Set first widget as active
		activeStackTabs.value[stackId] = String(targetWidget.i);
	}

	// Remove widget from stack
	function unstackWidget(widgetId: string | number) {
		const widget = widgets.value.find((w) => w.i === widgetId);
		if (!widget?.stackId) return;

		const stackId = widget.stackId;
		const stack = widgetStacks.value.stacks.get(stackId);
		if (!stack) return;

		widget.stackId = undefined;
		widget.stackOrder = 0;

		// If stack has only one widget left, remove its stackId too
		const remainingInStack = widgets.value.filter((w) => w.stackId === stackId);
		if (remainingInStack.length === 1) {
			remainingInStack[0].stackId = undefined;
			remainingInStack[0].stackOrder = 0;
			delete activeStackTabs.value[stackId];
		}
	}

	// Update all widgets in stack when one is moved/resized
	function syncStackPositionAndSize(
		widgetId: string | number,
		x: number,
		y: number,
		width: number,
		height: number,
	) {
		const widget = widgets.value.find((w) => w.i === widgetId);
		if (!widget?.stackId) {
			// If not in stack, just update this widget
			widget && (widget.x = x);
			widget && (widget.y = y);
			widget && (widget.w = width);
			widget && (widget.h = height);
			return;
		}

		// Update all widgets in the stack
		widgets.value.forEach((w) => {
			if (w.stackId === widget.stackId) {
				w.x = x;
				w.y = y;
				w.w = width;
				w.h = height;
			}
		});
	}

	// Check if widget can be stacked (overlaps with another widget)
	function canStackWith(
		draggedId: string | number,
		targetX: number,
		targetY: number,
	): WidgetItem | null {
		const dragged = widgets.value.find((w) => w.i === draggedId);
		if (!dragged) return null;

		// Find widget at target position
		for (const widget of gridItems.value) {
			if (widget.i === draggedId) continue;

			const overlap = !(
				targetX + dragged.w <= widget.x ||
				targetX >= widget.x + widget.w ||
				targetY + dragged.h <= widget.y ||
				targetY >= widget.y + widget.h
			);

			if (overlap) {
				return widget;
			}
		}

		return null;
	}

	return {
		widgetStacks,
		gridItems,
		activeStackTabs,
		getStackByWidgetId,
		getStackWidgets,
		isActiveInStack,
		setActiveTab,
		createStack,
		unstackWidget,
		syncStackPositionAndSize,
		canStackWith,
	};
}
