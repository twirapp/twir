<script setup lang="ts">
import TimerEdit from '~~/layers/dashboard/features/timers/timers-edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<TimerEdit />
</template>
