import {
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getPaginationRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { NTime } from 'naive-ui'
import { computed, h } from 'vue'

import { useCommunityEmotesDetails } from './use-community-emotes-details'

import type { EmotesStatisticsDetail } from '@/api/emotes-statistic'
import type { ColumnDef } from '@tanstack/vue-table'

import UsersTableCellUser
	from '@/features/admin-panel/manage-users/components/users-table-cell-user.vue'
import { valueUpdater } from '@/helpers/value-updater'

type UserUsage = NonNullable<EmotesStatisticsDetail['emotesStatisticEmoteDetailedInformation']>['usagesHistory'][number]

export const useCommunityEmotesDetailsUsersHistory = createGlobalState(() => {
	const { details, usagesPagination } = useCommunityEmotesDetails()

	const data = computed<UserUsage[]>(() => {
		return details.value?.emotesStatisticEmoteDetailedInformation?.usagesHistory ?? []
	})
	const total = computed(() => details.value?.emotesStatisticEmoteDetailedInformation?.usagesByUsersTotal ?? 0)
	const pageCount = computed(() => {
		return Math.ceil(total.value / usagesPagination.value.pageSize)
	})

	const columns = computed<ColumnDef<UserUsage>[]>(() => [
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
			accessorKey: 'time',
			header: '',
			cell: ({ row }) => {
				const date = new Date(row.original.date)
				const diff = Date.now() - date.getTime()
				return h(NTime, { type: 'relative', time: 0, to: diff })
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
				return usagesPagination.value
			},
		},
		manualPagination: true,
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, usagesPagination),
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
	})

	return {
		table,
		usagesPagination,
		total,
		pageCount,
	}
},
)
