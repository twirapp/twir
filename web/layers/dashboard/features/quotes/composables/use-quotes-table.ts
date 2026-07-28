import {
	type ColumnDef,
	getCoreRowModel,
	getPaginationRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { intlFormat } from 'date-fns'
import { computed, h, ref, watch } from 'vue'

import type { Quote } from '~~/layers/dashboard/api/quotes.js'

import { useQuotesApi } from '~~/layers/dashboard/api/quotes.js'
import { usePagination } from '~~/layers/dashboard/composables/use-pagination.js'
import { valueUpdater } from '~~/layers/dashboard/helpers/value-updater.js'

import QuotesTableActions from '../ui/quotes-table-actions.vue'

export const useQuotesTable = createGlobalState(() => {
	const { t } = useI18n()
	const quotesApi = useQuotesApi()
	const searchInput = ref('')
	const { pagination } = usePagination()
	const { data, fetching } = quotesApi.useQueryQuotes()
	const quotes = computed<Quote[]>(() => data.value?.quotes ?? [])
	const filteredQuotes = computed(() => {
		const search = searchInput.value.trim().toLocaleLowerCase()
		if (!search) return quotes.value

		return quotes.value.filter((quote) => quote.text.toLocaleLowerCase().includes(search))
	})
	const pageCount = computed(() => Math.ceil(filteredQuotes.value.length / pagination.value.pageSize))

	watch(searchInput, () => {
		pagination.value.pageIndex = 0
	})

	watch(pageCount, (count) => {
		const lastPageIndex = Math.max(count - 1, 0)
		if (pagination.value.pageIndex > lastPageIndex) {
			pagination.value.pageIndex = lastPageIndex
		}
	})

	const tableColumns = computed<ColumnDef<Quote>[]>(() => [
		{
			accessorKey: 'number',
			size: 8,
			header: () => h('div', {}, t('quotes.number')),
			cell: ({ row }) => h('span', `#${row.original.number}`),
		},
		{
			accessorKey: 'text',
			size: 40,
			header: () => h('div', {}, t('quotes.text')),
			cell: ({ row }) => h('span', {
				class: 'block max-w-xl truncate',
				title: row.original.text,
			}, row.original.text),
		},
		{
			accessorKey: 'creatorName',
			size: 16,
			header: () => h('div', {}, t('quotes.creator')),
			cell: ({ row }) => h('span', row.original.creatorName ?? '-'),
		},
		{
			accessorKey: 'gameName',
			size: 16,
			header: () => h('div', {}, t('quotes.game')),
			cell: ({ row }) => h('span', row.original.gameName ?? '-'),
		},
		{
			accessorKey: 'createdAt',
			size: 16,
			header: () => h('div', {}, t('quotes.createdAt')),
			cell: ({ row }) => h('span', intlFormat(new Date(row.original.createdAt), {
				day: 'numeric',
				hour: 'numeric',
				minute: 'numeric',
				month: 'numeric',
				year: 'numeric',
			})),
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => h(QuotesTableActions, { quote: row.original }),
		},
	])

	const table = useVueTable({
		get data() {
			return filteredQuotes.value
		},
		get columns() {
			return tableColumns.value
		},
		state: {
			get pagination() {
				return pagination.value
			},
		},
		onPaginationChange: updaterOrValue => valueUpdater(updaterOrValue, pagination),
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
	})

	return {
		filteredQuotes,
		isLoading: fetching,
		pagination,
		searchInput,
		table,
	}
})
