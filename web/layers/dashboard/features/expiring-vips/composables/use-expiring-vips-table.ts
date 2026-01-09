import type { ColumnDef } from '@tanstack/vue-table'

import { getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import type { ScheduledVip } from '#layers/dashboard/api/scheduled-vips.ts'

import { useScheduledVipsApi } from '#layers/dashboard/api/scheduled-vips.ts'
import UsersTableCellUser from '~/features/admin-panel/manage-users/ui/users-table-cell-user.vue'

import ExpiringVipsTableActions from '../ui/expiring-vips-table-actions.vue'

export const useExpiringVipsTable = createGlobalState(() => {
	const { t } = useI18n()
	const api = useScheduledVipsApi()
	const { data, fetching } = api.useQueryScheduledVips()

	const scheduledVips = computed(() => {
		return data.value?.scheduledVips ?? []
	})

	const tableColumns = computed<ColumnDef<ScheduledVip>[]>(() => [
		{
			accessorKey: 'user',
			size: 40,
			header: () => h('div', {}, t('sharedTexts.user')),
			cell: ({ row }) => {
				return h(
					'a',
					{
						class: 'flex flex-col',
						href: `https://twitch.tv/${row.original.twitchProfile.login}`,
						target: '_blank',
					},
					h(UsersTableCellUser, {
						avatar: row.original.twitchProfile.profileImageUrl,
						name: row.original.twitchProfile.login,
						displayName: row.original.twitchProfile.displayName,
					})
				)
			},
		},
		{
			accessorKey: 'createdAt',
			size: 10,
			header: () => h('div', {}, 'Created at'),
			cell: ({ row }) => {
				return h('span', new Date(row.original.createdAt).toLocaleString())
			},
		},
		{
			accessorKey: 'removeAt',
			size: 10,
			header: () => h('div', {}, 'Expire at'),
			cell: ({ row }) => {
				return h(
					'span',
					row.original.removeAt ? new Date(row.original.removeAt).toLocaleString() : '-'
				)
			},
		},
		{
			accessorKey: 'actions',
			size: 5,
			header: () => '',
			cell: ({ row }) => h(ExpiringVipsTableActions, { scheduledVip: row.original }),
		},
	])

	const table = useVueTable({
		get data() {
			return scheduledVips.value
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
