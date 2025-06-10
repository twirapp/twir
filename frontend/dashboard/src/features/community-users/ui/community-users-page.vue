<script setup lang="ts">
import CommunityUsersTableSearch from './community-users-table-search.vue'
import CommunityUsersTable from './community-users-table.vue'

import CommunityUsersResetButton
	from '@/features/community-users/ui/community-users-reset-button.vue'
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<div class="flex flex-wrap gap-4 w-full">
			<CommunityUsersTableSearch class="flex-grow" />
			<CommunityUsersResetButton />
		</div>
		<slot name="pagination" />
		<CommunityUsersTable />
		<slot name="pagination" />
	</div>
</template>
