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
import type { TableBaseColumn } from 'naive-ui/es/data-table/src/interface';
import { ref, computed, h, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
	useCommunityUsers,
	ComminityOrder,
	CommunitySortBy,
	GetCommunityUsersOpts,
	useTwitchGetUsers,
	useCommunityReset,
	CommunityResetStatsField,
	useProfile,
} from '@/api/index.js';

const { t } = useI18n();
const communityManager = useCommunityUsers();

const { data: profile } = useProfile();

const usersOpts = ref<GetCommunityUsersOpts>({
	page: 1,
	limit: 100,
	order: ComminityOrder.Desc,
	sortBy: CommunitySortBy.Watched,
	channelId: undefined,
});
const users = communityManager.getAll(usersOpts);

watch(() => profile.value?.selectedDashboardId, async (v) => {
	usersOpts.value.channelId = v;
	await users.refetch();
}, { immediate: true });

const usersIdsForRequest = computed(() => {
	return users.data?.value?.users.map((user) => user.id) ?? [];
});
const twitchUsers = useTwitchGetUsers({
	ids: usersIdsForRequest,
});

const HOUR = 1000 * 60 * 60;
const columns = ref<Array<TableBaseColumn & { resetableKey?: CommunityResetStatsField }>>([
	{
		title: '',
		key: 'avatar',
		width: 50,
		render(row) {
			const twitchUser = twitchUsers.data.value?.users.find((user) => user.id === row.id);

			return h(NAvatar, {
				src: twitchUser?.profileImageUrl,
				class: 'flex',
			});
		},
	},
	{
		title: t('community.users.table.user'),
		key: 'id',
		render(row) {
			const twitchUser = twitchUsers.data.value?.users.find((user) => user.id === row.id);
			return h(NTag, { type: 'info', bordered: false }, { default: () => twitchUser?.displayName ?? row.id });
		},
	},
	{
		title: t('community.users.table.watchedTime'),
		key: 'watched',
		render(row) {
			return `${(Number(row.watched) / HOUR).toFixed(1)}h`;
		},
		sorter: true,
		sortOrder: 'descend',
		resetableKey: CommunityResetStatsField.Watched,
	},
	{
		title: t('community.users.table.messages'),
		key: 'messages',
		sorter: true,
		sortOrder: false,
		resetableKey: CommunityResetStatsField.Messages,
	},
	{
		title: t('community.users.table.usedEmotes'),
		key: 'emotes',
		sorter: true,
		sortOrder: false,
		resetableKey: CommunityResetStatsField.Emotes,
	},
	{
		title: t('community.users.table.usedChannelPoints'),
		key: 'usedChannelPoints',
		sorter: true,
		sortOrder: false,
		resetableKey: CommunityResetStatsField.UsedChannelPoints,
	},
]);

const paginationOptions = computed<PaginationProps>(() => {
	return {
		page: usersOpts.value.page,
		pageSize: usersOpts.value.limit,
		itemCount: users.data?.value?.totalUsers ?? 0,
		prefix ({ itemCount }) {
			return t('community.users.total', { total: itemCount });
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
	<n-space justify="space-between" class="mb-4">
		<n-space>
			<n-popconfirm
				v-for="item of columns.filter(c => c.resetableKey !== undefined)"
				:key="item.resetableKey"
				:negative-text="t('sharedButtons.cancel')"
				:positive-text="t('sharedButtons.confirm')"
				@positive-click="handleReset(item.resetableKey!)"
			>
				<template #trigger>
					<n-button secondary type="warning">
						{{ t('community.users.reset.button') }} {{ (item.title! as string).toLowerCase() }}
					</n-button>
				</template>
				{{ t('community.users.reset.resetQuestion', { title: (item.title! as string).toLowerCase() }) }}
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
	<n-space justify="end" class="mt-4">
		<Pagination />
	</n-space>
</template>
