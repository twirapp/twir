import { getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { PlusIcon } from 'lucide-vue-next'
import { computed, h } from 'vue'

import type { Giveaway } from '@/api/giveaways.ts'
import type { ColumnDef } from '@tanstack/vue-table'

import { Button } from '@/components/ui/button'
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'

export const useGiveawaysListTable = createGlobalState(() => {
	const { giveawaysList, giveawaysListFetching } = useGiveaways()

	const tableColumns = computed<ColumnDef<Giveaway>[]>(() => {
		return [
			{
				accessorKey: 'text',
				size: 10,
				header: () => h('div', {}, 'Keyword'),
				cell: ({ row }) => h('span', row.original.keyword),
			},
			{
				accessorKey: 'createdAt',
				size: 10,
				header: () => h('div', {}, 'Created at'),
				cell: ({ row }) => h('span', new Date(row.original.createdAt).toLocaleString()),
			},
			{
				accessorKey: 'endedAt',
				size: 10,
				header: () => h('div', {}, 'Ended at'),
				cell: ({ row }) => h('span', new Date(row.original.endedAt).toLocaleString()),
			},
			{
				accessorKey: 'actions',
				size: 20,
				header: () => h(Button, { size: 'sm', class: 'flex gap-2 items-center' }, {
					default: () => [
						h(PlusIcon),
						'Create new',
					],
				}),
			},
		]
	})

	const table = useVueTable({
		get data() {
			return giveawaysList.value
		},
		get columns() {
			return tableColumns.value
		},
		getCoreRowModel: getCoreRowModel(),
	})

	return {
		isLoading: giveawaysListFetching,
		table,
	}
})
