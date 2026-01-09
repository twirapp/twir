import type { AdminShortUrl } from '#layers/dashboard/api/admin/short-urls.ts'
import type { ColumnDef } from '@tanstack/vue-table'

import { getCoreRowModel, getPaginationRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import UsersTableCellUser from '~/features/admin-panel/manage-users/ui/users-table-cell-user.vue'
import ShortUrlsActions from '~/features/admin-panel/short-urls/ui/short-urls-actions.vue'
import { TABLE_ACCESSOR_KEYS } from '~/features/community-users/composables/use-community-users-table.ts'

import { useAdminShortUrlsApi } from './use-admin-short-urls-api.ts'

export const useAdminShortUrlsTable = createGlobalState(() => {
	const { list, pagination, totalItems } = useAdminShortUrlsApi()

	const tableColumns = computed<ColumnDef<AdminShortUrl>[]>(() => [
		{
			accessorKey: 'shortId',
			size: 5,
			header: () => h('div', {}, 'Short ID'),
			cell: ({ row }) => {
				return h(
					'a',
					{
						href: `${window.location.origin}/s/${row.original.id}`,
						target: '_blank',
						class: 'underline',
					},
					row.original.id
				)
			},
		},
		{
			accessorKey: 'link',
			size: 10,
			header: () => h('div', {}, 'Link'),
			cell: ({ row }) => {
				return h(
					'a',
					{ href: row.original.link, target: '_blank', class: 'underline' },
					row.original.link
				)
			},
		},
		{
			accessorKey: 'views',
			size: 5,
			header: () => h('div', {}, 'Views'),
			cell: ({ row }) => {
				return h('span', row.original.views)
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
			accessorKey: TABLE_ACCESSOR_KEYS.user,
			size: 20,
			header: () => h('div', {}, 'User'),
			cell: ({ row }) => {
				if (!row.original.userProfile) return

				return h(
					'a',
					{
						class: 'flex flex-col',
						href: `https://twitch.tv/${row.original.userProfile.login}`,
						target: '_blank',
					},
					h(UsersTableCellUser, {
						avatar: row.original.userProfile.profileImageUrl,
						name: row.original.userProfile.login,
						displayName: row.original.userProfile.displayName,
					})
				)
			},
		},
		{
			accessorKey: 'actions',
			size: 5,
			header: () => '',
			cell: ({ row }) => h(ShortUrlsActions, { item: row.original }),
		},
	])

	const table = useVueTable({
		get pageCount() {
			return totalItems.value
		},
		get data() {
			return list.value
		},
		get columns() {
			return tableColumns.value
		},
		state: {
			get pagination() {
				return pagination.value
			},
		},
		getCoreRowModel: getCoreRowModel(),
		manualPagination: true,
		getPaginationRowModel: getPaginationRowModel(),
	})

	return {
		table,
	}
})
