<script lang="ts" setup>
import {
	type ColumnDef,
	getCoreRowModel,
	useVueTable,
	FlexRender,
	getExpandedRowModel, type Cell,
} from '@tanstack/vue-table';
import type { Command } from '@twir/api/messages/commands_unprotected/commands_unprotected';
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import { computed, h, onMounted } from 'vue';

import { useCommands } from '@/api/use-commands';
import TableRowsSkeleton from '@/components/TableRowsSkeleton.vue';
import { Badge } from '@/components/ui/badge';
import {
	Table,
	TableBody,
	TableHead,
	TableHeader,
	TableRow,
	TableCell,
} from '@/components/ui/table';
import CommandsCooldownCell from '@/pages/commands/commands-cooldown-cell.vue';
import CommandsNameCell from '@/pages/commands/commands-name-cell.vue';
import CommandsPermissionsCell, {
	permissionsIconsMapping,
} from '@/pages/commands/commands-permissions-cell.vue';
import CommandsResponsesCell from '@/pages/commands/commands-responses-cell.vue';
import { createGroups, type Group, isCommand } from '@/pages/commands/create-group';

const { data } = useCommands();

const commandsWithGroups = computed(() => createGroups(data.value?.commands ?? []));

const columns: ColumnDef<Command | Group>[] = [
	{
		accessorKey: 'Name',
		size: 10,
		cell: ({ row }) => isCommand(row.original) ? h(CommandsNameCell, {
			name: row.original.name,
			aliases: isCommand(row.original) ? row.original.aliases : [],
		}) : h('div', {}, row.original.name),
	},
	{
		accessorKey: 'Response',
		size: 80,
		cell: ({ row }) => isCommand(row.original) ? h(CommandsResponsesCell, {
			responses: row.original.responses,
			description: row.original.description,
		}) : null,
	},
	{
		accessorKey: 'Permissions',
		size: 5,
		cell: ({ row }) => isCommand(row.original) ? h(CommandsPermissionsCell, {
			permissions: row.original.permissions,
		}) : null,
	},
	{
		accessorKey: 'Cooldown',
		size: 5,
		cell: ({ row }) => isCommand(row.original) ? h(CommandsCooldownCell, {
			cooldown: row.original.cooldown,
			cooldownType: row.original.cooldownType,
		}) : null,
	},
];

const table = useVueTable({
	get data() {
		return commandsWithGroups.value;
	},
	get columns() {
		return columns;
	},
	getCoreRowModel: getCoreRowModel(),
	getExpandedRowModel: getExpandedRowModel(),
	getSubRows: (original) => {
		if ('commands' in original) {
			return original.commands;
		}
	},
});

onMounted(() => {
	table.toggleAllRowsExpanded();
});

function computeCellSpan(cell: Cell<Command | Group, unknown>) {
	const isGroup = !isCommand(cell.row.original);

	if (isGroup && cell.column.id === 'Name') {
		return columns.length;
	}

	return 1;
}

const breakpoints = useBreakpoints(breakpointsTailwind);
const isSmall = breakpoints.smaller('xl');
</script>

<template>
	<div v-if="isSmall" class="flex flex-col gap-2">
		<div
			v-for="command of data?.commands"
			:key="command.name"
			class="flex flex-col gap-2 rounded-md border p-5"
		>
			<h3>
				!{{ command.name }}
			</h3>

			<span v-if="command.description" class="text-sm text-muted-foreground break-all">{{ command.description }}</span>
			<template v-else>
				<span
					v-for="(r, idx) of command.responses" :key="idx"
					class="text-sm text-muted-foreground break-all"
				>{{ r }}</span>
			</template>
			<div class="flex flex-wrap gap-1">
				<Badge v-if="command.group || command.module !== 'CUSTOM'" variant="secondary">
					{{ command.group || command.module }}
				</Badge>
				<Badge variant="secondary">
					Cooldown | {{ command.cooldown }}s
				</Badge>
				<Badge v-for="perm of command.permissions" :key="perm.name" variant="secondary">
					<div class="flex items-center gap-1">
						<component :is="permissionsIconsMapping[perm.type]" class="w-4 h-4" />
						{{ perm.name.charAt(0).toUpperCase() + perm.name.slice(1).toLowerCase() }}
					</div>
				</Badge>
			</div>
		</div>
	</div>

	<div v-else class="rounded-md border">
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
						:style="{ width: `${header.column.columnDef.size}%` }"
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
				<TableBody v-if="!data">
					<table-rows-skeleton
						:rows="20"
						:colspan="table.getAllColumns().length"
					/>
				</TableBody>
				<TableBody v-else>
					<TableRow
						v-for="row in table.getRowModel().rows"
						:key="row.id"
					>
						<template
							v-for="cell in row.getVisibleCells()"
							:key="cell.id"
						>
							<TableCell
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
							</TableCell>
						</template>
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
</style>
