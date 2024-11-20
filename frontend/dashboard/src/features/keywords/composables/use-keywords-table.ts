import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import KeywordsTableActions from '../ui/keywords-table-actions.vue'

import { type KeywordResponse, useKeywordsApi } from '@/api/keywords.js'

export const useKeywordsTable = createGlobalState(() => {
	const { t } = useI18n()
	const keywordsApi = useKeywordsApi()

	const { data, fetching } = keywordsApi.useQueryKeywords()
	const greetings = computed<KeywordResponse[]>(() => {
		if (!data.value) return []
		return data.value.keywords
	})

	const tableColumns = computed<ColumnDef<KeywordResponse>[]>(() => [
		{
			accessorKey: 'text',
			size: 20,
			header: () => h('div', {}, t('keywords.triggerText')),
			cell: ({ row }) => h('span', row.original.text),
		},
		{
			accessorKey: 'response',
			size: 20,
			header: () => h('div', {}, t('sharedTexts.response')),
			cell: ({ row }) => h('span', row.original.response ?? ''),
		},
		{
			accessorKey: 'variables',
			size: 45,
			header: () => h('div', {}, 'Variable for command'),
			cell: ({ row }) => h('span', { class: 'cursor-pointer text-xs' }, { default: () => `$(keywords.counter|${row.original.id})` }),
		},
		{
			accessorKey: 'usages',
			size: 5,
			header: () => h('div', {}, t('keywords.usages')),
			cell: ({ row }) => h('span', row.original.usageCount),
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => h(KeywordsTableActions, { keyword: row.original }),
		},
	])

	const table = useVueTable({
		get data() {
			return greetings.value
		},
		get columns() {
			return tableColumns.value
		},
		getCoreRowModel: getCoreRowModel(),
	})

	return {
		isLoading: fetching,
		table,
	}
})
