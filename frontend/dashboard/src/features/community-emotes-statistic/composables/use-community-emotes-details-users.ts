import {
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getPaginationRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { NTime } from 'naive-ui'
import { defineStore, storeToRefs } from 'pinia'
import { computed, h } from 'vue'

import { useCommunityEmotesDetails } from './use-community-emotes-details'

import type { EmotesStatisticsDetail } from '@/api/emotes-statistic'
import type { ColumnDef } from '@tanstack/vue-table'

import UsersTableCellUser
	from '@/features/admin-panel/manage-users/components/users-table-cell-user.vue'
import { resolveUserName } from '@/helpers'
import { valueUpdater } from '@/helpers/value-updater'

type UserUsage = NonNullable<EmotesStatisticsDetail['emotesStatisticEmoteDetailedInformation']>['usagesByUsers'][number]

export const useCommunityEmotesDetailsUsers = defineStore(
	'features/community-emotes-statistic-table/details-users',
	() => {
		const { details, pagination } = storeToRefs(useCommunityEmotesDetails())

		const data = computed<UserUsage[]>(() => {
			return details.value?.emotesStatisticEmoteDetailedInformation?.usagesByUsers ?? []
		})
		const total = computed(() => details.value?.emotesStatisticEmoteDetailedInformation?.usagesByUsersTotal ?? 0)
		const pageCount = computed(() => {
			return Math.ceil(total.value / pagination.value.pageSize)
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
						userId: row.original.userId,
						name: resolveUserName(row.original.twitchProfile.login, row.original.twitchProfile.displayName),
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
					return pagination.value
				},
			},
			manualPagination: true,
			onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
			getCoreRowModel: getCoreRowModel(),
			getPaginationRowModel: getPaginationRowModel(),
			getFacetedRowModel: getFacetedRowModel(),
			getFacetedUniqueValues: getFacetedUniqueValues(),
		})

		return {
			table,
			pagination,
			total,
			pageCount,
		}
	},
)
