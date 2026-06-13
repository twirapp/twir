<script setup lang="ts">
import VariableEdit from '~/layers/dashboard/features/variables/variable-edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<VariableEdit />
</template>
