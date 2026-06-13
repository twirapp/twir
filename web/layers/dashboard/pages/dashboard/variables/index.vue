<script setup lang="ts">
import VariablesList from '~~/layers/dashboard/features/variables/variables.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<VariablesList />
</template>
