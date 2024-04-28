<script setup lang="ts">
import { IconChevronRight, IconChevronDown } from '@tabler/icons-vue';
import type { ColumnDef } from '@tanstack/vue-table';
import {
	FlexRender,
	getCoreRowModel,
	useVueTable,
	getExpandedRowModel,
} from '@tanstack/vue-table';
import { rgbToHex, hexToRgb, colorBrightness, type Rgb } from '@zero-dependency/utils';
import { useThemeVars } from 'naive-ui';
import { h, computed } from 'vue';

import EditModal from './edit-modal.vue';
import ColumnActions from './list-actions.vue';
import { type Group, isCommand, createGroups } from './list-groups.js';

import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import type { Command } from '@/gql/graphql';

const themeVars = useThemeVars();

const props = withDefaults(defineProps<{
	commands: Command[]
	enableGroups?: boolean,
	showBackground?: boolean
}>(), {
	showHeader: false,
	enableGroups: false,
});

const columns: ColumnDef<Command | Group>[] = [
  {
    accessorKey: 'name',
		size: 10,
    header: () => h('div', {}, 'Name'),
    cell: ({ row }) => {
			const chevron = row.getCanExpand() ? h(row.getIsExpanded() ? IconChevronDown : IconChevronRight) : null;

			if (isCommand(row.original)) {
				return h(
					'div',
					{ class: 'flex gap-2 items-center select-none' },
					[chevron, '!' + row.getValue('name') as string],
				);
			}

			let rgbColor: Rgb | null = null;
			if (row.original.color) {
				rgbColor = hexToRgb(rgbToHex(row.original.color));
			}

			const color = rgbColor
				? (colorBrightness(rgbColor) >= 128 ? '#000' : '#fff')
				: 'var(--n-text-color)';

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
			);
    },
  },
	{
    accessorKey: 'responses',
    header: () => h('div', {  }, 'Responses'),
		size: 85,
    cell: ({ row }) => {
			if (!isCommand(row.original)) {
				return;
			}

			const responses: Command['responses'] = row.getValue('responses');
			if (!responses?.length) {
				return row.original.description;
			}

			const mappedResponses = responses.map((r) => h('span', {}, r.text));
      return h('div', { class: 'flex flex-col' }, mappedResponses);
    },
  },
	{
		id: 'actions',
		size: 5,
    cell: ({ row }) => {
			if (!isCommand(row.original)) {
				return;
			}

			return h(
				ColumnActions,
				{
					row: row.original,
				},
			);
    },
  },
];

const tableValue = computed(() => props.enableGroups ? createGroups(props.commands) : props.commands);

const table = useVueTable({
	get data() {
		return tableValue.value;
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
</script>

<template>
	<edit-modal />

	<div
		class="border rounded-md"
		:style="{
			backgroundColor: props.showBackground ? themeVars.cardColor : 'inherit',
			color: themeVars.textColor2
		}"
	>
		<Table>
			<TableHeader>
				<TableRow
					v-for="headerGroup in table.getHeaderGroups()"
					:key="headerGroup.id"
					class="border-b"
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
			<TableBody>
				<template v-if="table.getRowModel().rows?.length">
					<TableRow
						v-for="row in table.getRowModel().rows" :key="row.id"
						:data-state="row.getIsSelected() ? 'selected' : undefined"
						class="border-b"
						:class="{ 'cursor-pointer': !isCommand(row.original) }"
					>
						<TableCell
							v-for="cell in row.getVisibleCells()"
							:key="cell.id"
							@click="() => {
								if (row.getCanExpand()) {
									row.getToggleExpandedHandler()()
								}
							}"
						>
							<FlexRender
								:render="cell.column.columnDef.cell"
								:props="cell.getContext()"
							/>
						</TableCell>
					</TableRow>
				</template>
				<template v-else>
					<TableRow>
						<TableCell :colSpan="columns.length" class="h-24 text-center">
							No commands
						</TableCell>
					</TableRow>
				</template>
			</TableBody>
		</Table>
	</div>
</template>
