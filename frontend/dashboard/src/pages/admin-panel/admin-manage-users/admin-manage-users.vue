<script setup lang="ts">
import {
	type ColumnDef,
	FlexRender,
	getCoreRowModel,
	useVueTable,
} from '@tanstack/vue-table';
import { GetTwirStreamersResponse_Streamer as Streamer } from '@twir/api/messages/stats/stats';
import { SearchIcon } from 'lucide-vue-next';
import { useThemeVars } from 'naive-ui';
import { computed, h } from 'vue';

import ActionsButton from './actions-button.vue';

import { useStreamers } from '@/api/streamers';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';

const themeVars = useThemeVars();

const { data: streamers } = useStreamers();

const columns: ColumnDef<Streamer>[] = [
	{
		accessorKey: 'userDisplayName',
		size: 40,
		header: () => h('div', {}, 'Name'),
		cell: ({ row }) => {
			return h('a',
				{
					class: 'flex items-center gap-4 flex-wrap max-sm:justify-center',
					href: `https://twitch.tv/${row.original.userLogin}`,
					target: '_blank',
				},
				[
					h('img', { class: 'h-9 w-9', src: row.original.avatar, loading: 'lazy' }),
					row.original.userDisplayName,
				],
			);
		},
	},
	{
		accessorKey: 'userId',
		size: 25,
		header: () => h('div', {}, 'Id'),
		cell: ({ row }) => {
			return h('span', row.original.userId);
		},
	},
	{
		accessorKey: 'followersCount',
		enableSorting: true,
		size: 25,
		header: () => h('div', {}, 'Followers'),
		// header: ({ column }) => {
		//   return h(Button, {
		//     variant: 'ghost',
		//     onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
		//   }, () => ['Followers', h(ArrowUpDown, { class: 'ml-2 h-4 w-4' })]);
		// },
		cell: ({ row }) => {
			return h('span', row.original.followersCount);
		},
	},
	{
		accessorKey: 'actions',
		size: 10,
		header: () => '',
		cell: () => {
			return h(ActionsButton);
		},
	},
];

const tableStreamers = computed(() => {
	if (!streamers.value) return [];
	// TODO: remove mock
	return Array.from({ length: 10 }, () => streamers.value.streamers[0]);
});

const table = useVueTable({
	get data() {
		return tableStreamers.value;
	},
	get columns() {
		return columns;
	},
	getCoreRowModel: getCoreRowModel(),
});
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<div class="relative w-full items-center">
			<Input id="search" type="text" placeholder="Search..." class="pl-10" />
			<span class="absolute start-2 inset-y-0 flex items-center justify-center px-2">
				<SearchIcon class="size-4 text-muted-foreground" />
			</span>
		</div>
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

		<div class="flex items-center justify-end space-x-2 py-4">
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
	</div>
</template>
