<script setup lang="ts">
import { FlexRender } from '@tanstack/vue-table'

import { useCommunityUsersTable } from '~/features/community-users/composables/use-community-users-table'

definePageMeta({
	layout: 'public',
})

const store = useCommunityUsersTable()
await useAsyncData('communityUsers', async () => store.fetchUsers().then(() => true))
</script>

<template>
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
</template>
