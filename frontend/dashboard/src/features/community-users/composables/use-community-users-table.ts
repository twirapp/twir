import { getCoreRowModel, useVueTable, type ColumnDef } from '@tanstack/vue-table';
import { defineStore } from 'pinia';
import { computed, h } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommunityTableActions } from './use-community-table-actions.js';

import { useProfile } from '@/api';
import { useCommunityUsersApi, type CommunityUser } from '@/api/community-users.js';
import { usePagination } from '@/composables/use-pagination.js';
import UsersTableCellUser from '@/features/admin-panel/manage-users/components/users-table-cell-user.vue';
import { type CommunityUsersOpts } from '@/gql/graphql';
import { resolveUserName } from '@/helpers/resolveUserName.js';

const ONE_HOUR = 60 * 60 * 1000;

export const useCommunityUsersTable = defineStore('features/community-users-table', () => {
	const communityUsersApi = useCommunityUsersApi();
	const { data: profile } = useProfile();
	const { t } = useI18n();

	const tableActions = useCommunityTableActions();

	const { pagination, setPagination } = usePagination();
	const params = computed<CommunityUsersOpts>((prevParams) => {
		// reset pagination on search change
		if (prevParams?.query !== tableActions.debouncedSearchInput) {
			pagination.value.pageIndex = 0;
		}

		return {
			channelId: profile.value?.selectedDashboardId ?? '',
			query: tableActions.debouncedSearchInput,
			page: pagination.value.pageIndex,
			perPage: pagination.value.pageSize,
			order: tableActions.tableOrder,
			sortBy: tableActions.tableSortBy,
		};
	});

	const { data, fetching } = communityUsersApi.useCommunityUsers(params);
	const communityUsers = computed<CommunityUser[]>(() => {
		if (!data.value) return [];
		return data.value.communityUsers.users;
	});

	const totalUsers = computed(() => data.value?.communityUsers.total ?? 0);
	const pageCount = computed(() => {
		return Math.ceil(totalUsers.value / pagination.value.pageSize);
	});

	const tableColumns = computed<ColumnDef<CommunityUser>[]>(() => [
		{
			accessorKey: 'id',
			size: 20,
			header: () => h('div', t('community.users.table.user')),
			cell: ({ row }) => {
				return h('a',
					{
						class: 'flex flex-col',
						href: `https://twitch.tv/${row.original.twitchProfile.login}`,
						target: '_blank',
					},
					h(UsersTableCellUser, {
						avatar: row.original.twitchProfile.profileImageUrl,
						userId: row.original.id,
						name: resolveUserName(row.original.twitchProfile.login, row.original.twitchProfile.displayName),
					}),
				);
			},
		},
		{
			accessorKey: 'messages',
			size: 20,
			header: () => h('div', t('community.users.table.messages')),
			cell: ({ row }) => {
				return h('div', row.original.messages);
			},
		},
		{
			accessorKey: 'usedChannelPoints',
			size: 20,
			header: () => h('div', t('community.users.table.usedChannelPoints')),
			cell: ({ row }) => {
				return h('div', row.original.usedChannelPoints);
			},
		},
		{
			accessorKey: 'usedEmotes',
			size: 20,
			header: () => h('div', t('community.users.table.usedEmotes')),
			cell: ({ row }) => {
				return h('div', row.original.usedEmotes);
			},
		},
		{
			accessorKey: 'watchedMs',
			size: 20,
			header: () => h('div', t('community.users.table.watchedTime')),
			cell: ({ row }) => {
				return h('div', `${(Number(row.original.watchedMs) / ONE_HOUR).toFixed(1)}h`);
			},
		},
	]);

	const table = useVueTable({
		get pageCount() {
			return pageCount.value;
		},
		state: {
			pagination: pagination.value,
		},
		get data() {
			return communityUsers.value;
		},
		get columns() {
			return tableColumns.value;
		},
		manualPagination: true,
		getCoreRowModel: getCoreRowModel(),
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				setPagination(updater(pagination.value));
			} else {
				setPagination(updater);
			}
		},
	});

	return {
		isLoading: fetching,
		table,
		totalUsers,
		pagination,
		setPagination,
	};
});
