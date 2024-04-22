import { refDebounced } from '@vueuse/core';
import { defineStore } from 'pinia';
import { ref } from 'vue';

import { CommunityUsersOrder, CommunityUsersSortBy } from '@/gql/graphql';

export const useCommunityTableActions = defineStore('features/community-actions', () => {
	const searchInput = ref('');
	const debouncedSearchInput = refDebounced(searchInput, 500);

	const tableOrder = ref<CommunityUsersOrder>(CommunityUsersOrder.Desc);
	function setOrder(order: CommunityUsersOrder) {
		tableOrder.value = order;
	}

	const tableSortBy = ref<CommunityUsersSortBy>(CommunityUsersSortBy.Watched);
	function setTableSortBy(sortBy: CommunityUsersSortBy) {
		tableSortBy.value = sortBy;
	}

	const selectedColumns = ref<string[]>([]);

	return {
		searchInput,
		debouncedSearchInput,

		tableOrder,
		setOrder,

		tableSortBy,
		setTableSortBy,
	};
});
