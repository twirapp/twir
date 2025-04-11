import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import type { ChatWall } from '@/api/moderation-chat-wall.ts'

import { useModerationChatWall } from '@/api/moderation-chat-wall.ts'
import ChatWallAction from '@/features/moderation/ui/chat-wall-action.vue'
import ChatWallLog from '@/features/moderation/ui/chat-wall-log.vue'

export const useModerationWallTable = createGlobalState(() => {
	const { t } = useI18n()
	const api = useModerationChatWall()
	const { data, fetching } = api.useList()

	const list = computed(() => {
		return data.value?.chatWalls ?? []
	})

	const tableColumns = computed<ColumnDef<ChatWall>[]>(() => [
		{
			accessorKey: 'phrase',
			size: 20,
			header: () => t('chatWall.table.phrase'),
			cell: ({ row }) => {
				return h('span', row.original.phrase)
			},
		},
		{
			accessorKey: 'createdAt',
			size: 20,
			header: () => t('chatWall.table.createdAt'),
			cell: ({ row }) => {
				return h('span', new Date(row.original.createdAt).toLocaleString())
			},
		},
		{
			accessorKey: 'enabled',
			size: 5,
			header: () => t('chatWall.table.inProcess'),
			cell: ({ row }) => {
				return h('span', row.original.enabled ? t('chatWall.table.yes') : t('chatWall.table.no'))
			},
		},
		{
			accessorKey: 'action',
			size: 30,
			header: () => t('chatWall.table.action'),
			cell: ({ row }) => {
				return h(ChatWallAction, { chatWall: row.original })
			},
		},
		{
			accessorKey: 'info',
			header: () => '',
			size: 10,
			cell: ({ row }) => {
				return h(ChatWallLog, { chatWall: row.original })
			},
		},
	])

	const table = useVueTable({
		get data() {
			return list.value
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
