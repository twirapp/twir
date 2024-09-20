<script setup lang="ts">
import { useVariablesTable } from './composables/use-variables-table'

import Table from '@/components/table.vue'
import VariablesInformation from '@/features/variables/ui/variables-information.vue'

const variablesTable = useVariablesTable()
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<VariablesInformation />
		<Table :table="variablesTable.table" :is-loading="variablesTable.isLoading.value" />
	</div>
</template>
