<script setup lang="ts">
import { FlexRender } from '@tanstack/vue-table';
import { useI18n } from 'vue-i18n';

import { useUsersTable } from '../composables/use-users-table.js';

import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';

const { t } = useI18n();
const usersTable = useUsersTable();
</script>

<template>
	<Table>
		<TableHeader>
			<TableRow v-for="headerGroup in usersTable.table.getHeaderGroups()" :key="headerGroup.id" class="border-b">
				<TableHead v-for="header in headerGroup.headers" :key="header.id" :style="{ width: `${header.getSize()}%` }">
					<FlexRender
						v-if="!header.isPlaceholder" :render="header.column.columnDef.header"
						:props="header.getContext()"
					/>
				</TableHead>
			</TableRow>
		</TableHeader>
		<TableBody :class="[usersTable.isLoading ? 'animate-pulse' : '']">
			<template v-if="usersTable.table.getRowModel().rows?.length">
				<!-- TODO: highlight banned users -->
				<TableRow
					v-for="row in usersTable.table.getRowModel().rows" :key="row.id"
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
					<TableCell :colSpan="usersTable.tableColumns.length" class="h-24 text-center">
						{{ t('adminPanel.manageUsers.noUsers') }}
					</TableCell>
				</TableRow>
			</template>
		</TableBody>
	</Table>
</template>
