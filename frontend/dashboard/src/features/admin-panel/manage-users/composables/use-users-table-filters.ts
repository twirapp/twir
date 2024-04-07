import { refDebounced } from '@vueuse/core';
import { defineStore, storeToRefs } from 'pinia';
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useBadges } from '../../manage-badges/composables/use-badges';

interface Filter {
	group: string;
	list: {
		label: string;
		key: string;
		image?: string;
	}[];
}

export const useUsersTableFilters = defineStore('manage-users/users-table-filters', () => {
	const { t } = useI18n();

	const searchInput = ref('');
	const debounceSearchInput = refDebounced<string>(searchInput, 500);

	const { badges } = storeToRefs(useBadges());

	const filtersList = computed<Filter[]>(() => [
		{
			group: t('adminPanel.manageUsers.statusGroup'),
			list: [
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
			],
		},
		{
			group: t('adminPanel.manageUsers.badgesGroup'),
			list: badges.value.map((badge) => ({
				label: badge.name,
				key: badge.id,
				image: badge.fileUrl,
			})),
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
