import {
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getPaginationRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import { useCommunityEmotesDetails } from './use-community-emotes-details'

import type { EmotesStatisticsDetail } from '@/api/emotes-statistic'
import type { ColumnDef } from '@tanstack/vue-table'

import UsersTableCellUser
	from '@/features/admin-panel/manage-users/components/users-table-cell-user.vue'
import { valueUpdater } from '@/helpers/value-updater'

type User = NonNullable<EmotesStatisticsDetail['emotesStatisticEmoteDetailedInformation']>['topUsers'][number]

export const useCommunityEmotesDetailsUsersTop = createGlobalState(() => {
	const { details, topPagination } = useCommunityEmotesDetails()

	const data = computed<User[]>(() => {
		return details.value?.emotesStatisticEmoteDetailedInformation?.topUsers ?? []
	})
	const total = computed(() => details.value?.emotesStatisticEmoteDetailedInformation?.topUsersTotal ?? 0)
	const pageCount = computed(() => {
		return Math.ceil(total.value / topPagination.value.pageSize)
	})

	const columns = computed<ColumnDef<User>[]>(() => [
		{
			accessorKey: 'user',
			size: 50,
			header: () => '',
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
			accessorKey: 'count',
			header: '',
			cell: ({ row }) => {
				return h('span', null, row.original.count)
			},
		},
	])

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		get data() {
			return data.value
		},
		get columns() {
			return columns.value
		},
		state: {
			get pagination() {
				return topPagination.value
			},
		},
		manualPagination: true,
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, topPagination),
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
	})

	return {
		table,
		topPagination,
		total,
		pageCount,
	}
},
)
