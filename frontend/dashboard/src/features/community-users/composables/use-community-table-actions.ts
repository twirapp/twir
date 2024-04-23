import type { ColumnFiltersState, SortingState, VisibilityState } from '@tanstack/vue-table';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { TABLE_ACCESSOR_KEYS } from './use-community-users-table';

import { CommunityUsersOrder, CommunityUsersSortBy } from '@/gql/graphql';

export const useCommunityTableActions = defineStore('features/community-actions', () => {
	const sorting = ref<SortingState>([
		{
			desc: true,
			id: TABLE_ACCESSOR_KEYS.watchedMs, // accessorKey
		},
	]);
	const columnFilters = ref<ColumnFiltersState>([]);
	const columnVisibility = ref<VisibilityState>({});
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
		tableOrder,
		tableSortBy,

		sorting,
		columnFilters,
		columnVisibility,
		rowSelection,
	};
});
