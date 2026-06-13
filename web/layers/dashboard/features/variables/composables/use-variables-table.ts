import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import VariablesActions from '../ui/variables-actions.vue'

import type { CustomVariable } from '@/api/variables.js'

import { useVariablesApi } from '@/api/variables.js'
import { Badge } from '@/components/ui/badge/index.js'
import { VariableType } from '@/gql/graphql.js'

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
		{
			accessorKey: 'type',
			size: 10,
			header: () => t('variables.type'),
			cell: ({ row }) => h(Badge, { variant: 'secondary' }, row.original.type),
		},
		{
			accessorKey: 'response',
			size: 50,
			header: () => h('div', {}, t('sharedTexts.response')),
			cell: ({ row }) => {
				if (row.original.type === VariableType.Script) {
					return '<SCRIPT>'
				}

				return h('div', { class: 'truncate' }, row.original.response)
			},
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => h(VariablesActions, { row: row.original }),
		},
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
