import { type ColumnDef, getCoreRowModel, getFacetedRowModel, getFacetedUniqueValues, getFilteredRowModel, getPaginationRowModel, getSortedRowModel, useVueTable } from '@tanstack/vue-table'

import { useCommunityUsersSorting } from './use-community-users-sorting.js'
import CommunityUsersTableColumn from '../ui/community-users-table-column.vue'

import UserCell from '~/components/table/cells/user-cell.vue'
import { graphql } from '~/gql/gql.js'
import { type CommunityUser, type CommunityUsersOpts, CommunityUsersResetType } from '~/gql/graphql.js'
import { valueUpdater } from '~/lib/utils.js'

const ONE_HOUR = 60 * 60 * 1000
export const TABLE_ACCESSOR_KEYS = {
	user: 'user',
	messages: 'messages',
	watchedMs: 'watchedMs',
	usedEmotes: 'usedEmotes',
	usedChannelPoints: 'usedChannelPoints',
}

export const useCommunityUsersTable = defineStore('community-users', () => {
	const currentChannelId = useCurrentChannelId()

	const {
		sorting,
		columnFilters,
		columnVisibility,
		rowSelection,
		tableOrder,
		tableSortBy,
		debouncedSearchInput,
	} = storeToRefs(useCommunityUsersSorting())
	const { pagination, setPagination } = usePagination()

	const communityUsersOpts = computed<CommunityUsersOpts>((prevParams) => {
		if (prevParams?.search !== debouncedSearchInput.value) {
			pagination.value.pageIndex = 0
		}

		return {
			search: debouncedSearchInput.value,
			channelId: currentChannelId.value || '',
			page: pagination.value.pageIndex,
			perPage: pagination.value.pageSize,
			order: tableOrder.value,
			sortBy: tableSortBy.value,
		}
	})

	const communityUsers = useQuery({
		get variables() {
			return {
				opts: communityUsersOpts.value,
			}
		},
		query: graphql(`
			query GetAllCommunityUsers($opts: CommunityUsersOpts!) {
				communityUsers(opts: $opts) {
					total
					users {
						id
						twitchProfile {
							login
							displayName
							profileImageUrl
						}
						watchedMs
						messages
						usedEmotes
						usedChannelPoints
					}
				}
			}
		`),
	})

	const totalUsers = computed(() => communityUsers.data.value?.communityUsers.total ?? 0)
	const pageCount = computed(() => {
		return Math.ceil(totalUsers.value / pagination.value.pageSize)
	})

	const tableColumns = computed<ColumnDef<CommunityUser>[]>(() => [
		{
			// accessorKey: t('community.users.table.user'),
			accessorKey: TABLE_ACCESSOR_KEYS.user,
			size: 20,
			header: () => h('div', {}, 'User'),
			cell: ({ row }) => {
				return h('a', {
					class: 'flex flex-col',
					href: `https://twitch.tv/${row.original.twitchProfile.login}`,
					target: '_blank',
				}, h(UserCell, {
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
					title: 'Messages',
				})
			},
			cell: ({ row }) => {
				return h('div', row.original.messages)
			},
		},
		{
			accessorKey: TABLE_ACCESSOR_KEYS.usedChannelPoints,
			size: 20,
			header: ({ column }) => {
				return h(CommunityUsersTableColumn, {
					column,
					columnType: CommunityUsersResetType.UsedChannelsPoints,
					title: 'Used channel points',
				})
			},
			cell: ({ row }) => {
				return h('div', row.original.usedChannelPoints)
			},
		},
		{
			accessorKey: TABLE_ACCESSOR_KEYS.usedEmotes,
			size: 20,
			header: ({ column }) => {
				return h(CommunityUsersTableColumn, {
					column,
					columnType: CommunityUsersResetType.UsedEmotes,
					title: 'Used emotes',
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
					title: 'Watched time',
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
			return communityUsers.data.value?.communityUsers.users ?? []
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
		table,
		fetchUsers: communityUsers.executeQuery,
	}
})

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useCommunityUsersTable, import.meta.hot))
}
