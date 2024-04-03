<script setup lang="ts">
import {
	type ColumnDef,
	FlexRender,
	getCoreRowModel,
	useVueTable,
} from '@tanstack/vue-table';
import { addZero } from '@zero-dependency/utils';
import { useThemeVars } from 'naive-ui';
import { computed, h } from 'vue';

import ActionsButton from './actions-button.vue';
import CreatedAtCell from './created-at-cell.vue';

import { Button } from '@/components/ui/button';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';

interface Notification {
	id: number
	message: string
	url?: string
	createdAt: Date
}

const themeVars = useThemeVars();

const columns: ColumnDef<Notification>[] = [
	{
		accessorKey: 'message',
		size: 90,
		header: () => h('div', {}, 'Message'),
		cell: ({ row }) => {
			return h('span', row.original.message);
		},
	},
	{
		accessorKey: 'createdAt',
		size: 5,
		header: () => h('div', {}, 'Created At'),
		cell: ({ row }) => {
			return h(CreatedAtCell, { time: row.original.createdAt });
		},
	},
	{
		accessorKey: 'actions',
		size: 5,
		header: () => '',
		cell: ({ row }) => {
			return h(ActionsButton, { notificationId: row.original.id });
		},
	},
];

const tableNotifications = computed<Notification[]>(() => {
	return Array.from({ length: 20 }, (_, id) => ({
		id: id++,
		message: 'orem Ipsum - это текст-"рыба", часто используемый в печати и вэб-дизайне. Lorem Ipsum является стандартной "рыбой" для текстов на латинице с начала XVI века.',
		url: 'https://twir.app',
		createdAt: new Date(`2024-04-02T${id < 10 ? addZero(id) : id}:48:40.096Z`),
	}));
});

const table = useVueTable({
	get data() {
		return tableNotifications.value;
	},
	get columns() {
		return columns;
	},
	getCoreRowModel: getCoreRowModel(),
});
</script>

<template>
	<div
		class="flex-wrap border rounded-md" :style="{
			backgroundColor: themeVars.cardColor,
			color: themeVars.textColor2
		}"
	>
		<Table>
			<TableHeader>
				<TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id" class="border-b">
					<TableHead v-for="header in headerGroup.headers" :key="header.id" :style="{ width: `${header.getSize()}%` }">
						<FlexRender
							v-if="!header.isPlaceholder" :render="header.column.columnDef.header"
							:props="header.getContext()"
						/>
					</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				<template v-if="table.getRowModel().rows?.length">
					<TableRow
						v-for="row in table.getRowModel().rows" :key="row.id"
						:data-state="row.getIsSelected() ? 'selected' : undefined" class="border-b"
					>
						<TableCell
							v-for="cell in row.getVisibleCells()" :key="cell.id" @click="() => {
								if (row.getCanExpand()) {
									row.getToggleExpandedHandler()()
								}
							}"
						>
							<FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
						</TableCell>
					</TableRow>
				</template>
				<template v-else>
					<TableRow>
						<TableCell :colSpan="columns.length" class="h-24 text-center">
							No streamers
						</TableCell>
					</TableRow>
				</template>
			</TableBody>
		</Table>
	</div>

	<div class="flex w-full items-center justify-end space-x-2 py-4">
		<div class="flex-1 text-sm text-muted-foreground">
			{{ table.getFilteredSelectedRowModel().rows.length }} of
			{{ table.getFilteredRowModel().rows.length }} row(s) selected.
		</div>
		<div class="space-x-2">
			<Button
				variant="outline"
				size="sm"
				:disabled="!table.getCanPreviousPage()"
				@click="table.previousPage()"
			>
				Previous
			</Button>
			<Button
				variant="outline"
				size="sm"
				:disabled="!table.getCanNextPage()"
				@click="table.nextPage()"
			>
				Next
			</Button>
		</div>
	</div>
</template>
