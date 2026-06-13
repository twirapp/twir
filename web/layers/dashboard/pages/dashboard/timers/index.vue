<script setup lang="ts">
import TimersList from '~~/layers/dashboard/features/timers/timers.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<TimersList />
</template>
