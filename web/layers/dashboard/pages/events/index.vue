<script setup lang="ts">
import EventsList from '~~/layers/dashboard/features/events/events-list.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<EventsList />
</template>
