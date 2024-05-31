import {
	type ColumnDef,
	getCoreRowModel,
	getPaginationRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import { useUsersTableFilters } from './use-users-table-filters.js'
import { useUsers } from './use-users.js'
import UsersActionSelector from '../ui/users-action-selector.vue'
import UsersBadgeSelector from '../ui/users-badge-selector.vue'
import UsersTableCellUser from '../ui/users-table-cell-user.vue'

import type { User } from '@/api/admin/users.js'
import type { TwirUsersSearchParams } from '@/gql/graphql'

import { usePagination } from '@/composables/use-pagination.js'
import { valueUpdater } from '@/helpers/value-updater.js'

export const useUsersTable = createGlobalState(() => {
	const { t } = useI18n()

	const { pagination } = usePagination()

	const tableFilters = useUsersTableFilters()

	const tableParams = computed<TwirUsersSearchParams>((prevParams) => {
		// reset pagination on search change
		if (prevParams?.search !== tableFilters.debounceSearchInput.value) {
			pagination.value.pageIndex = 0
		}

		return {
			...tableFilters.selectedStatuses.value,
			search: tableFilters.debounceSearchInput.value,
			page: pagination.value.pageIndex,
			perPage: pagination.value.pageSize,
			badges: tableFilters.selectedBadges.value,
		}
	})

	const { usersApi } = useUsers()
	const { data, fetching } = usersApi.useQueryUsers(tableParams)

	const users = computed<User[]>(() => {
		if (!data.value) return []
		return data.value.twirUsers.users
	})

	const totalUsers = computed(() => data.value?.twirUsers.total ?? 0)

	const pageCount = computed(() => {
		return Math.ceil(totalUsers.value / pagination.value.pageSize)
	})

	const tableColumns = computed<ColumnDef<User>[]>(() => [
		{
			accessorKey: 'user',
			size: 60,
			header: () => h('div', {}, t('adminPanel.manageUsers.user')),
			cell: ({ row }) => {
				return h('a', {
					class: 'flex flex-col',
					href: `https://twitch.tv/${row.original.twitchProfile.login}`,
					target: '_blank',
				}, h(UsersTableCellUser, {
					avatar: row.original.twitchProfile.profileImageUrl,
					name: row.original.twitchProfile.login,
					displayName: row.original.twitchProfile.displayName,
				}))
			},
		},
		{
			accessorKey: 'userId',
			size: 30,
			header: () => h('div', {}, t('adminPanel.manageUsers.userId')),
			cell: ({ row }) => {
				return h('span', row.original.id)
			},
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => {
				return h(
					'div',
					{ class: 'flex items-center justify-end gap-2' },
					[
						h(UsersBadgeSelector, {
							userId: row.original.id,
						}),
						h(UsersActionSelector, {
							userId: row.original.id,
							isBanned: row.original.isBanned,
							isBotAdmin: row.original.isBotAdmin,
						}),
					],
				)
			},
		},
	])

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		get data() {
			return users.value
		},
		get columns() {
			return tableColumns.value
		},
		state: {
			get pagination() {
				return pagination.value
			},
		},
		manualPagination: true,
		enableRowSelection: true,
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
	})

	return {
		isLoading: fetching,
		pagination,
		totalUsers,
		table,
		tableColumns,
	}
})
