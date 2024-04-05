<script setup lang="ts">
import {
	FlexRender,
} from '@tanstack/vue-table';
import { SearchIcon } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';

import { useNotificationsTable } from '../composables/use-notifications-table';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import ShadcnLayout from '@/layout/shadcn-layout.vue';

const { t } = useI18n();
const notificationsTable = useNotificationsTable();
</script>

<template>
	<div class="flex flex-wrap w-full items-center justify-between gap-2">
		<div class="flex gap-2 max-sm:w-full">
			<div class="relative w-full items-center">
				<Input id="search" type="text" placeholder="Search..." class="h-9 pl-10 max-sm:w-full" />
				<span class="absolute start-2 inset-y-0 flex items-center justify-center px-2">
					<SearchIcon class="size-4 text-muted-foreground" />
				</span>
			</div>
			<Select v-model="notificationsTable.notificationsFilter">
				<SelectTrigger class="h-9 w-[120px]">
					<SelectValue />
				</SelectTrigger>
				<SelectContent>
					<SelectGroup>
						<SelectItem value="globals">
							{{ t('adminPanel.notifications.globals') }}
						</SelectItem>
						<SelectItem value="users">
							{{ t('adminPanel.notifications.users') }}
						</SelectItem>
					</SelectGroup>
				</SelectContent>
			</Select>
		</div>

		<div class="flex items-center gap-2 max-sm:w-full">
			<div class="flex-1 text-sm text-muted-foreground">
				{{ t('sharedTexts.pagination', {
					page: notificationsTable.table.getState().pagination.pageIndex + 1,
					total: notificationsTable.table.getPageCount().toLocaleString(),
				}) }}
			</div>
			<Button
				variant="outline"
				size="sm"
				:disabled="!notificationsTable.table.getCanPreviousPage()"
				@click="notificationsTable.table.previousPage()"
			>
				{{ t('sharedButtons.previous') }}
			</Button>
			<Button
				variant="outline"
				size="sm"
				:disabled="!notificationsTable.table.getCanNextPage()"
				@click="notificationsTable.table.nextPage()"
			>
				{{ t('sharedButtons.next') }}
			</Button>
		</div>
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
		</Table>
	</shadcn-layout>
</template>
