import { refDebounced } from '@vueuse/core';
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

export const useUsersTableFilters = defineStore('manage-users/users-table-filters', () => {
	const { t } = useI18n();

	const searchInput = ref('');
	const debounceSearchInput = refDebounced<string>(searchInput, 500);

	const filtersList = computed(() => [
		{
			label: t('adminPanel.manageUsers.isAdmin'),
			key: 'isAdmin',
		},
		{
			label: t('adminPanel.manageUsers.isBanned'),
			key: 'isBanned',
		},
		{
			label: t('adminPanel.manageUsers.isBotEnabled'),
			key: 'isBotEnabled',
		},
	]);

	const selectedFilters = ref<Record<string, true | undefined>>({});
	const selectedFiltersCount = computed(() => {
		return Object.keys(selectedFilters.value).length;
	});

	function clearFilters() {
		selectedFilters.value = {};
	}

	function setFilterValue(filterKey: string) {
		if (selectedFilters.value[filterKey]) {
			delete selectedFilters.value[filterKey];
			return;
		}

		selectedFilters.value[filterKey] = true;
	}

	return {
		searchInput,
		debounceSearchInput,
		filtersList,
		selectedFilters,
		selectedFiltersCount,
		setFilterValue,
		clearFilters,
	};
});
