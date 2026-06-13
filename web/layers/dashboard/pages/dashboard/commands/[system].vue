<script setup lang="ts">
import CommandsBySystem from '~~/layers/dashboard/features/commands/commands-by-system.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<CommandsBySystem />
</template>
