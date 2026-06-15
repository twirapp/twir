import {
	type ColumnDef,
	getCoreRowModel,
	getPaginationRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import { useUsersTableFilters } from './use-users-table-filters.js'
import { useUsers } from './use-users.js'
import UsersActionSelector from '../ui/users-action-selector.vue'
import UsersBadgeSelector from '../ui/users-badge-selector.vue'
import UsersTableCellUser from '../ui/users-table-cell-user.vue'

import type { User } from '~~/layers/dashboard/app/api/admin/users.js'
import type { TwirUsersSearchParams } from '~/gql/graphql.js'

import { usePagination } from '~~/layers/dashboard/app/composables/use-pagination.js'
import { resolveProfile } from '~~/layers/dashboard/app/helpers/resolveProfile.js'
import { valueUpdater } from '~~/layers/dashboard/app/helpers/value-updater.js'

export const useUsersTable = createGlobalState(() => {
	const { t } = useI18n()

	const { pagination } = usePagination()

	const tableFilters = useUsersTableFilters()

	const tableParams = computed<TwirUsersSearchParams>((prevParams) => {
		const currentSearch = tableFilters.debounceSearchInput.value
		const currentPlatforms = tableFilters.selectedPlatforms.value

		// reset pagination on search change
		if (
			prevParams?.search !== currentSearch
			|| JSON.stringify(prevParams?.platforms ?? []) !== JSON.stringify(currentPlatforms)
		) {
			pagination.value.pageIndex = 0
		}

		return {
			...tableFilters.selectedStatuses.value,
			search: currentSearch,
			page: pagination.value.pageIndex,
			perPage: pagination.value.pageSize,
			badges: tableFilters.selectedBadges.value,
			platforms: [...currentPlatforms],
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
				const profile = resolveProfile({
					profileImageUrl: row.original.avatar,
					login: row.original.login,
					displayName: row.original.displayName,
					platform: row.original.platform,
				})

				return h(UsersTableCellUser, {
					avatar: profile.avatar,
					name: profile.login,
					displayName: profile.displayName,
					url: profile.url,
					platform: row.original.platform,
				})
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
