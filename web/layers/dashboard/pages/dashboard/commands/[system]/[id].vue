<script setup lang="ts">
import CommandEdit from '~~/layers/dashboard/features/commands/commands-edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<CommandEdit />
</template>
