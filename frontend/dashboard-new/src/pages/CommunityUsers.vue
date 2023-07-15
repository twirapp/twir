<script setup lang='ts'>
import type { GetUsersResponse_User } from '@twir/grpc/generated/api/api/community';
import {
	type DataTableColumns,
	NAvatar,
	NDataTable,
	NTag,
} from 'naive-ui';
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
const columns: DataTableColumns<GetUsersResponse_User> = [
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
	},
	{
		title: 'Messages',
		key: 'messages',
	},
	{
		title: 'Used emotes',
		key: 'emotes',
	},
	{
		title: 'Used channel points',
		key: 'usedChannelPoints',
	},
];
</script>

<template>
  <n-data-table
    :loading="users.isLoading.value || twitchUsers.isLoading.value"
    :columns="columns"
    :data="users.data.value?.users ?? []"
  />
</template>
