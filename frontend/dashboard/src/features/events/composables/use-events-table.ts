import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import EventsTableActions from '../ui/events-table-actions.vue'

import { type Event, useEventsApi } from '@/api/events'
import { flatEvents, getEventName } from '@/features/events/constants/helpers'

export const useEventsTable = createGlobalState(() => {
	const { t } = useI18n()
	const eventsApi = useEventsApi()

	const { data, fetching } = eventsApi.useQueryEvents()
	const events = computed<Event[]>(() => {
		if (!data.value) return []
		return data.value.events
	})

	const tableColumns = computed<ColumnDef<Event>[]>(() => [
		{
			accessorKey: 'type',
			size: 10,
			header: () => h('div', {}, t('events.type')),
			cell: ({ row }) => h('div', { class: 'flex items-center gap-2' }, [
				h(flatEvents[row.original.type]?.icon ?? 'div'),
				h('span', {}, getEventName(row.original.type)),
			]),
		},
		{
			accessorKey: 'description',
			size: 20,
			header: () => h('div', {}, t('events.description')),
			cell: ({ row }) => h('span', row.original.description),
		},
		{
			accessorKey: 'operations',
			size: 60,
			header: () => h('div', {}, t('events.operations')),
			cell: ({ row }) => h('span', {}, `${row.original.operations.length} ${t('events.operations')}`),
		},
		{
			accessorKey: 'actions',
			size: 5,
			header: () => '',
			cell: ({ row }) => h(EventsTableActions, { event: row.original }),
		},
	])

	const table = useVueTable({
		get data() {
			return events.value
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
