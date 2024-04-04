<script setup lang="ts">
import {
	type ColumnDef,
	FlexRender,
	getCoreRowModel,
	useVueTable,
} from '@tanstack/vue-table';
import { Notification } from '@twir/api/messages/admin_notifications/admin_notifications';
import { useThemeVars } from 'naive-ui';
import { computed, h } from 'vue';
import { useI18n } from 'vue-i18n';

import ActionsButton from './actions-button.vue';
import CreatedAtCell from './created-at-cell.vue';

import { useAdminNotifications } from '@/api/notifications';
import { Button } from '@/components/ui/button';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';

const { t } = useI18n();
const themeVars = useThemeVars();

const notificationsCrud = useAdminNotifications();
const notifications = notificationsCrud.getAll({});
console.log(notificationsCrud.getAll);

const columns: ColumnDef<Notification>[] = [
	{
		accessorKey: 'message',
		size: 80,
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
			return h(CreatedAtCell, { time: new Date(row.original.createdAt) });
		},
	},
	{
		accessorKey: 'actions',
		size: 10,
		header: () => '',
		cell: ({ row }) => {
			return h(ActionsButton, {
				onDelete: () => onDeleteNotification(row.original.id),
				onEdit: () => onEditNotification(row.original.id),
			});
		},
	},
];

const notificationsList = computed<Notification[]>(() => {
	return notifications.data.value?.notifications ?? [];
});

const table = useVueTable({
	get data() {
		return notificationsList.value;
	},
	get columns() {
		return columns;
	},
	getCoreRowModel: getCoreRowModel(),
});

async function onDeleteNotification(notificationId: string) {
	console.log(notificationId);
}

async function onEditNotification(notificationId: string) {

}
</script>

<template>
	<div
		class="flex-wrap w-full border rounded-md" :style="{
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
							{{ t('adminPanel.notifications.emptyNotifications') }}
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
