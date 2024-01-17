<script setup lang="ts">
import { useVueTable, FlexRender, createColumnHelper, getCoreRowModel, type SortingState } from '@tanstack/vue-table';
import { type GetUsersResponse_User } from '@twir/api/messages/community/community';
import { computed, h, ref } from 'vue';

import { useCommunityUsers, type SortKey } from '@/api/community.js';
import { useTwitchGetUsers } from '@/api/users.js';

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

const { data } = useCommunityUsers(usersOpts);

const usersIdsForRequest = computed(() => {
	return data.value?.users.map((user) => user.id) ?? [];
});
const { data: twitchUsers } = useTwitchGetUsers(usersIdsForRequest);

const columnHelper = createColumnHelper<GetUsersResponse_User>();
const HOUR = 1000 * 60 * 60;

const table = useVueTable({
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
				const user = twitchUsers.value?.users.find(u => u.id === ctx.row.original.id );
				return h('img', { src: user?.profileImageUrl, class: 'rounded-full w-8 h-8' });
			},
			enableSorting: false,
		}),

		columnHelper.accessor('id', {
			header: 'User',
			cell: (ctx) => {
				const user = twitchUsers.value?.users.find(u => u.id === ctx.getValue() );
				return h('span', {}, {
					default: () => user?.displayName.toLocaleLowerCase() === user?.login
						? user?.displayName
						: user?.login,
				});
			},
			enableSorting: false,
		}),

		columnHelper.accessor('watched', {
			header: 'Watched time',
			cell: (ctx) => {
				return h('span', {}, {
					default: () => `${(Number(ctx.row.original.watched) / HOUR).toFixed(1)}h`,
				});
			},
		}),

		columnHelper.accessor('messages', {
			header: 'Messages',
		}),

		columnHelper.accessor('emotes', {
			header: 'Used emotes',
		}),

		columnHelper.accessor('usedChannelPoints', {
			header: 'Used channel points',
		}),
	],
});

const pagesCount = computed(() => Math.ceil((data.value?.totalUsers ?? 0) / 100));
</script>

<template>
	<div>
		<div class="flex justify-between items-center gap-2 pb-2">
			<div class="text-slate-200">
				Total	{{ data?.totalUsers }} users
			</div>
			<div class="flex items-center gap-2">
				<div class="text-slate-200">
					Page {{ page }} of {{ pagesCount }}
				</div>
				<button
					:disabled="page === 1"
					class="p-1 rounded shadow-lg"
					:class="[page === 1 ? 'bg-neutral-700 cursor-not-allowed' : 'bg-neutral-600']"
					@click="page--"
				>
					<svg xmlns="http://www.w3.org/2000/svg" class="text-slate-200" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
						<path stroke="none" d="M0 0h24v24H0z" fill="none" />
						<path d="M14 6l-6 6l6 6v-12" />
					</svg>
				</button>
				<button
					:disabled="page === pagesCount"
					class="p-1 rounded shadow-lg"
					:class="[page === pagesCount ? 'bg-neutral-700 cursor-not-allowed' : 'bg-neutral-600']"
					@click="page++"
				>
					<svg xmlns="http://www.w3.org/2000/svg" class="text-slate-200" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
						<path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
						<path d="M10 18l6 -6l-6 -6v12"></path>
					</svg>
				</button>
			</div>
		</div>
		<div class="overflow-auto overflow-y-hidden rounded-lg border-gray-200 shadow-lg">
			<table class="w-full border-collapse text-left text-sm text-slate-200 relative">
				<thead class="bg-neutral-700 text-slate-200">
					<tr
						v-for="headerGroup in table.getHeaderGroups()"
						:key="headerGroup.id"
					>
						<th
							v-for="header in headerGroup.headers"
							:key="header.id"
							:colSpan="header.colSpan"
							scope="col" class="px-6 py-4 font-medium"
							:class="header.column.getCanSort() ? 'cursor-pointer select-none' : ''"
							@click="header.column.getToggleSortingHandler()?.($event)"
						>
							<template v-if="!header.isPlaceholder">
								<FlexRender
									:render="header.column.columnDef.header"
									:props="header.getContext()"
								/>

								{{
									{ asc: ' 🔼', desc: ' 🔽' }[
										header.column.getIsSorted() as string
									]
								}}
							</template>
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-neutral-600 border-t border-neutral-600 bg-neutral-700">
					<tr v-for="row in table.getRowModel().rows" :key="row.id" class="hover:bg-neutral-600">
						<td v-for="cell in row.getVisibleCells()" :key="cell.id" class="px-6 py-4">
							<FlexRender
								:render="cell.column.columnDef.cell"
								:props="cell.getContext()"
							/>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</template>
