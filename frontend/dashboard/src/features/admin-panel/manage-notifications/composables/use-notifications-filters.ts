import { refDebounced } from '@vueuse/core';
import { defineStore } from 'pinia';
import { ref } from 'vue';

import { NotificationType } from '@/gql/graphql';

export const useNotificationsFilters = defineStore('admin-panel/notifications-filters', () => {
	const searchInput = ref('');
	const debounceSearchInput = refDebounced<string>(searchInput, 500);

	const filterInput = ref<NotificationType>(NotificationType.Global);

	return {
		searchInput,
		debounceSearchInput,
		filterInput,
	};
});
