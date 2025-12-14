<script setup lang="ts" generic="T extends RowData">
import { computed } from 'vue'

import type { PaginationState, RowData, Table } from '@tanstack/vue-table'

const props = defineProps<{
	total: number
	table: Table<T>
	pagination: PaginationState
}>()

const emits = defineEmits<{
	(event: 'update:page', page: number): void
	(event: 'update:pageSize', pageSize: number): void
}>()

const currentPage = computed(() => {
	if (props.pagination.pageIndex < 0) return 1
	return props.pagination.pageIndex + 1
})

function handleGoToPage(event: any) {
	const page = event.target.value ? Number(event.target.value) - 1 : 0
	if (Number.isNaN(page)) return
	emits('update:page', page < 0 ? 0 : page)
}

function handlePageSizeChange(pageSize: string) {
	emits('update:page', 0)
	emits('update:pageSize', Number(pageSize))
}
</script>

<template>
	<div class="flex justify-between max-sm:flex-col gap-4">
		<div class="flex gap-2 items-center">
			<div class="text-sm text-muted-foreground text-nowrap">
				{{ table.getPageCount() }} page(s) / {{ formatNumber(total) }} item(s)
			</div>
			<UiSelect
				default-value="10"
				:model-value="pagination.pageSize.toString()"
				@update:model-value="handlePageSizeChange"
			>
				<UiSelectTrigger class="h-9 justify-between gap-2">
					<div>
						Per page
						<UiSelectValue class="flex-none" />
					</div>
				</UiSelectTrigger>
				<UiSelectContent>
					<UiSelectItem
						v-for="pageSize in ['10', '20', '50', '100']"
						:key="pageSize"
						:value="pageSize"
					>
						{{ pageSize }}
					</UiSelectItem>
				</UiSelectContent>
			</UiSelect>
		</div>
		<div class="flex gap-2 items-center">
			<div class="flex gap-2 max-sm:justify-end max-sm:w-full">
				<UiButton
					class="w-9 h-9 min-w-9 max-sm:w-full"
					variant="outline"
					size="icon"
					:disabled="!table.getCanPreviousPage()"
					@click="table.previousPage()"
				>
					<Icon name="lucide:chevron-left" class="h-4 w-4" />
				</UiButton>
				<UiInput
					class="w-20 h-9 max-sm:w-full"
					:min="1"
					:max="table.getPageCount()"
					:model-value="currentPage"
					inputmode="numeric"
					type="number"
					:disabled="!table.getCanPreviousPage() && !table.getCanNextPage()"
					@input="handleGoToPage"
				/>
				<UiButton
					class="w-9 h-9 min-w-9 max-sm:w-full"
					variant="outline"
					size="icon"
					:disabled="!table.getCanNextPage()"
					@click="table.nextPage()"
				>
					<Icon name="lucide:chevron-right" class="h-4 w-4" />
				</UiButton>
			</div>
		</div>
	</div>
</template>
