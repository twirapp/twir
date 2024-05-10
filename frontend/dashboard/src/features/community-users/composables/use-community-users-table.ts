import {
	type ColumnDef,
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getFilteredRowModel,
	getPaginationRowModel,
	getSortedRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommunityTableActions } from './use-community-table-actions.js'
import CommunityUsersTableColumn from '../components/community-users-table-column.vue'

import { useProfile } from '@/api/auth.js'
import { type CommunityUser, useCommunityUsersApi } from '@/api/community-users.js'
import { usePagination } from '@/composables/use-pagination.js'
import UsersTableCellUser
	from '@/features/admin-panel/manage-users/components/users-table-cell-user.vue'
import { type CommunityUsersOpts, CommunityUsersResetType } from '@/gql/graphql.js'
import { valueUpdater } from '@/helpers/value-updater.js'

const ONE_HOUR = 60 * 60 * 1000
export const TABLE_ACCESSOR_KEYS = {
	user: 'user',
	messages: 'messages',
	watchedMs: 'watchedMs',
	usedEmotes: 'usedEmotes',
	usedChannelPoints: 'usedChannelPoints',
}

export const useCommunityUsersTable = createGlobalState(() => {
	const { t } = useI18n()

	const { data: profile } = useProfile()
	const communityUsersApi = useCommunityUsersApi()

	const {
		sorting,
		columnFilters,
		columnVisibility,
		rowSelection,
		tableOrder,
		tableSortBy,
		debouncedSearchInput,
	} = useCommunityTableActions()

	const { pagination, setPagination } = usePagination()
	const params = computed<CommunityUsersOpts>((prevParams) => {
		// reset pagination on search change
		if (prevParams?.search !== debouncedSearchInput.value) {
			pagination.value.pageIndex = 0
		}

		return {
			search: debouncedSearchInput.value,
			channelId: profile.value?.selectedDashboardId ?? '',
			page: pagination.value.pageIndex,
			perPage: pagination.value.pageSize,
			order: tableOrder.value,
			sortBy: tableSortBy.value,
		}
	})

	const { data, fetching } = communityUsersApi.useCommunityUsers(params)
	const communityUsers = computed<CommunityUser[]>(() => {
		if (!data.value) return []
		return data.value.communityUsers.users
	})

	const totalUsers = computed(() => data.value?.communityUsers.total ?? 0)
	const pageCount = computed(() => {
		return Math.ceil(totalUsers.value / pagination.value.pageSize)
	})

	const tableColumns = computed<ColumnDef<CommunityUser>[]>(() => [
		{
			// accessorKey: t('community.users.table.user'),
			accessorKey: TABLE_ACCESSOR_KEYS.user,
			size: 20,
			header: () => h('div', {}, t('community.users.table.user')),
			// header: ({ column }) => {
			// 	return h(CommunityUsersTableColumn, {
			// 		column,
			// 		title: t('community.users.table.user'),
			// 	});
			// },
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
			// accessorKey: t('community.users.table.messages'),
			accessorKey: TABLE_ACCESSOR_KEYS.messages,
			size: 20,
			header: ({ column }) => {
				return h(CommunityUsersTableColumn, {
					column,
					columnType: CommunityUsersResetType.Messages,
					title: t('community.users.table.messages'),
				})
			},
			cell: ({ row }) => {
				return h('div', row.original.messages)
			},
		},
		{
			// accessorKey: t('community.users.table.usedChannelPoints'),
			accessorKey: TABLE_ACCESSOR_KEYS.usedChannelPoints,
			size: 20,
			header: ({ column }) => {
				return h(CommunityUsersTableColumn, {
					column,
					columnType: CommunityUsersResetType.UsedChannelsPoints,
					title: t('community.users.table.usedChannelPoints'),
				})
			},
			cell: ({ row }) => {
				return h('div', row.original.usedChannelPoints)
			},
		},
		{
			// accessorKey: t('community.users.table.usedEmotes'),
			accessorKey: TABLE_ACCESSOR_KEYS.usedEmotes,
			size: 20,
			header: ({ column }) => {
				return h(CommunityUsersTableColumn, {
					column,
					columnType: CommunityUsersResetType.UsedEmotes,
					title: t('community.users.table.usedEmotes'),
				})
			},
			cell: ({ row }) => {
				return h('div', row.original.usedEmotes)
			},
		},
		{
			// accessorKey: t('community.users.table.watchedTime'),
			accessorKey: TABLE_ACCESSOR_KEYS.watchedMs,
			size: 20,
			header: ({ column }) => {
				return h(CommunityUsersTableColumn, {
					column,
					columnType: CommunityUsersResetType.Watched,
					title: t('community.users.table.watchedTime'),
				})
			},
			cell: ({ row }) => {
				return h('div', `${(Number(row.original.watchedMs) / ONE_HOUR).toFixed(1)}h`)
			},
		},
	])

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		get data() {
			return communityUsers.value
		},
		get columns() {
			return tableColumns.value
		},
		state: {
			get sorting() {
				return sorting.value
			},
			get columnFilters() {
				return columnFilters.value
			},
			get columnVisibility() {
				return columnVisibility.value
			},
			get rowSelection() {
				return rowSelection.value
			},
			get pagination() {
				return pagination.value
			},
		},
		manualPagination: true,
		enableRowSelection: true,
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
		onSortingChange: updaterOrValue => valueUpdater(updaterOrValue, sorting),
		onColumnFiltersChange: updaterOrValue => valueUpdater(updaterOrValue, columnFilters),
		onColumnVisibilityChange: updaterOrValue => valueUpdater(updaterOrValue, columnVisibility),
		onRowSelectionChange: updaterOrValue => valueUpdater(updaterOrValue, rowSelection),
		getCoreRowModel: getCoreRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
	})

	return {
		isLoading: fetching,
		table,
		totalUsers,
		pagination,
		setPagination,
	}
})
