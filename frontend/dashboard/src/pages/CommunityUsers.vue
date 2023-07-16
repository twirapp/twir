<script setup lang='ts'>
import {
	type PaginationProps,
	NAvatar,
	NDataTable,
	NTag,
	NPagination,
	NSpace,
	NButton,
	NPopconfirm,
} from 'naive-ui';
import type { TableColumn } from 'naive-ui/es/data-table/src/interface';
import { ref, computed, h } from 'vue';

import {
	useCommunityUsers,
	ComminityOrder,
	CommunitySortBy,
	GetCommunityUsersOpts,
	useTwitchGetUsers,
	useCommunityReset,
	CommunityResetStatsField,
} from '@/api/index.js';

const communityManager = useCommunityUsers();

const usersOpts = ref<GetCommunityUsersOpts>({
	page: 1,
	limit: 100,
	order: ComminityOrder.Desc,
	sortBy: CommunitySortBy.Watched,
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

			return h(NAvatar, {
				src: twitchUser?.profileImageUrl,
				style: 'display: flex;',
			});
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

function handleSorterChange(sorter: { columnKey: string; order: 'ascend' | 'descend' | false }) {
	const column: any = columns.value.find((column: any) => column.key === sorter.columnKey);
	for (const column of columns.value) {
		(column as any).sortOrder = false;
	}

	column.sortOrder = sorter.order;

	if (sorter.order === 'ascend') {
		usersOpts.value.order = ComminityOrder.Asc;
	} else if (sorter.order === 'descend') {
		usersOpts.value.order = ComminityOrder.Desc;
	} else {
		usersOpts.value.order = ComminityOrder.Desc;
	}

	if (sorter.columnKey === CommunitySortBy.Watched) {
		usersOpts.value.sortBy = CommunitySortBy.Watched;
	} else if (sorter.columnKey === CommunitySortBy.Messages) {
		usersOpts.value.sortBy = CommunitySortBy.Messages;
	} else if (sorter.columnKey === CommunitySortBy.Emotes) {
		usersOpts.value.sortBy = CommunitySortBy.Emotes;
	} else if (sorter.columnKey === CommunitySortBy.UsedChannelPoints) {
		usersOpts.value.sortBy = CommunitySortBy.UsedChannelPoints;
	} else {
		usersOpts.value.sortBy = CommunitySortBy.Watched;
	}
}

const resetter = useCommunityReset();
async function handleReset(field: CommunityResetStatsField) {
	await resetter.mutateAsync(field);
	usersOpts.value.page = 1;
	usersOpts.value.order = ComminityOrder.Desc;
	usersOpts.value.sortBy = CommunitySortBy.Watched;
}
</script>

<template>
  <n-space justify="space-between" style="margin-bottom: 15px;">
    <n-space>
      <n-popconfirm @positive-click="handleReset(CommunityResetStatsField.Watched)">
        <template #trigger>
          <n-button secondary type="warning">
            Reset watched
          </n-button>
        </template>
        Are you sure?
      </n-popconfirm>
      <n-popconfirm @positive-click="handleReset(CommunityResetStatsField.Messages)">
        <template #trigger>
          <n-button secondary type="warning">
            Reset messages
          </n-button>
        </template>
        Are you sure?
      </n-popconfirm>
      <n-popconfirm @positive-click="handleReset(CommunityResetStatsField.Emotes)">
        <template #trigger>
          <n-button secondary type="warning">
            Reset emotes
          </n-button>
        </template>
        Are you sure?
      </n-popconfirm>
      <n-popconfirm @positive-click="handleReset(CommunityResetStatsField.UsedChannelPoints)">
        <template #trigger>
          <n-button secondary type="warning">
            Reset points
          </n-button>
        </template>
        Are you sure?
      </n-popconfirm>
    </n-space>
    <Pagination />
  </n-space>
  <n-data-table
    :loading="users.isLoading.value || twitchUsers.isLoading.value"
    :columns="columns as any"
    :data="users.data.value?.users ?? []"
    remote
    @update:sorter="handleSorterChange"
  />
  <n-space justify="end" style="margin-top: 15px;">
    <Pagination />
  </n-space>
</template>
