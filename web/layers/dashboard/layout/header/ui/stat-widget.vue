<script setup lang="ts">
interface Props {
	label?: string
	value?: string | number
	interactive?: boolean
}

defineProps<Props>()
</script>

<template>
	<div
		class="header-widget"
		:class="{ 'cursor-pointer': interactive }"
	>
		<slot>
			<div class="header-widget-content">
				<p class="header-widget-value">
					{{ value }}
				</p>
				<p class="header-widget-label">
					{{ label }}
				</p>
			</div>
		</slot>
	</div>
</template>

<!-- Extracted 1:1 from header.vue scoped styles; kept global so slotted
     content in layout components can reuse the same class names. -->
<style>
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
