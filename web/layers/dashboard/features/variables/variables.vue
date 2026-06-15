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
	<PageLayout>
		<template #title>
			{{ t('sidebar.variables') }}
		</template>

		<template #action>
			<VariablesCreateButton />
		</template>

		<template #content>
			<div class="flex flex-col gap-4">
				<VariablesInformation />
				<Table
					:table="variablesTable.table"
					:is-loading="variablesTable.isLoading.value"
				/>
			</div>
		</template>
	</PageLayout>
</template>
