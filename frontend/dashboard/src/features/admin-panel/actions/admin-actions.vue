<script setup lang="ts">
import DropSessions from './ui/drop-sessions.vue'
import EventsubSubscribe from './ui/eventsub-subscribe.vue'
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<EventsubSubscribe />
		<DropSessions />
	</div>
</template>
