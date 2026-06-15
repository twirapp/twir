<script setup lang="ts">
import EventEdit from '~~/layers/dashboard/app/features/events/event-form.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<EventEdit />
</template>
