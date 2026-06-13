<script setup lang="ts">
import CommandsBySystem from '~~/layers/dashboard/features/commands/commands.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<CommandsBySystem />
</template>
