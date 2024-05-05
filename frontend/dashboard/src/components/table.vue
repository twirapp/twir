<script setup lang="ts" generic="T extends RowData">
import { FlexRender, type RowData, type Table } from '@tanstack/vue-table'

import {
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	Table as TableRoot,
	TableRow,
} from '@/components/ui/table'
import ShadcnLayout from '@/layout/shadcn-layout.vue'

defineProps<{
	table: Table<T>
	isLoading: boolean
	hideHeader?: boolean
}>()
</script>

<template>
	<ShadcnLayout>
		<TableRoot>
			<TableHeader v-if="!hideHeader">
				<TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id" class="border-b">
					<TableHead v-for="header in headerGroup.headers" :key="header.id" :style="{ width: `${header.getSize()}%` }">
						<FlexRender
							v-if="!header.isPlaceholder" :render="header.column.columnDef.header"
							:props="header.getContext()"
						/>
					</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody :class="[isLoading ? 'animate-pulse' : '']">
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
						<TableCell :colSpan="table.getAllColumns().length" class="h-24 text-center">
							<slot name="empty-message" />
						</TableCell>
					</TableRow>
				</template>
			</TableBody>
		</TableRoot>
	</ShadcnLayout>
</template>
