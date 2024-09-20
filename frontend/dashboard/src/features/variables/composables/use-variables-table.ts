import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import type { CustomVariable } from '@/api/variables.js'

import { useVariablesApi } from '@/api/variables.js'

export const useVariablesTable = createGlobalState(() => {
	const { t } = useI18n()
	const variablesApi = useVariablesApi()

	const { data, fetching } = variablesApi.variablesQuery
	const variables = computed<CustomVariable[]>(() => {
		if (!data.value) return []
		return data.value.variables
	})

	const tableColumns = computed<ColumnDef<CustomVariable>[]>(() => [
		{
			accessorKey: 'name',
			size: 30,
			header: () => h('div', {}, t('sharedTexts.name')),
			cell: ({ row }) => h('span', row.original.name),
		},
		// {
		// 	accessorKey: 'actions',
		// 	size: 10,
		// 	header: () => '',
		// 	cell: ({ row }) => h(GreetingsTableActions, { greetings: row.original }),
		// },
	])

	const table = useVueTable({
		get data() {
			return variables.value
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
