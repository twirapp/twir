import { type ColumnDef, getCoreRowModel, useVueTable, type PaginationState } from '@tanstack/vue-table';
import type {
	UsersGetRequest,
	UsersGetResponse_UsersGetResponseUser as User,
} from '@twir/api/messages/admin_users/admin_users';
import { defineStore } from 'pinia';
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUsersTableFilters } from './use-users-table-filters';
import UsersTableActions from '../components/users-table-actions.vue';
import UsersTableAvatar from '../components/users-table-avatar.vue';

import { useAdminUsers } from '@/api/manage-users';

export const useUsersTable = defineStore('manage-users/users-table', () => {
	const { t } = useI18n();

	const pagination = ref({
		pageIndex: 0,
		pageSize: 10,
	});

	function setPagination({
		pageIndex,
		pageSize,
	}: PaginationState): PaginationState {
		pagination.value.pageIndex = pageIndex;
		pagination.value.pageSize = pageSize;
		return { pageIndex, pageSize };
	}

	const tableFilters = useUsersTableFilters();

	const tableParams = computed<UsersGetRequest>(() => ({
		...tableFilters.selectedFilters,
		search: tableFilters.debounceSearchInput,
		page: pagination.value.pageIndex,
		perPage: pagination.value.pageSize,
	}));

	const { data, isFetching } = useAdminUsers(tableParams);

	const users = computed<User[]>(() => {
		if (!data.value) return [];
		return data.value.users;
	});

	const totalUsers = computed(() => data.value?.total ?? 0);

	const pageCount = computed(() => {
		return Math.ceil((data.value?.total ?? 0) / pagination.value.pageSize);
	});

	const tableColumns = computed<ColumnDef<User>[]>(() => [
		{
			accessorKey: 'userLogin',
			size: 60,
			header: () => h('div', {}, t('adminPanel.manageUsers.user')),
			cell: ({ row }) => {
				return h('a',
					{
						class: 'flex items-center gap-4 flex-wrap max-sm:justify-center',
						href: `https://twitch.tv/${row.original.userName}`,
						target: '_blank',
					},
					[
						h(UsersTableAvatar, { avatar: row.original.avatar, name: row.original.userName }),
						row.original.userDisplayName,
					],
				);
			},
		},
		{
			accessorKey: 'userId',
			size: 30,
			header: () => h('div', {}, t('adminPanel.manageUsers.userId')),
			cell: ({ row }) => {
				return h('span', row.original.id);
			},
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => {
				return h(UsersTableActions, {
					userId: row.original.id,
					isBanned: row.original.isBanned,
					isAdmin: row.original.isAdmin,
				});
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
			return users.value;
		},
		get columns() {
			return tableColumns.value;
		},
		manualPagination: true,
		getCoreRowModel: getCoreRowModel(),
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				setPagination(
					updater({
						pageIndex: pagination.value.pageIndex,
						pageSize: pagination.value.pageSize,
					}),
				);
			} else {
				setPagination(updater);
			}
		},
	});

	return {
		isLoading: isFetching,
		totalUsers,
		table,
		tableColumns,
	};
});
