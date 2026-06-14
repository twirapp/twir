<script setup lang="ts">
import EventsList from '~~/layers/dashboard/features/events/events-list.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<EventsList />
</template>
