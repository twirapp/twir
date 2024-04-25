import type { ColumnFiltersState, SortingState, VisibilityState } from '@tanstack/vue-table';
import { refDebounced, useLocalStorage } from '@vueuse/core';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { TABLE_ACCESSOR_KEYS } from './use-community-users-table';

import { CommunityUsersOrder, CommunityUsersSortBy } from '@/gql/graphql';

const COLUMN_VISIBLE_STORAGE_KEY = 'twirCommunityUsersColumnVisibility';

export const useCommunityTableActions = defineStore('features/community-table-actions', () => {
	const searchInput = ref('');
	const debouncedSearchInput = refDebounced<string>(searchInput, 500);

	const sorting = ref<SortingState>([
		{
			desc: true,
			id: TABLE_ACCESSOR_KEYS.watchedMs, // accessorKey
		},
	]);
	const columnFilters = ref<ColumnFiltersState>([]);
	const columnVisibility = useLocalStorage<VisibilityState>(COLUMN_VISIBLE_STORAGE_KEY, {});
	const rowSelection = ref({});

	const tableOrder = computed(() => {
		return sorting.value[0].desc
			? CommunityUsersOrder.Desc
			: CommunityUsersOrder.Asc;
	});

	const tableSortBy = computed(() => {
		const sortingItem = sorting.value[0];
		switch (sortingItem.id) {
			case TABLE_ACCESSOR_KEYS.messages:
				return CommunityUsersSortBy.Messages;
			case TABLE_ACCESSOR_KEYS.usedEmotes:
				return CommunityUsersSortBy.UsedEmotes;
			case TABLE_ACCESSOR_KEYS.usedChannelPoints:
				return CommunityUsersSortBy.UsedChannelsPoints;
			case TABLE_ACCESSOR_KEYS.watchedMs:
			default:
				return CommunityUsersSortBy.Watched;
		}
	});

	return {
		searchInput,
		debouncedSearchInput,

		tableOrder,
		tableSortBy,

		sorting,
		columnFilters,
		columnVisibility,
		rowSelection,
	};
});
