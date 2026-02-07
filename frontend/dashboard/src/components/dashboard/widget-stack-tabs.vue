<script setup lang="ts">
import { X } from "lucide-vue-next";
import type { WidgetItem } from "./widgets.ts";

interface Props {
	stack: WidgetItem[];
	activeWidgetId: string | number;
	stackId: string;
}

defineProps<Props>();

const emit = defineEmits<{
	(e: "tabChange", widgetId: string | number): void;
	(e: "unstack", widgetId: string | number): void;
}>();

function getWidgetLabel(widgetId: string | number): string {
	const labels: Record<string, string> = {
		chat: "Chat",
		stream: "Stream",
		events: "Events",
		"audit-logs": "Audit Logs",
	};
	return labels[String(widgetId)] || String(widgetId);
}
</script>

<template>
	<div class="absolute bottom-0 left-0 right-0 z-50 flex items-center gap-1 px-2 pb-1">
		<div
			v-for="widget in stack"
			:key="widget.i"
			class="flex items-center gap-1 px-3 py-1 rounded-t-md text-sm cursor-pointer transition-colors"
			:class="[
				widget.i === activeWidgetId
					? 'bg-accent text-accent-foreground'
					: 'bg-muted text-muted-foreground hover:bg-muted/80',
			]"
			@click="emit('tabChange', widget.i)"
		>
			<span>{{ getWidgetLabel(widget.i) }}</span>
			<button
				class="ml-1 hover:text-destructive transition-colors"
				@click.stop="emit('unstack', widget.i)"
			>
				<X class="size-3" />
			</button>
		</div>
	</div>
</template>
