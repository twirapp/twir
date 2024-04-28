<script setup lang="ts">
import CommunityUsersPage from './components/community-users-page.vue'
import { useCommunityUsersTable } from './composables/use-community-users-table.js'

import Pagination from '@/components/pagination.vue'

const communityUsersTable = useCommunityUsersTable()
</script>

<template>
	<CommunityUsersPage>
		<template #pagination>
			<Pagination
				:total="communityUsersTable.totalUsers"
				:table="communityUsersTable.table"
				:pagination="communityUsersTable.pagination"
				@update:page="(page) => communityUsersTable.pagination.pageIndex = page"
				@update:page-size="(pageSize) => communityUsersTable.pagination.pageSize = pageSize"
			/>
		</template>
	</CommunityUsersPage>
</template>
