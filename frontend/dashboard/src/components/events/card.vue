<script setup lang="ts">
import { useThrottleFn } from '@vueuse/core';
import { NCard, NTag, NSpace, NText, NRow, NTooltip, NSwitch } from 'naive-ui';

import { EVENTS } from './events.js';
import { OPERATIONS } from './operations.js';
import { EditableEvent } from './types.js';

import { useEventsManager } from '@/api/index.js';

const props = defineProps<{
	event: EditableEvent
}>();

const eventsManager = useEventsManager();
const eventsPatcher = eventsManager.patch!;

const getEventName = (eventType: string) => EVENTS[eventType]?.name ?? eventType;
const getOperationName = (operationType: string) => {
	return OPERATIONS[operationType]?.name ?? operationType;
};

const throttledSwitchState = useThrottleFn((v: boolean) => {
	eventsPatcher.mutate({ id: props.event.id!, enabled: v });
}, 500);
</script>

<template>
	<n-card segmented header-style="padding: 10px; height: 100%;">
		<template #header>
			<n-space align="center">
				<component :is="EVENTS[event.type]!.icon" v-if="EVENTS[event.type].icon" style="display: flex" />
				<n-text>{{ getEventName(event.type) }}</n-text>
			</n-space>
		</template>

		<template #header-extra>
			<n-switch
				:value="event.enabled"
				@update-value="(v) => throttledSwitchState(v)"
			/>
		</template>

		<n-space vertical>
			<n-text>{{ event.description }}</n-text>
		</n-space>

		<template #footer>
			<n-row style="gap: 8px;">
				<n-tooltip v-for="(operation, index) of event.operations" :key="index" :disabled="!operation.input">
					<template #trigger>
						<n-tag :disabled="!operation.enabled" :bordered="false" type="info">
							{{ getOperationName(operation.type) }}
						</n-tag>
					</template>
					<n-space vertical>
						<n-text>{{ operation.input }}</n-text>
						<n-text>Delay: {{ operation.delay }} | Repeat: {{ operation.repeat }}</n-text>
					</n-space>
				</n-tooltip>
			</n-row>
		</template>
	</n-card>
</template>
