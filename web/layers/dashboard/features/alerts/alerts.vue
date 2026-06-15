<script setup lang="ts">
import Table from '~~/layers/dashboard/components/table.vue'

import { useAlertsTable } from './composables/use-alerts-table.js'
import AlertsInformation from './ui/alerts-information.vue'

const { t } = useI18n()
const { table, isLoading } = useAlertsTable()
</script>

<template>
	<div class="flex w-full flex-col gap-4">
		<AlertsInformation />
		<Table
			:table="table"
			:is-loading="isLoading"
		>
			<template #empty-message>
				{{ t('alerts.tableEmpty') }}
			</template>
		</Table>
	</div>
</template>
