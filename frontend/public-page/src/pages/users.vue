<script setup lang="ts">
import {
	useVueTable,
	FlexRender,
	createColumnHelper,
	getCoreRowModel,
	type SortingState,
} from '@tanstack/vue-table';
import { type GetUsersResponse_User } from '@twir/api/messages/community/community';
import { computed, h, ref } from 'vue';

import TableRowsSkeleton from '@/components/TableRowsSkeleton.vue';
import { Button } from '@/components/ui/button';
import {
	Table,
	TableBody,
	TableHead,
	TableHeader,
	TableRow,
	TableCell,
} from '@/components/ui/table';
import { useCommunityUsers, type SortKey } from '@/composables/use-community';
import { useTwitchGetUsers } from '@/composables/use-twitch-users';

const props = defineProps<{
	channelId: string
	channelName: string
}>();

const tableSorting = ref<SortingState>([]);

const computedOrderAndSorting = computed<{ sortBy: SortKey, desc: boolean }>(() => {
	const value = tableSorting.value.at(0);
	if (!value) return {
		sortBy: 'watched',
		desc: true,
	};

	return {
		sortBy: value.id as SortKey,
		desc: value.desc,
	};
});
const page = ref(1);

const usersOpts = computed(() => ({
	limit: 100,
	channelId: props.channelId,
	desc: computedOrderAndSorting.value.desc,
	page: page.value,
	sortBy: computedOrderAndSorting.value.sortBy,
}));

const { data, isLoading } = useCommunityUsers(usersOpts);

const usersIdsForRequest = computed(() => {
	return data.value?.users.map((user) => user.id) ?? [];
});
const { data: twitchUsers } = useTwitchGetUsers(usersIdsForRequest);

const columnHelper = createColumnHelper<GetUsersResponse_User>();
const HOUR = 1000 * 60 * 60;

const vueTable = useVueTable({
	get data() {
		return data.value?.users ?? [];
	},

	state: {
		get sorting() {
			return tableSorting.value;
		},
	},

	manualPagination: true,
	manualSorting: true,
	onSortingChange: updaterOrValue => {
		tableSorting.value =
			typeof updaterOrValue === 'function'
				? updaterOrValue(tableSorting.value)
				: updaterOrValue;
	},

	getCoreRowModel: getCoreRowModel(),

	columns: [
		columnHelper.display({
			header: '',
			id: 'avatar',
			cell: (ctx) => {
				const user = twitchUsers.value?.users.find(u => u.id === ctx.row.original.id);
				return h('img', { src: user?.profileImageUrl, class: 'rounded-full w-8 h-8' });
			},
			size: 5,
			enableSorting: false,
		}),

		columnHelper.accessor('id', {
			header: 'User',
			cell: (ctx) => {
				const user = twitchUsers.value?.users.find(u => u.id === ctx.getValue());
				return h('span', {}, {
					default: () => user?.displayName.toLocaleLowerCase() === user?.login
						? user?.displayName
						: user?.login,
				});
			},
			enableSorting: false,
			size: 30,
		}),

		columnHelper.accessor('watched', {
			header: 'Watched time',
			cell: (ctx) => {
				return h('span', {}, {
					default: () => `${(Number(ctx.row.original.watched) / HOUR).toFixed(1)}h`,
				});
			},
			size: 5,
		}),

		columnHelper.accessor('messages', {
			header: 'Messages',
			size: 20,
		}),

		columnHelper.accessor('emotes', {
			header: 'Used emotes',
			size: 20,
		}),

		columnHelper.accessor('usedChannelPoints', {
			header: 'Used channel points',
			size: 20,
		}),
	],
});

const pagesCount = computed(() => Math.ceil((data.value?.totalUsers ?? 0) / 100));
</script>

<template>
	<div>
		<div class="flex justify-between items-center gap-2 pb-2">
			<div class="text-slate-200">
				Total {{ data?.totalUsers ?? 0 }} users
			</div>
			<div class="flex items-center gap-2">
				<div class="text-slate-200">
					Page {{ page }} of {{ pagesCount }}
				</div>
				<Button
					variant="outline"
					:disabled="page === 1"
					size="icon"
					@click="page--"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg" class="text-slate-200" width="24" height="24"
						viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
						stroke-linecap="round" stroke-linejoin="round"
					>
						<path stroke="none" d="M0 0h24v24H0z" fill="none" />
						<path d="M14 6l-6 6l6 6v-12" />
					</svg>
				</Button>
				<Button
					variant="outline"
					:disabled="page === pagesCount"
					size="icon"
					@click="page++"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg" class="text-slate-200" width="24" height="24"
						viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
						stroke-linecap="round" stroke-linejoin="round"
					>
						<path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
						<path d="M10 18l6 -6l-6 -6v12"></path>
					</svg>
				</Button>
			</div>
		</div>

		<div class="rounded-md border">
			<Table>
				<TableHeader>
					<TableRow
						v-for="headerGroup in vueTable.getHeaderGroups()"
						:key="headerGroup.id"
						class="text-slate-50"
					>
						<TableHead
							v-for="header in headerGroup.headers"
							:key="header.id"
							:style="{ width: `${header.getSize()}px` }"
						>
							<FlexRender
								v-if="!header.isPlaceholder"
								:render="header.column.columnDef.header"
								:props="header.getContext()"
							/>
						</TableHead>
					</TableRow>
				</TableHeader>
				<Transition name="table-rows" appear mode="out-in">
					<TableBody v-if="isLoading">
						<table-rows-skeleton :rows="20" :colspan="5" />
					</TableBody>
					<TableBody v-else>
						<TableRow v-for="row in vueTable.getRowModel().rows" :key="row.id">
							<TableCell
								v-for="cell in row.getVisibleCells()"
								:key="cell.id"
							>
								<FlexRender
									:render="cell.column.columnDef.cell"
									:props="cell.getContext()"
								/>
							</TableCell>
						</TableRow>
					</TableBody>
				</Transition>
			</Table>
		</div>
	</div>
</template>

<style scoped>
.table-rows-enter-active,
.table-rows-leave-active {
	transition: opacity 0.5s ease;
}

.table-rows-enter-from,
.table-rows-leave-to {
	opacity: 0;
}
</style>
