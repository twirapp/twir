<script setup lang="ts">
import { IconSettings } from '@tabler/icons-vue';
import { useThrottleFn } from '@vueuse/core';
import { NText, NButton, NTooltip, NTag, NRow, NSpace, NSwitch } from 'naive-ui';

import { EVENTS } from './events.js';
import { OPERATIONS } from './operations.js';
import { EditableEvent } from './types.js';

import { useEventsManager } from '@/api/index.js';
import Card from '@/components/card/card.vue';

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
			</n-space>
		</template>

		<template #footer>
			<n-button secondary size="large" @click="$emit('openSettings')">
				<span>Settings</span>
				<IconSettings />
			</n-button>
			<n-button secondary type="error" size="large" @click="$emit('openSettings')">
				<span>Delete</span>
			</n-button>
		</template>
	</card>
</template>
