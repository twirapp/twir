import {
	type ColumnDef,
	getCoreRowModel,
	useVueTable,
	getPaginationRowModel,
} from '@tanstack/vue-table';
import type { UsersGetRequest, UsersGetResponse_UsersGetResponseUser as User } from '@twir/api/messages/admin_users/admin_users';
import { defineStore } from 'pinia';
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import UsersTableActions from '../components/users-table-actions.vue';

import { useAdminUsers } from '@/api/manage-users';

export const useUsersTable = defineStore('manage-users/users-table', () => {
	const { t } = useI18n();

	const searchInput = ref('');
	const searchFilters = [
		{ value: 'isAdmin', label: 'Admin' },
		{ value: 'isBotEnabled', label: 'Bot enabled' },
	];
	const selectedFilters = ref<Record<string, true | undefined>>({
		isAdmin: undefined,
		isBotEnabled: undefined,
	});
	const pagination = ref({
		pageIndex: 0,
		pageSize: 10,
	});
	const tableParams = computed<UsersGetRequest>(() => ({
		...selectedFilters.value,
		search: searchInput.value,
		page: pagination.value.pageIndex,
		perPage: pagination.value.pageSize,
	}));

	const { data } = useAdminUsers(tableParams);

	const users = computed(() => {
		if (!data.value) return [];
		return data.value.users;
	});

	const pageCount = computed(() => {
		if (!data.value) return 0;
		return data.value.total;
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
						h('img', { class: 'h-9 w-9 rounded-full', src: row.original.avatar, loading: 'lazy' }),
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
		pageCount: pageCount.value,
		state: {
			pagination: pagination.value,
		},
		get data() {
			return users.value;
		},
		get columns() {
			return tableColumns.value;
		},
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
	});

	return {
		table,
		tableColumns,
		searchInput,
		searchFilters,
		selectedFilters,
	};
});
