<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { useAlertsTable } from './composables/use-alerts-table.js'
import AlertsInformation from './ui/alerts-information.vue'

import Table from '@/components/table.vue'

const { t } = useI18n()
const { table, isLoading } = useAlertsTable()
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<AlertsInformation />
		<Table :table="table" :is-loading="isLoading">
			<template #empty-message>
				{{ t('alerts.tableEmpty') }}
			</template>
		</Table>
	</div>
</template>
