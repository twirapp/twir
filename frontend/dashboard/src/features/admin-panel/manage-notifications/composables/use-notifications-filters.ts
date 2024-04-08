import { refDebounced } from '@vueuse/core';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

export const useNotificationsFilters = defineStore('admin-panel/notifications-filters', () => {
	const searchInput = ref('');
	const debounceSearchInput = refDebounced<string>(searchInput, 500);

	// globals or users
	const filterInput = ref('globals');
	const isUsersFilter = computed(() => filterInput.value === 'users');

	return {
		searchInput,
		debounceSearchInput,
		filterInput,
		isUsersFilter,
	};
});
