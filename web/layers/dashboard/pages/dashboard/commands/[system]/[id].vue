<script setup lang="ts">
import CommandEdit from '~~/layers/dashboard/features/commands/commands-edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<CommandEdit />
</template>
