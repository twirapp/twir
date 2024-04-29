import { type Ref, ref } from 'vue'

import type { PaginationState } from '@tanstack/vue-table'

export interface UsePagination {
	pagination: Ref<PaginationState>
	setPagination: (state: PaginationState) => PaginationState
}

export function usePagination(): UsePagination {
	const pagination = ref<PaginationState>({
		pageIndex: 0,
		pageSize: 10
	})

	function setPagination({
		pageIndex,
		pageSize
	}: PaginationState): PaginationState {
		pagination.value.pageIndex = pageIndex
		pagination.value.pageSize = pageSize
		return { pageIndex, pageSize }
	}

	return {
		pagination,
		setPagination
	}
}
