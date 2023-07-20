<script setup lang="ts">
import { IconSettings } from '@tabler/icons-vue';
import { useThrottleFn } from '@vueuse/core';
import { NText, NButton, NTooltip, NTag, NRow, NSpace, NSwitch, NPopconfirm } from 'naive-ui';

import { EVENTS } from './events.js';
import { OPERATIONS, Operation } from './operations.js';
import { EditableEvent } from './types.js';

import { useEventsManager } from '@/api/index.js';
import Card from '@/components/card/card.vue';

const props = defineProps<{
	event: EditableEvent
}>();

defineEmits<{
	openSettings: [id: string]
}>();

const eventsManager = useEventsManager();
const eventsPatcher = eventsManager.patch!;
const eventsDeleter = eventsManager.deleteOne;

const getEventName = (eventType: string) => EVENTS[eventType]?.name ?? eventType;

const flatOperations = Object.entries(OPERATIONS).reduce((acc, curr) => {
	if (curr[1].type === 'group' && curr[1].childrens) {
		Object.entries(curr[1].childrens).forEach(([key, value]) => acc[key] = value);
		return acc;
	}

	acc[curr[0]] = curr[1];
	return acc;
}, {} as Record<string, Operation>);

console.log(flatOperations);

const getOperation = (operationType: string) => {
	return flatOperations[operationType] ?? null;
};

const throttledSwitchState = useThrottleFn((v: boolean) => {
	eventsPatcher.mutate({ id: props.event.id!, enabled: v });
}, 500);
</script>

<template>
	<card :icon="EVENTS[event.type]?.icon" :title="getEventName(event.type)">
		<template #headerExtra>
			<n-switch
				:value="event.enabled"
				@update-value="(v) => throttledSwitchState(v)"
			/>
		</template>

		<template #content>
			<n-space vertical>
				<n-text>{{ event.description }}</n-text>
				<n-row style="gap: 8px;">
					<n-tooltip v-for="(operation, index) of event.operations" :key="index" :disabled="!operation.input">
						<template #trigger>
							<n-tag
								:disabled="!operation.enabled"
								:bordered="false"
								:type="getOperation(operation.type)?.color ?? 'info'"
							>
								{{ getOperation(operation.type).description }}
							</n-tag>
						</template>
						<n-space vertical>
							<n-text>{{ operation.input }}</n-text>
							<n-text>Delay: {{ operation.delay }} | Repeat: {{ operation.repeat }}</n-text>
						</n-space>
					</n-tooltip>
				</n-row>
			</n-space>
		</template>

		<template #footer>
			<n-button secondary size="large" @click="$emit('openSettings', event.id)">
				<span>Settings</span>
				<IconSettings />
			</n-button>
			<n-popconfirm @positive-click="eventsDeleter.mutateAsync({ id: event.id! })">
				<template #trigger>
					<n-button secondary type="error" size="large">
						<span>Delete</span>
					</n-button>
				</template>
				Are you sure?
			</n-popconfirm>
		</template>
	</card>
</template>
