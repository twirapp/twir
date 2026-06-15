<script setup lang="ts">
import CommandsBySystem from '~~/layers/dashboard/app/features/commands/commands.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<CommandsBySystem />
</template>
