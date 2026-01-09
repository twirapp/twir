import type { EmotesStatisticsDetail } from '#layers/dashboard/api/emotes-statistic'
import type { ColumnDef } from '@tanstack/vue-table'

import { valueUpdater } from '#layers/dashboard/helpers/value-updater'
import {
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getPaginationRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { UseTimeAgo } from '@vueuse/components'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import UsersTableCellUser from '~/features/admin-panel/manage-users/ui/users-table-cell-user.vue'

import { useCommunityEmotesDetails } from './use-community-emotes-details'

type UserUsage = NonNullable<
	EmotesStatisticsDetail['emotesStatisticEmoteDetailedInformation']
>['usagesHistory'][number]

export const useCommunityEmotesDetailsUsersHistory = createGlobalState(() => {
	const { details, usagesPagination } = useCommunityEmotesDetails()

	const data = computed<UserUsage[]>(() => {
		return details.value?.emotesStatisticEmoteDetailedInformation?.usagesHistory ?? []
	})
	const total = computed(
		() => details.value?.emotesStatisticEmoteDetailedInformation?.usagesByUsersTotal ?? 0
	)
	const pageCount = computed(() => {
		return Math.ceil(total.value / usagesPagination.value.pageSize)
	})

	const columns = computed<ColumnDef<UserUsage>[]>(() => [
		{
			accessorKey: 'user',
			size: 50,
			header: () => '',
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
			accessorKey: 'time',
			header: '',
			cell: ({ row }) => {
				const date = new Date(row.original.date)
				return h(
					UseTimeAgo,
					{ time: date },
					{
						default: ({ timeAgo }: { timeAgo: string }) => timeAgo,
					}
				)
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
})
