<script setup lang="ts" generic="T extends RowData">
import type { PaginationState, RowData, Table } from '@tanstack/vue-table';
import { ChevronLeft, ChevronRight } from 'lucide-vue-next';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectTrigger, SelectValue, SelectItem } from '@/components/ui/select';
import { formatNumber } from '@/helpers/format-number.js';

const props = defineProps<{
	total: number
	table: Table<T>
	pagination: PaginationState
}>();

const currentPage = computed(() => {
	if (props.pagination.pageIndex < 0) return 1;
	return props.pagination.pageIndex + 1;
});

const emits = defineEmits<{
	(event: 'update:page', page: number): void
	(event: 'update:pageSize', pageSize: number): void
}>();

const { t } = useI18n();

function handleGoToPage(event: any) {
  const page = event.target.value ? Number(event.target.value) - 1 : 0;
	if (Number.isNaN(page)) return;
	emits('update:page', page < 0 ? 0 : page);
}

function handlePageSizeChange(pageSize: string) {
	emits('update:page', 0);
  emits('update:pageSize', Number(pageSize));
}
</script>

<template>
	<div class="flex items-center justify-between gap-2 py-4 w-full">
		<Input
			class="w-20 h-9"
			:min="1"
			:max="table.getPageCount()"
			:model-value="currentPage"
			inputmode="numeric"
			type="number"
			@input="handleGoToPage"
		/>
		<div class="flex-1 text-sm text-muted-foreground">
			{{ t('sharedTexts.pagination', {
				total: table.getPageCount(),
				items: formatNumber(total),
			}) }}
		</div>
		<div class="flex gap-2">
			<Button
				class="h-9 w-9"
				variant="outline"
				size="icon"
				:disabled="!table.getCanPreviousPage()"
				@click="table.previousPage()"
			>
				<ChevronLeft />
			</Button>
			<Select default-value="10" @update:model-value="handlePageSizeChange">
				<SelectTrigger class="w-20 h-9">
					<SelectValue />
				</SelectTrigger>
				<SelectContent>
					<SelectItem v-for="pageSize in ['10', '20', '50', '100']" :key="pageSize" :value="pageSize">
						{{ pageSize }}
					</SelectItem>
				</SelectContent>
			</Select>
			<Button
				class="h-9 w-9"
				variant="outline"
				size="icon"
				:disabled="!table.getCanNextPage()"
				@click="table.nextPage()"
			>
				<ChevronRight />
			</Button>
		</div>
	</div>
</template>
