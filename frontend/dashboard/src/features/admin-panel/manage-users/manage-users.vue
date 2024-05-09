<script setup lang="ts">
import UsersPage from './components/users-page.vue'
import { useUsersTable } from './composables/use-users-table.js'

import Pagination from '@/components/pagination.vue'

const usersTable = useUsersTable()
</script>

<template>
	<UsersPage>
		<template #pagination>
			<Pagination
				:total="usersTable.totalUsers.value"
				:table="usersTable.table"
				:pagination="usersTable.pagination.value"
				@update:page="(page) => usersTable.pagination.value.pageIndex = page"
				@update:page-size="(pageSize) => usersTable.pagination.value.pageSize = pageSize"
			/>
		</template>
	</UsersPage>
</template>
