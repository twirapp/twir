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
import { formatDistance } from 'date-fns'
import { computed, h } from 'vue'

import CommunityRewardsTableRewardCell from '../ui/cells/community-rewards-history-table-reward-cell.vue'

import type { Redemption } from '@/api/community-rewards'

import { useCommunityRewardsApi } from '@/api/community-rewards'
import UsersTableCellUser from '@/features/admin-panel/manage-users/ui/users-table-cell-user.vue'
import {
	useCommunityRewardsHistoryQuery,
} from '@/features/community-rewards-history/composables/community-rewards-history-query.ts'
import { valueUpdater } from '@/helpers/value-updater'

export const useCommunityRewardsTable = createGlobalState(() => {
	const communityRewardsApi = useCommunityRewardsApi()

	const { query, pagination } = useCommunityRewardsHistoryQuery()
	const historyResult = communityRewardsApi.useHistory(query)

	const history = computed<Redemption[]>(() => {
		return historyResult.data.value?.rewardsRedemptionsHistory.redemptions ?? []
	})
	const total = computed(() => {
		return historyResult.data.value?.rewardsRedemptionsHistory.total ?? 0
	})
	const pageCount = computed(() => {
		return Math.ceil(total.value / pagination.value.pageSize)
	})

	const tableColumns = computed<ColumnDef<Redemption>[]>(() => [
		{
			accessorKey: 'reward',
			size: 20,
			header: () => 'Reward',
			cell: ({ row }) => {
				return h(CommunityRewardsTableRewardCell, {
					name: row.original.reward.title,
					imageUrl: row.original.reward.imageUrls?.at(-1),
				})
			},
		},
		{
			accessorKey: 'user',
			size: 20,
			header: () => 'User',
			cell: ({ row }) => {
				return h('a', {
					class: 'flex flex-col',
					href: `https://twitch.tv/${row.original.user.login}`,
					target: '_blank',
				}, h(UsersTableCellUser, {
					avatar: row.original.user.profileImageUrl,
					name: row.original.user.login,
					displayName: row.original.user.displayName,
				}))
			},
		},
		{
			accessorKey: 'cost',
			size: 5,
			header: () => 'Cost',
			cell: ({ row }) => row.original.reward.cost,
		},
		{
			accessorKey: 'input',
			size: 25,
			header: () => 'User input',
			cell: ({ row }) => row.original.prompt,
		},
		{
			accessorKey: 'redeemedAt',
			size: 25,
			header: () => 'Redemed at',
			cell: ({ row }) => {
				return formatDistance(new Date(row.original.redeemedAt), new Date(), {
					addSuffix: true,
				})
			},
		},
	])

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		get data() {
			return history.value
		},
		get columns() {
			return tableColumns.value
		},
		manualPagination: true,
		enableRowSelection: true,
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
		// onSortingChange: updaterOrValue => valueUpdater(updaterOrValue, sorting),
		// onColumnFiltersChange: updaterOrValue => valueUpdater(updaterOrValue, columnFilters),
		// onColumnVisibilityChange: updaterOrValue => valueUpdater(updaterOrValue, columnVisibility),
		// onRowSelectionChange: updaterOrValue => valueUpdater(updaterOrValue, rowSelection),
		getCoreRowModel: getCoreRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
	})

	function refresh() {
		historyResult.executeQuery({ requestPolicy: 'cache-and-network' })
	}

	return {
		table,
		total,
		pageCount,
		isLoading: historyResult.fetching,
		refresh,
	}
})
