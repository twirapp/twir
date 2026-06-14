<script setup lang="ts">
import VariableEdit from '~~/layers/dashboard/features/variables/variables-edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<VariableEdit />
</template>
