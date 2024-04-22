<script setup lang="ts">
import {
	FlexRender,
} from '@tanstack/vue-table';
import { useI18n } from 'vue-i18n';

import NotificationsTableSearch from './notifications-table-search.vue';
import { useNotificationsTable } from '../composables/use-notifications-table';

import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { TooltipProvider } from '@/components/ui/tooltip';
import ShadcnLayout from '@/layout/shadcn-layout.vue';

const { t } = useI18n();
const notificationsTable = useNotificationsTable();
</script>

<template>
	<div class="flex flex-wrap w-full items-center justify-between gap-2">
		<notifications-table-search />
		<slot name="pagination" />
	</div>

	<shadcn-layout>
		<Table>
			<TableHeader>
				<TableRow v-for="headerGroup in notificationsTable.table.getHeaderGroups()" :key="headerGroup.id" class="border-b">
					<TableHead v-for="header in headerGroup.headers" :key="header.id" :style="{ width: `${header.getSize()}%` }">
						<FlexRender
							v-if="!header.isPlaceholder" :render="header.column.columnDef.header"
							:props="header.getContext()"
						/>
					</TableHead>
				</TableRow>
			</TableHeader>
			<TooltipProvider :delay-duration="100">
				<TableBody>
					<template v-if="notificationsTable.table.getRowModel().rows?.length">
						<TableRow
							v-for="row in notificationsTable.table.getRowModel().rows" :key="row.id"
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
							<TableCell :colSpan="notificationsTable.tableColumns.length" class="h-24 text-center">
								{{ t('adminPanel.notifications.emptyNotifications') }}
							</TableCell>
						</TableRow>
					</template>
				</TableBody>
			</TooltipProvider>
		</Table>
	</shadcn-layout>

	<slot name="pagination" />
</template>
