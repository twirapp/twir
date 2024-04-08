import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

export const useNotificationsFilters = defineStore('admin-panel/notifications-filters', () => {
	const searchInput = ref('');

	// globals or users
	const filterInput = ref('globals');
	const isUsersFilter = computed(() => filterInput.value === 'users');

	return {
		searchInput,
		filterInput,
		isUsersFilter,
	};
});
