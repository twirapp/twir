<script lang="ts" setup>
import { IconGripVertical, IconEyeOff } from '@tabler/icons-vue';
import { NCard, NButton } from 'naive-ui';
import { type CSSProperties, useAttrs } from 'vue';

import { useWidgets, type WidgetItem } from './widgets.js';

defineSlots<{
	default: any,
	action?: any
	'header-extra'?: any
}>();

withDefaults(defineProps<{
	contentStyle?: CSSProperties
}>(), {
	contentStyle: () => ({ padding: '0px' }),
});

const widgets = useWidgets();

const attrs = useAttrs() as { item: WidgetItem, [x: string]: unknown };

const hideItem = () => {
	const item = widgets.value.find(item => item.i === attrs.item.i);
	if (!item) return;
	item.visible = false;
};
</script>

<template>
	<n-card
		:segmented="{
			content: true,
			footer: 'soft'
		}"
		header-style="padding: 5px;"
		:content-style="contentStyle"
		style="width: 100%; height: 100%"
		v-bind="$attrs"
	>
		<template #header>
			<div class="widgets-draggable-handle" style="display: flex; align-items: center">
				<IconGripVertical style="width: 20px; height: 20px;" />
				{{ attrs.item.i }}
			</div>
		</template>

		<template #header-extra>
			<div style="display: flex; gap: 5px">
				<slot name="header-extra" />
				<n-button text size="small" @click="hideItem">
					<IconEyeOff />
				</n-button>
			</div>
		</template>

		<slot />

		<template #action>
			<slot name="action" />
		</template>
	</n-card>
</template>
