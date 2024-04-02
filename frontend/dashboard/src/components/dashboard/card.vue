<script lang="ts" setup>
import { IconGripVertical, IconEyeOff } from '@tabler/icons-vue';
import { NCard, NButton } from 'naive-ui';
import { type CSSProperties, useAttrs } from 'vue';
import { useI18n } from 'vue-i18n';

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

const { t } = useI18n();
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
			<div class="widgets-draggable-handle flex items-center">
				<IconGripVertical class="w-5 h-5" />
				{{ t(`dashboard.widgets.${attrs.item.i}.title`) }}
			</div>
		</template>

		<template #header-extra>
			<div class="flex gap-1">
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
