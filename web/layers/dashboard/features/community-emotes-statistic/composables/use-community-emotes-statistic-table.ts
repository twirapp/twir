import {
	type EmotesStatistics,
	useEmotesStatisticQuery,
} from '#layers/dashboard/api/emotes-statistic'
import { valueUpdater } from '#layers/dashboard/helpers/value-updater'
import {
	type ColumnDef,
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getPaginationRowModel,
	getSortedRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import { usePagination } from '~/composables/use-pagination'
import CommunityEmotesTableColumnEmote from '~/features/community-emotes-statistic/ui/community-emotes-table-column-emote.vue'
import { EmoteStatisticRange, type EmotesStatisticsOpts } from '~/gql/graphql'

import CommunityEmotesTableColumnChartRange from '../ui/community-emotes-table-column-chart-range.vue'
import CommunityEmotesTableColumnChart from '../ui/community-emotes-table-column-chart.vue'
import CommunityEmotesTableColumn from '../ui/community-emotes-table-column.vue'
import { useCommunityEmotesStatisticFilters } from './use-community-emotes-statistic-filters.js'

export const useCommunityEmotesStatisticTable = createGlobalState(() => {
	const { t } = useI18n()
	const { pagination } = usePagination()
	const { debouncedSearchInput, tableRange, sortingState, tableOrder } =
		useCommunityEmotesStatisticFilters()

	const emotesQueryOptions = computed<EmotesStatisticsOpts>((prevParams) => {
		if (prevParams?.search !== debouncedSearchInput.value) {
			pagination.value.pageIndex = 0
		}

		return {
			search: debouncedSearchInput.value,
			perPage: pagination.value.pageSize,
			page: pagination.value.pageIndex,
			graphicRange: tableRange.value,
			order: tableOrder.value,
		}
	})
	const { data, fetching } = useEmotesStatisticQuery(emotesQueryOptions)

	const emotes = computed<EmotesStatistics>(() => {
		if (!data.value) return []
		return data.value.emotesStatistics.emotes
	})
	const totalEmotes = computed(() => data.value?.emotesStatistics.total ?? 0)
	const pageCount = computed(() => {
		return Math.ceil(totalEmotes.value / pagination.value.pageSize)
	})

	const statsColumn = computed<ColumnDef<EmotesStatistics[0]>[]>(() => [
		{
			accessorKey: 'name',
			size: 5,
			header: () => h('div', {}, t('community.emotesStatistic.table.emote')),
			cell: ({ row }) => {
				return h(CommunityEmotesTableColumnEmote, { emoteName: row.original.emoteName })
			},
		},
		{
			accessorKey: 'usages',
			size: 5,
			header: ({ column }) => {
				return h(CommunityEmotesTableColumn, {
					column,
					title: t('community.emotesStatistic.table.usages'),
				})
			},
			cell: ({ row }) => {
				return h('div', `${row.original.totalUsages}`)
			},
		},
		{
			accessorKey: 'chart',
			size: 80,
			header: () => h(CommunityEmotesTableColumnChartRange),
			cell: ({ row }) => {
				return h(CommunityEmotesTableColumnChart, {
					isDayRange: tableRange.value === EmoteStatisticRange.LastDay,
					usages: row.original.graphicUsages,
				})
			},
		},
	])

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		get data() {
			return emotes.value
		},
		get columns() {
			return statsColumn.value
		},
		state: {
			get sorting() {
				return sortingState.value
			},
			get pagination() {
				return pagination.value
			},
		},
		manualPagination: true,
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		onSortingChange: (updaterOrValue) => valueUpdater(updaterOrValue, sortingState),
		getSortedRowModel: getSortedRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
	})

	return {
		isLoading: fetching,
		table,
		totalEmotes,
		pageCount,
		pagination,
	}
})
