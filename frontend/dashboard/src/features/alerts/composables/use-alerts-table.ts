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
			accessorKey: 'rewards',
			size: 25,
			header: () => h('div', {}, t('alerts.rewards')),
			cell: ({ row }) => h(
				'div',
				{ class: 'flex flex-col gap-0.5' },
				row.original.rewardIds?.map((id) => {
					const reward = data.value?.twitchGetChannelRewards.rewards.find((r) => r.id === id)
					return h('span', reward?.title)
				}),
			),
		},
		{
			accessorKey: 'commands',
			size: 25,
			header: () => h('div', {}, t('alerts.commands')),
			cell: ({ row }) => h(
				'div',
				{ class: 'flex flex-col gap-0.5' },
				row.original.commandIds?.map((id) => {
					const command = data.value?.commands.find((c) => c.id === id)
					return h('span', `!${command?.name}`)
				}),
			),
		},
		{
			accessorKey: 'keywords',
			size: 25,
			header: () => h('div', {}, t('alerts.trigger.keywords')),
			cell: ({ row }) => h(
				'div',
				{ class: 'flex flex-col gap-0.5' },
				row.original.keywordsIds?.map((id) => {
					const keyword = data.value?.keywords.find((k) => k.id === id)
					return h('span', keyword?.text)
				}),
			),
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
