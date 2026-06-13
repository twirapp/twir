<script setup lang="ts">
import TimersList from '~/layers/dashboard/features/timers/timers-list.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<TimersList />
</template>
