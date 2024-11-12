import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import type { TimerResponse } from '@/api/timers'

import { useTimersApi } from '@/api/timers'
import TimersTableActions from '@/features/timers/ui/timers-table-actions.vue'

export const useTimersTable = createGlobalState(() => {
	const { t } = useI18n()
	const timersApi = useTimersApi()

	const { data, fetching } = timersApi.useQueryTimers()
	const timers = computed<TimerResponse[]>(() => {
		if (!data.value) return []
		return data.value.timers
	})

	const tableColumns = computed<ColumnDef<TimerResponse>[]>(() => {
		return [
			{
				accessorKey: 'text',
				size: 10,
				header: () => h('div', {}, t('sharedTexts.name')),
				cell: ({ row }) => h('span', row.original.name),
			},
			{
				accessorKey: 'responses',
				size: 75,
				header: () => h('div', {}, t('sharedTexts.responses')),
				cell: ({ row }) => h(
					'div',
					{ class: 'flex flex-col gap-0.5' },
					row.original.responses.map(r => h('span', { class: 'truncate md:whitespace-normal' }, r.text)),
				),
			},
			{
				accessorKey: 'timeInterval',
				size: 5,
				header: () => h('div', {}, t('timers.table.columns.intervalInMinutes')),
				cell: ({ row }) => h('span', row.original.timeInterval),
			},
			{
				accessorKey: 'messageInterval',
				size: 5,
				header: () => h('div', {}, t('timers.table.columns.intervalInMessages')),
				cell: ({ row }) => h('span', row.original.messageInterval),
			},
			{
				accessorKey: 'actions',
				header: '',
				size: 5,
				cell: ({ row }) => h(TimersTableActions, { timer: row.original }),
			},
		]
	})

	const table = useVueTable({
		get data() {
			return timers.value
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
})
