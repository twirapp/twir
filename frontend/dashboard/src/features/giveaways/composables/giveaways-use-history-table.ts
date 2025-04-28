import { getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { EyeIcon } from 'lucide-vue-next'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import type { Giveaway } from '@/api/giveaways.ts'
import type { ColumnDef } from '@tanstack/vue-table'

import { Button } from '@/components/ui/button'
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'

export const useGiveawaysHistoryTable = createGlobalState(() => {
	const { t } = useI18n()
	const {
		archivedGiveaways,
		giveawaysListFetching,
		viewGiveaway,
	} = useGiveaways()

	const tableColumns = computed<ColumnDef<Giveaway>[]>(() => {
		return [
			{
				accessorKey: 'keyword',
				size: 20,
				header: () => h('div', {}, t('giveaways.keyword')),
				cell: ({ row }) => h('div', { class: 'flex items-center gap-2' }, [
					h('span', {}, row.original.keyword),
					// row.original.archivedAt ?
					//   h(Badge, { variant: 'outline' }, () => 'Archived') :
					//   h(Badge, { variant: 'secondary' }, () => 'Ended')
				]),
			},
			{
				accessorKey: 'createdAt',
				size: 20,
				header: () => h('div', {}, t('giveaways.createdAt')),
				cell: ({ row }) => h('span', {}, new Date(row.original.createdAt).toLocaleString()),
			},
			{
				accessorKey: 'startedAt',
				size: 20,
				header: () => h('div', {}, t('giveaways.startedAt')),
				cell: ({ row }) => h('span', {}, row.original.startedAt ? new Date(row.original.startedAt).toLocaleString() : '-'),
			},
			{
				accessorKey: 'stoppedAt',
				size: 20,
				header: () => h('div', {}, t('giveaways.stoppedAt')),
				cell: ({ row }) => h('span', {}, row.original.stoppedAt ? new Date(row.original.stoppedAt).toLocaleString() : '-'),
			},
			{
				accessorKey: 'winners',
				size: 20,
				header: () => h('div', {}, t('giveaways.winners')),
				cell: ({ row }) => h('span', {}, row.original.winners?.length || 0),
			},
			{
				accessorKey: 'actions',
				size: 20,
				header: () => h('div', {}),
				cell: ({ row }) => h('div', { class: 'flex gap-2 justify-end' }, [
					// View button
					h(Button, {
						size: 'sm',
						variant: 'outline',
						class: 'flex gap-2 items-center',
						onClick: () => viewGiveaway(row.original.id),
					}, {
						default: () => [
							h(EyeIcon, { class: 'size-4' }),
							t('giveaways.view'),
						],
					}),
				]),
			},
		]
	})

	const table = useVueTable({
		get data() {
			return archivedGiveaways.value
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
