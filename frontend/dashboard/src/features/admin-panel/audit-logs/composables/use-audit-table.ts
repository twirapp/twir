import {
	type ColumnDef,
	getCoreRowModel,
	getPaginationRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import { useAuditFilters } from './use-audit-filters.js'

import type { BadgeVariants } from '@/components/ui/badge'
import type { AdminAuditLogsInput, AdminAuditLogsQuery } from '@/gql/graphql'

import { useAdminAuditLogs } from '@/api/admin/audit-logs'
import { Badge } from '@/components/ui/badge'
import { usePagination } from '@/composables/use-pagination.js'
import AuditTableValue from '@/features/admin-panel/audit-logs/ui/audit-table-value.vue'
import UsersTableCellUser from '@/features/admin-panel/manage-users/ui/users-table-cell-user.vue'
import { AuditOperationType } from '@/gql/graphql'
import { valueUpdater } from '@/helpers/value-updater.js'

function computeOperationBadgeVariant(operation: AuditOperationType): BadgeVariants['variant'] {
	switch (operation) {
		case AuditOperationType.Create:
			return 'success'
		case AuditOperationType.Delete:
			return 'destructive'
		case AuditOperationType.Update:
			return 'default'
	}
}

export const useAuditTable = createGlobalState(() => {
	const { pagination } = usePagination()
	const { searchType, searchUserId, selectedFilters } = useAuditFilters()

	const auditParams = computed<AdminAuditLogsInput>(() => {
		const userSearchKey = searchType.value === 'channel' ? 'channelId' : 'userId'
		const params: AdminAuditLogsInput = {
			page: pagination.value.pageIndex,
			perPage: pagination.value.pageSize,
			[userSearchKey]: searchUserId.value || undefined,
			operationType: selectedFilters.value['operation-type'],
			system: selectedFilters.value.system,
		}
		return params
	})

	const { data, fetching } = useAdminAuditLogs(auditParams)

	const logs = computed<AdminAuditLogsQuery['adminAuditLogs']['logs']>(() => {
		if (!data.value) return []
		return data.value.adminAuditLogs.logs
	})

	const total = computed(() => data.value?.adminAuditLogs.total ?? 0)

	const pageCount = computed(() => {
		return Math.ceil(total.value / pagination.value.pageSize)
	})

	const tableColumns = computed<ColumnDef<AdminAuditLogsQuery['adminAuditLogs']['logs'][0]>[]>(() => [
		{
			accessorKey: 'channel',
			size: 15,
			minSize: 15,
			header: () => h('div', {}, { default: () => 'Channel' }),
			cell: ({ row }) => {
				if (!row.original.channel) {
					return null
				}

				return h('a', {
					class: 'flex flex-col',
					href: `https://twitch.tv/${row.original.channel.login}`,
					target: '_blank',
				}, h(UsersTableCellUser, {
					avatar: row.original.channel.profileImageUrl,
					name: row.original.channel.login,
					displayName: row.original.channel.displayName,
				}))
			},
		},
		{
			accessorKey: 'user',
			size: 15,
			minSize: 15,
			header: () => h('div', {}, { default: () => 'Actor' }),
			cell: ({ row }) => {
				if (!row.original.user) {
					return null
				}

				return h('a', {
					class: 'flex flex-col',
					href: `https://twitch.tv/${row.original.user.login}`,
					target: '_blank',
				}, h(UsersTableCellUser, {
					avatar: row.original.user.profileImageUrl,
					name: row.original.user.login,
					displayName: row.original.user.displayName,
				}))
			},
		},
		{
			accessorKey: 'system',
			size: 15,
			minSize: 15,
			header: () => h('div', {}, { default: () => 'System' }),
			cell: ({ row }) => {
				return row.original.system?.toLowerCase()
			},
		},
		{
			accessorKey: 'operationType',
			size: 15,
			minSize: 15,
			header: () => h('div', {}, { default: () => 'Operation' }),
			cell: ({ row }) => {
				return h(
					Badge,
					{ variant: computeOperationBadgeVariant(row.original.operationType) },
					{ default: () => row.original.operationType.toLowerCase() },
				)
			},
		},
		{
			accessorKey: 'value',
			size: 40,
			minSize: 40,
			header: () => h('div', {}, { default: () => 'Values' }),
			cell: ({ row }) => {
				return h(AuditTableValue, { log: row.original })
			},
		},
	])

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		get data() {
			return logs.value
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
		pagination,
		total,
		table,
		tableColumns,
	}
})
