<script lang="ts" setup>
import { type ColumnDef, getCoreRowModel, useVueTable, FlexRender } from '@tanstack/vue-table';
import type { Command } from '@twir/api/messages/commands_unprotected/commands_unprotected';
import { computed, h } from 'vue';

import TableRowsSkeleton from '@/components/TableRowsSkeleton.vue';
import {
	Table,
	TableBody,
	TableHead,
	TableHeader,
	TableRow,
	TableCell,
} from '@/components/ui/table';
import { useCommands } from '@/composables/use-commands';
import CommandsCooldownCell from '@/pages/commands/commands-cooldown-cell.vue';
import CommandsNameCell from '@/pages/commands/commands-name-cell.vue';
import CommandsPermissionsCell from '@/pages/commands/commands-permissions-cell.vue';
import CommandsResponsesCell from '@/pages/commands/commands-responses-cell.vue';

const { data, isLoading: isCommandsLoading } = useCommands();

const commands = computed(() => data.value?.commands ?? []);

const columns: ColumnDef<Command>[] = [
	{
		accessorKey: 'Name',
		size: 10,
		cell: ({ row }) => h(CommandsNameCell, {
			name: row.original.name,
			aliases: row.original.aliases,
		}),
	},
	{
		accessorKey: 'Response',
		size: 80,
		cell: ({ row }) => h(CommandsResponsesCell, {
			responses: row.original.responses,
			description: row.original.description,
		}),
	},
	{
		accessorKey: 'Permissions',
		size: 5,
		cell: ({ row }) => h(CommandsPermissionsCell, { permissions: row.original.permissions }),
	},
	{
		accessorKey: 'Cooldown',
		size: 5,
		cell: ({ row }) => h(CommandsCooldownCell, {
			cooldown: row.original.cooldown,
			cooldownType: row.original.cooldownType,
		}),
	},
];

const table = useVueTable({
	get data() {
		return commands.value;
	},
	get columns() {
		return columns;
	},
	getCoreRowModel: getCoreRowModel(),
});
</script>

<template>
	<div class="rounded-md border">
		<Table>
			<TableHeader>
				<TableRow
					v-for="headerGroup in table.getHeaderGroups()"
					:key="headerGroup.id"
					class="text-slate-50"
				>
					<TableHead
						v-for="header in headerGroup.headers"
						:key="header.id"
						:style="{ width: `${header.getSize()}%` }"
					>
						<FlexRender
							v-if="!header.isPlaceholder"
							:render="header.column.columnDef.header"
							:props="header.getContext()"
						/>
					</TableHead>
				</TableRow>
			</TableHeader>
			<Transition name="table-rows" appear mode="out-in">
				<TableBody v-if="isCommandsLoading">
					<table-rows-skeleton :rows="20" :colspan="4" />
				</TableBody>
				<TableBody v-else>
					<TableRow
						v-for="row in table.getRowModel().rows" :key="row.id"
					>
						<TableCell
							v-for="cell in row.getVisibleCells()"
							:key="cell.id"
						>
							<FlexRender
								:render="cell.column.columnDef.cell"
								:props="cell.getContext()"
							/>
						</TableCell>
					</TableRow>
				</TableBody>
			</Transition>
		</Table>
	</div>
</template>

<style scoped>
.table-rows-enter-active,
.table-rows-leave-active {
	transition: opacity 0.5s ease;
}

.table-rows-enter-from,
.table-rows-leave-to {
	opacity: 0;
}

.cooldownIcon {
	height: 18px;
	width: 18px;
}
</style>
