<script setup lang="ts">
import { FlexRender } from '@tanstack/vue-table'

import { useCommunityUsersTable } from '~/features/community-users/composables/use-community-users-table'
import CommunityUsersTableSearch
	from '~/features/community-users/ui/community-users-table-search.vue'

definePageMeta({
	layout: 'public',
})

const store = useCommunityUsersTable()
await useAsyncData('communityUsers', async () => store.fetchUsers().then(() => true))
</script>

<template>
	<div class="flex flex-col gap-4">
		<CommunityUsersTableSearch />
		<Pagination
			:total="store.totalUsers"
			:table="store.table"
			:pagination="store.pagination"
			@update:page="(page) => store.pagination.pageIndex = page"
			@update:page-size="(pageSize) => store.pagination.pageSize = pageSize"
		/>

		<div class="flex-wrap w-full border rounded-md" style="background-color: rgb(24, 24, 28)">
			<UiTable>
				<UiTableHeader>
					<UiTableRow
						v-for="headerGroup in store.table.getHeaderGroups()"
						:key="headerGroup.id"
						class="text-slate-50"
					>
						<UiTableHead
							v-for="header in headerGroup.headers"
							:key="header.id"
							:style="{ width: `${header.column.columnDef.size}%` }"
						>
							<FlexRender
								v-if="!header.isPlaceholder"
								:render="header.column.columnDef.header"
								:props="header.getContext()"
							/>
						</UiTableHead>
					</UiTableRow>
				</UiTableHeader>

				<UiTableBody>
					<UiTableRow
						v-for="row in store.table.getRowModel().rows"
						:key="row.id"
					>
						<template
							v-for="cell in row.getVisibleCells()"
							:key="cell.id"
						>
							<UiTableCell>
								<FlexRender
									:render="cell.column.columnDef.cell"
									:props="cell.getContext()"
								/>
							</UiTableCell>
						</template>
					</UiTableRow>
				</UiTableBody>
			</UiTable>
		</div>

		<Pagination
			:total="store.totalUsers"
			:table="store.table"
			:pagination="store.pagination"
			@update:page="(page) => store.pagination.pageIndex = page"
			@update:page-size="(pageSize) => store.pagination.pageSize = pageSize"
		/>
	</div>
</template>
