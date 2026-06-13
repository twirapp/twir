import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import { useAdminToxicMessagesApi } from './use-admin-toxic-mesasges-api.ts'

import type { AdminToxicMessagesQuery } from '@/gql/graphql.ts'

import UsersTableCellUser from '@/features/admin-panel/manage-users/ui/users-table-cell-user.vue'

export const useAdminToxicMessagesTable = createGlobalState(() => {
	const { list, pagination, totalItems } = useAdminToxicMessagesApi()

	const tableColumns = computed<ColumnDef<AdminToxicMessagesQuery['adminToxicMessages']['items'][0]>[]>(() => [
		{
			accessorKey: 'channel',
			size: 15,
			minSize: 15,
			header: () => h('div', {}, { default: () => 'Channel' }),
			cell: ({ row }) => {
				if (!row.original.channelProfile) {
					return null
				}

				return h('a', {
					class: 'flex flex-col',
					href: `https://twitch.tv/${row.original.channelProfile.login}`,
					target: '_blank',
				}, h(UsersTableCellUser, {
					avatar: row.original.channelProfile.profileImageUrl,
					name: row.original.channelProfile.login,
					displayName: row.original.channelProfile.displayName,
				}))
			},
		},
		{
			accessorKey: 'text',
			size: 60,
			header: () => h('div', {}, { default: () => 'Text' }),
			cell: ({ row }) => {
				return h('span', row.original.text)
			},
		},
		{
			accessorKey: 'createdAt',
			size: 10,
			header: () => h('div', {}, { default: () => 'Created at' }),
			cell: ({ row }) => {
				return h('span', new Date(row.original.createdAt).toLocaleString())
			},
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
	})

	return {
		table,
	}
})
