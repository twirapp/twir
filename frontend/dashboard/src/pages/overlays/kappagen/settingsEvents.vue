<script setup lang="ts">
import { TwirEventType } from '@twir/api/messages/events/events';
import { NCheckbox, NGrid, NGridItem, NSwitch, NTabPane, NTabs } from 'naive-ui';
import { watch } from 'vue';

import { useKappagenFormSettings } from './store.js';

import { flatEvents } from '@/components/events/helpers.js';

const availableEvents = Object.values(flatEvents)
	.filter(e => e.enumValue !== undefined && TwirEventType[e.enumValue])
	.map(e => {
		return {
			name: e.name,
			value: e.enumValue,
		};
	}) as Array<{ name: string, value: TwirEventType }>;
const { settings: formValue } = useKappagenFormSettings();

watch(formValue.value.events, (v) => {
	for (const event of availableEvents) {
		const exists = v.find(e => e.event === event.value);
		if (exists) continue;

		formValue.value.events.push({
			event: event.value,
			disabledStyles: [],
			enabled: true,
		});
	}
}, { immediate: true });
</script>

<template>
	<n-tabs type="line" placement="left">
		<n-tab-pane
			v-for="(event) of formValue.events" :key="event.event" :name="event.event"
			:tab="availableEvents.find(e => e.value === event.event)?.name"
		>
			<template #tab>
				<div class="flex justify-between w-full gap-3">
					<span>
						{{ availableEvents.find(e => e.value === event.event)?.name }}
					</span>
					<n-switch v-model:value="event.enabled" />
				</div>
			</template>

			<n-grid :cols="2" :x-gap="8" :y-gap="8" responsive="self">
				<n-grid-item v-for="animation of formValue.animations" :key="animation.style" :span="1">
					<n-checkbox
						:checked="!event.disabledStyles.includes(animation.style)"
						@update:checked="(checked: boolean) => {
							if (checked) event.disabledStyles = event.disabledStyles.filter(s => s !== animation.style)
							else event.disabledStyles.push(animation.style)
						}"
					>
						{{ animation.style }}
					</n-checkbox>
				</n-grid-item>
			</n-grid>
		</n-tab-pane>
	</n-tabs>
</template>

<style scoped>
:deep(.n-tabs-tab__label) {
	width: 100%;
}
</style>
