<script setup lang="ts">
import EventEdit from '~~/layers/dashboard/features/events/event-form.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<EventEdit />
</template>
