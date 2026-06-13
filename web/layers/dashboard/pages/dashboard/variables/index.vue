<script setup lang="ts">
import VariablesList from '~~/layers/dashboard/features/variables/variables-list.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<VariablesList />
</template>
