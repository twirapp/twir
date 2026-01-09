<script setup lang="ts">
import { useUsersTable } from './composables/use-users-table.js'
import UsersPage from './ui/users-page.vue'

import Pagination from '#layers/dashboard/components/pagination.vue'

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
