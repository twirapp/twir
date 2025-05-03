<script setup lang="ts">
import { TwirEventType } from '@twir/api/messages/events/events'
import { NCheckbox, NGrid, NGridItem, NSwitch, NTabPane, NTabs } from 'naive-ui'
import { watch } from 'vue'

import { useKappagenFormSettings } from './store.js'

import { flatEvents } from '@/features/events/constants/helpers.js'

const availableEvents = Object.values(flatEvents)
	.filter(e => e.enumValue !== undefined && TwirEventType[e.enumValue])
	.map(e => {
		return {
			name: e.name,
			value: e.enumValue,
		}
	}) as Array<{ name: string, value: TwirEventType }>
const { settings: formValue } = useKappagenFormSettings()

watch(formValue.value.events, (v) => {
	for (const event of availableEvents) {
		const exists = v.find(e => e.event === event.value)
		if (exists) continue

		formValue.value.events.push({
			event: event.value,
			disabledStyles: [],
			enabled: true,
		})
	}
}, { immediate: true })
</script>

<template>
	<NTabs type="line" placement="left">
		<NTabPane
			v-for="(event) of formValue.events" :key="event.event" :name="event.event"
			:tab="availableEvents.find(e => e.value === event.event)?.name"
		>
			<template #tab>
				<div class="flex justify-between w-full gap-3">
					<span>
						{{ availableEvents.find(e => e.value === event.event)?.name }}
					</span>
					<NSwitch v-model:value="event.enabled" />
				</div>
			</template>

			<NGrid :cols="2" :x-gap="8" :y-gap="8" responsive="self">
				<NGridItem v-for="animation of formValue.animations" :key="animation.style" :span="1">
					<NCheckbox
						:checked="!event.disabledStyles.includes(animation.style)"
						@update:checked="(checked: boolean) => {
							if (checked) event.disabledStyles = event.disabledStyles.filter(s => s !== animation.style)
							else event.disabledStyles.push(animation.style)
						}"
					>
						{{ animation.style }}
					</NCheckbox>
				</NGridItem>
			</NGrid>
		</NTabPane>
	</NTabs>
</template>

<style scoped>
:deep(.n-tabs-tab__label) {
	width: 100%;
}
</style>
