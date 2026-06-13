<script setup lang="ts">
import { useCommunityUsersTable } from './composables/use-community-users-table.js'
import CommunityUsersPage from './ui/community-users-page.vue'

import Pagination from '@/components/pagination.vue'

const communityUsersTable = useCommunityUsersTable()
</script>

<template>
	<CommunityUsersPage>
		<template #pagination>
			<Pagination
				:total="communityUsersTable.totalUsers.value"
				:table="communityUsersTable.table"
				:pagination="communityUsersTable.pagination.value"
				@update:page="(page) => communityUsersTable.pagination.value.pageIndex = page"
				@update:page-size="(pageSize) => communityUsersTable.pagination.value.pageSize = pageSize"
			/>
		</template>
	</CommunityUsersPage>
</template>
