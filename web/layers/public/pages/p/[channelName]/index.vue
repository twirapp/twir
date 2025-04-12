<script lang="ts" setup>
import { type Cell, type ColumnDef, FlexRender, getCoreRowModel, getExpandedRowModel, useVueTable } from '@tanstack/vue-table'

import { useCommands } from '~~/layers/public/api/use-commands'
import CommandsCooldownCell from '~~/layers/public/components/commands/commands-cooldown-cell.vue'
import CommandsNameCell from '~~/layers/public/components/commands/commands-name-cell.vue'
import CommandsPermissionsCell from '~~/layers/public/components/commands/commands-permissions-cell.vue'
import CommandsResponsesCell from '~~/layers/public/components/commands/commands-responses-cell.vue'
import { type Command, type Group, createGroups, isCommand } from '~~/layers/public/lib/commands'

definePageMeta({
	layout: 'public',
	alias: ['/p/:channelName/commands'],
})

const { data } = await useCommands()

const commandsWithGroups = computed(() => createGroups(data.value?.commandsPublic ?? []))

const columns: ColumnDef<Command | Group>[] = [
	{
		accessorKey: 'Name',
		size: 10,
		cell: ({ row }) => isCommand(row.original)
			? h(CommandsNameCell, {
				name: row.original.name,
				aliases: isCommand(row.original) ? row.original.aliases : [],
			})
			: h('div', {}, row.original.name),
	},
	{
		accessorKey: 'Response',
		size: 80,
		cell: ({ row }) => isCommand(row.original)
			? h(CommandsResponsesCell, {
				responses: row.original.responses,
				description: row.original.description,
			})
			: null,
	},
	{
		accessorKey: 'Permissions',
		size: 5,
		cell: ({ row }) => isCommand(row.original)
			? h(CommandsPermissionsCell, {
				permissions: row.original.permissions,
			})
			: null,
	},
	{
		accessorKey: 'Cooldown',
		size: 5,
		cell: ({ row }) => isCommand(row.original)
			? h(CommandsCooldownCell, {
				cooldown: row.original.cooldown,
				cooldownType: row.original.cooldownType,
			})
			: null,
	},
]

const table = useVueTable({
	get data() {
		return commandsWithGroups.value
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
	initialState: {
		expanded: true,
	},
})

function computeCellSpan(cell: Cell<Command | Group, unknown>) {
	const isGroup = !isCommand(cell.row.original)

	if (isGroup && cell.column.id === 'Name') {
		return columns.length
	}

	return 1
}
</script>

<template>
	<div class="flex-wrap w-full border rounded-md" style="background-color: rgb(24, 24, 28)">
		<UiTable>
			<UiTableHeader>
				<UiTableRow
					v-for="headerGroup in table.getHeaderGroups()"
					:key="headerGroup.id"
					class="text-slate-50"
				>
					<UiTableHead
						v-for="header in headerGroup.headers"
						:key="header.id"
						:style="{ width: `${header.column.columnDef.size}%` }"
					>
						<FlexRender
							v-if="!header.isPlaceholder"
							:render="header.column.columnDef.header"
							:props="header.getContext()"
						/>
					</UiTableHead>
				</UiTableRow>
			</UiTableHeader>

			<UiTableBody>
				<UiTableRow
					v-for="row in table.getRowModel().rows"
					:key="row.id"
				>
					<template
						v-for="cell in row.getVisibleCells()"
						:key="cell.id"
					>
						<UiTableCell
							v-if="isCommand(cell.row.original) || cell.column.id === 'Name'"
							:colspan="computeCellSpan(cell)"
						>
							<FlexRender
								:render="cell.column.columnDef.cell"
								:props="cell.getContext()"
								:class="{
									'flex items-center justify-center': !isCommand(cell.row.original),
								}"
							/>
						</UiTableCell>
					</template>
				</UiTableRow>
			</UiTableBody>
		</UiTable>
	</div>
</template>
