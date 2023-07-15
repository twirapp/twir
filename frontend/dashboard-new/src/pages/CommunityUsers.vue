<script setup lang='ts'>
import type { GetUsersResponse_User } from '@twir/grpc/generated/api/api/community';
import {
	type DataTableColumns,
	type PaginationProps,
	NAvatar,
	NDataTable,
	NTag,
	NPagination,
	NSpace,
} from 'naive-ui';
import type { TableColumn } from 'naive-ui/es/data-table/src/interface';
import { ref, computed, h } from 'vue';

import {
	useCommunityUsers,
	UsersOrder,
	UsersSortBy,
	type GetCommunityUsersOpts,
	useTwitchGetUsers,
} from '@/api/index.js';

const communityManager = useCommunityUsers();

const usersOpts = ref<GetCommunityUsersOpts>({
	page: 1,
	limit: 100,
	order: UsersOrder.Desc,
	sortBy: UsersSortBy.Watched,
});
const users = communityManager.getAll(usersOpts);

const usersIdsForRequest = computed(() => {
	return users.data?.value?.users.map((user) => user.id) ?? [];
});
const twitchUsers = useTwitchGetUsers({
	ids: usersIdsForRequest,
});

const HOUR = 1000 * 60 * 60;
const columns = ref<Array<TableColumn>>([
	{
		title: '',
		key: 'avatar',
		width: 50,
		render(row) {
			const twitchUser = twitchUsers.data.value?.users.find((user) => user.id === row.id);

			return h(NAvatar, { src: twitchUser?.profileImageUrl });
		},
	},
	{
		title: 'User',
		key: 'id',
		render(row) {
			const twitchUser = twitchUsers.data.value?.users.find((user) => user.id === row.id);
			return h(NTag, { type: 'info', bordered: false }, { default: () => twitchUser?.displayName ?? row.id });
		},
	},
	{
		title: 'Watched time',
		key: 'watched',
		render(row) {
			return `${(Number(row.watched) / HOUR).toFixed(1)}h`;
		},
		sorter: true,
		sortOrder: 'descend',
	},
	{
		title: 'Messages',
		key: 'messages',
		sorter: true,
		sortOrder: false,
	},
	{
		title: 'Used emotes',
		key: 'emotes',
		sorter: true,
		sortOrder: false,
	},
	{
		title: 'Used channel points',
		key: 'usedChannelPoints',
		sorter: true,
		sortOrder: false,
	},
]);

const paginationOptions = computed<PaginationProps>(() => {
	return {
		page: usersOpts.value.page,
		pageSize: usersOpts.value.limit,
		itemCount: users.data?.value?.totalUsers ?? 0,
		prefix ({ itemCount }) {
			return `Total ${itemCount}`;
		},
	};
});

const Pagination = () => h(NPagination, {
	...paginationOptions.value,
	onUpdatePage: (page: number) => {
		handlePageChange(page);
	},
});

function handlePageChange(page: number) {
	usersOpts.value.page = page;
}

function handleFiltersChange(filters: any) {
	console.log(filters);
}

function handleSorterChange(sorter: { columnKey: string; order: 'ascend' | 'descend' | false }) {
	console.log(sorter);
	const column: any = columns.value.find((column: any) => column.key === sorter.columnKey);
	for (const column of columns.value) {
		(column as any).sortOrder = false;
	}

	column.sortOrder = sorter.order;

	if (sorter.order === 'ascend') {
		usersOpts.value.order = UsersOrder.Asc;
	} else if (sorter.order === 'descend') {
		usersOpts.value.order = UsersOrder.Desc;
	} else {
		usersOpts.value.order = UsersOrder.Desc;
	}

	if (sorter.columnKey === UsersSortBy.Watched) {
		usersOpts.value.sortBy = UsersSortBy.Watched;
	} else if (sorter.columnKey === UsersSortBy.Messages) {
		usersOpts.value.sortBy = UsersSortBy.Messages;
	} else if (sorter.columnKey === UsersSortBy.Emotes) {
		usersOpts.value.sortBy = UsersSortBy.Emotes;
	} else if (sorter.columnKey === UsersSortBy.UsedChannelPoints) {
		usersOpts.value.sortBy = UsersSortBy.UsedChannelPoints;
	} else {
		usersOpts.value.sortBy = UsersSortBy.Watched;
	}
}
</script>

<template>
  <n-space justify="end" style="margin-bottom: 15px;">
    <Pagination />
  </n-space>
  <n-data-table
    :loading="users.isLoading.value || twitchUsers.isLoading.value"
    :columns="columns as any"
    :data="users.data.value?.users ?? []"
    remote
    @update:filters="handleFiltersChange"
    @update:sorter="handleSorterChange"
  />
  <n-space justify="end" style="margin-top: 15px;">
    <Pagination />
  </n-space>
</template>
