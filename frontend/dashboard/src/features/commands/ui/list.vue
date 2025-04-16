<script setup lang="ts">
import {
	getCoreRowModel,
	getExpandedRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { type Rgb, colorBrightness, hexToRgb, rgbToHex } from '@zero-dependency/utils'
import { ChevronDownIcon, ChevronRightIcon } from 'lucide-vue-next'
import { computed, h } from 'vue'

import ColumnActions from './list-actions.vue'
import { type Group, createGroups, isCommand } from './list-groups.js'

import type { Command } from '@/gql/graphql'
import type { ColumnDef } from '@tanstack/vue-table'

import Table from '@/components/table.vue'

const props = withDefaults(defineProps<{
	commands: Command[]
	enableGroups?: boolean
	showBackground?: boolean
}>(), {
	showHeader: false,
	enableGroups: false,
})

const columns: ColumnDef<Command | Group>[] = [
	{
		accessorKey: 'name',
		size: 10,
		header: () => h('div', {}, 'Name'),
		cell: ({ row }) => {
			const chevron = row.getCanExpand() ? h(row.getIsExpanded() ? ChevronDownIcon : ChevronRightIcon) : null

			if (isCommand(row.original)) {
				return h(
					'div',
					{ class: 'flex gap-2 items-center' },
					[chevron, `!${row.getValue('name')}` as string],
				)
			}

			let rgbColor: Rgb | null = null
			if (row.original.color) {
				rgbColor = hexToRgb(rgbToHex(row.original.color))
			}

			const color = rgbColor
				? (colorBrightness(rgbColor) >= 128 ? '#000' : '#fff')
				: 'var(--n-text-color)'

			return h(
				'div',
				{ class: 'flex gap-2 items-center select-none' },
				[
					chevron,
					h(
						'span',
						{
							class: 'p-1 rounded',
							style: `background-color: ${row.original.color}; color: ${color}`,
						},
						row.original.name.charAt(0).toLocaleUpperCase() + row.original.name.slice(1),
					),
				],
			)
		},
	},
	{
		accessorKey: 'responses',
		header: () => h('div', { }, 'Responses'),
		size: 85,
		cell: ({ row }) => {
			if (!isCommand(row.original)) {
				return
			}

			const responses: Command['responses'] = row.getValue('responses')
			if (!responses?.length) {
				return row.original.description
			}

			const mappedResponses = responses.map((r) => h('span', { class: 'truncate md:whitespace-normal' }, r.text))
			return h('div', { class: 'flex flex-col' }, mappedResponses)
		},
	},
	{
		id: 'actions',
		size: 5,
		cell: ({ row }) => {
			if (!isCommand(row.original)) {
				return
			}

			return h(
				ColumnActions,
				{
					row: row.original,
				},
			)
		},
	},
]

const tableValue = computed(() => props.enableGroups ? createGroups(props.commands) : props.commands)

const table = useVueTable({
	get data() {
		return tableValue.value
	},
	get columns() {
		return columns
	},
	getCoreRowModel: getCoreRowModel(),
	getExpandedRowModel: getExpandedRowModel(),
	getSubRows: (original) => {
		if ('commands' in original) {
			return original.commands
		}
	},
})
</script>

<template>
	<Table :table="table" :is-loading="false" />
</template>
