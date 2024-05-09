import { type ColumnDef, getCoreRowModel, getPaginationRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import GreetingsTableActions from '../ui/greetings-table-actions.vue'

import { type Greetings, useGreetingsApi } from '@/api/greetings.js'
import { usePagination } from '@/composables/use-pagination.js'
import UsersTableCellUser from '@/features/admin-panel/manage-users/components/users-table-cell-user.vue'
import { valueUpdater } from '@/helpers/value-updater.js'

export const useGreetingsTable = createGlobalState(() => {
	const { t } = useI18n()
	const greetingsApi = useGreetingsApi()
	const { pagination } = usePagination()

	const { data, fetching } = greetingsApi.useQueryGreetings()
	const greetings = computed<Greetings[]>(() => {
		if (!data.value) return []
		return data.value.greetings
	})

	const pageCount = computed(() => {
		return Math.ceil(greetings.value.length / pagination.value.pageSize)
	})

	const tableColumns = computed<ColumnDef<Greetings>[]>(() => [
		{
			accessorKey: 'user',
			size: 60,
			header: () => h('div', {}, t('sharedTexts.user')),
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
			accessorKey: 'text',
			size: 30,
			header: () => h('div', {}, t('sharedTexts.response')),
			cell: ({ row }) => h('span', row.original.text),
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => h(GreetingsTableActions, { greetings: row.original }),
		},
	])

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		get data() {
			return greetings.value
		},
		get columns() {
			return tableColumns.value
		},
		state: {
			get pagination() {
				return pagination.value
			},
		},
		manualPagination: true,
		enableRowSelection: true,
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
	})

	return {
		isLoading: fetching,
		table,
	}
})
