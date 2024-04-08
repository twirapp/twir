import type { PaginationState } from '@tanstack/vue-table';
import { ref } from 'vue';

export const usePagination = () => {
	const pagination = ref({
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
