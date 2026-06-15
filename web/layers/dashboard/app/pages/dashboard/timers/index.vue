<script setup lang="ts">
import TimersList from '~~/layers/dashboard/app/features/timers/timers.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<TimersList />
</template>
