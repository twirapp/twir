<script setup lang="ts">
import { EditableEvent } from './types';

defineProps<{
	event: EditableEvent
}>();
</script>

<template>
	<n-card :title="event.type">
		qwe
	</n-card>
</template>
