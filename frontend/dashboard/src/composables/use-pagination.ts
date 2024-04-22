import type { PaginationState } from '@tanstack/vue-table';
import { ref, type Ref } from 'vue';

export interface UsePagination {
	pagination: Ref<PaginationState>;
	setPagination(state: PaginationState): PaginationState;
}

export const usePagination = (): UsePagination => {
	const pagination = ref<PaginationState>({
		pageIndex: 0,
		pageSize: 10,
	});

	function setPagination({
		pageIndex,
		pageSize,
	}: PaginationState): PaginationState {
		pagination.value.pageIndex = pageIndex;
		pagination.value.pageSize = pageSize;
		return { pageIndex, pageSize };
	}

	return {
		pagination,
		setPagination,
	};
};
