<script setup lang="ts">
import { NCard, NTag } from 'naive-ui';

import { EVENTS } from './events.js';
import { OPERATIONS } from './operations.js';
import { EditableEvent } from './types.js';

defineProps<{
	event: EditableEvent
}>();
</script>

<template>
	<n-card :title="EVENTS[event.type]?.shortName ?? event.type">
		{{ event.description }}

		<template #footer>
			<n-tag v-for="(operation, index) of event.operations" :key="index">
				{{ OPERATIONS[operation.type]?.shortName ?? operation.type }}
			</n-tag>
		</template>
	</n-card>
</template>
