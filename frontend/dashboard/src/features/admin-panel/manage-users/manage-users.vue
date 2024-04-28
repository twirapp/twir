<script setup lang="ts">
import UsersPage from './components/users-page.vue';
import { useUsersTable } from './composables/use-users-table';

import Pagination from '@/components/pagination.vue';

const usersTable = useUsersTable();
</script>

<template>
	<users-page>
		<template #pagination>
			<pagination
				:total="usersTable.totalUsers"
				:table="usersTable.table"
				:pagination="usersTable.pagination"
				@update:page="(page) => usersTable.pagination.pageIndex = page"
				@update:page-size="(pageSize) => usersTable.pagination.pageSize = pageSize"
			/>
		</template>
	</users-page>
</template>
