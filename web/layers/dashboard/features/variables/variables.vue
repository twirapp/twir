<script setup lang="ts">
import Table from '~~/layers/dashboard/components/table.vue'
import VariablesInformation from '~~/layers/dashboard/features/variables/ui/variables-information.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

import { useVariablesTable } from './composables/use-variables-table'
import VariablesCreateButton from './ui/variables-create-button.vue'

const variablesTable = useVariablesTable()
const { t } = useI18n()
</script>

<template>
	<div>
		<div class="flex flex-col gap-4">
			<div class="flex flex-row gap-10">
				<VariablesInformation />
				<VariablesCreateButton />
			</div>
			<Table
				:table="variablesTable.table"
				:is-loading="variablesTable.isLoading.value"
			/>
		</div>
	</div>
</template>
