<script setup lang="ts">
import EventEdit from '~/layers/dashboard/features/events/event-edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<EventEdit />
</template>
