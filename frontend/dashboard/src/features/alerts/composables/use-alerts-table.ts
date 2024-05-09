import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import AlertsTableActions from '../ui/alerts-table-actions.vue'

import { type Alert, useAlertsQuery } from '@/api/alerts.js'
import { Badge } from '@/components/ui/badge'

export function useAlertsTable() {
	const { t } = useI18n()

	const { data, fetching } = useAlertsQuery()
	const greetings = computed<Alert[]>(() => {
		if (!data.value) return []
		return data.value.channelAlerts
	})

	const tableColumns = computed<ColumnDef<Alert>[]>(() => [
		{
			accessorKey: 'name',
			size: 25,
			header: () => h('div', {}, t('alerts.name')),
			cell: ({ row }) => {
				return h(Badge, {}, { default: () => row.original.name })
			},
		},
		{
			accessorKey: 'reward_ids',
			size: 25,
			header: () => h('div', {}, t('alerts.rewards')),
			cell: ({ row }) => h('span', row.original.reward_ids?.join(', ')),
		},
		{
			accessorKey: 'command_ids',
			size: 25,
			header: () => h('div', {}, t('alerts.commands')),
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => h(AlertsTableActions, { alert: row.original }),
		},
	])

	const table = useVueTable({
		get data() {
			return greetings.value
		},
		get columns() {
			return tableColumns.value
		},
		getCoreRowModel: getCoreRowModel(),
	})

	return {
		isLoading: fetching,
		table,
	}
}
