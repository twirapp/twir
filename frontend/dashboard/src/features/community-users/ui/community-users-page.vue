<script setup lang="ts">
import CommunityUsersTableSearch from './community-users-table-search.vue'
import CommunityUsersTable from './community-users-table.vue'
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<CommunityUsersTableSearch />
		<slot name="pagination" />
		<CommunityUsersTable />
		<slot name="pagination" />
	</div>
</template>
