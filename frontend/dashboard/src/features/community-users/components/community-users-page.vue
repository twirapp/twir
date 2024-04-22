<script setup lang="ts">
import CommunityTableActions from './community-table-actions.vue';
import CommunityUsersTable from './community-users-table.vue';
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<community-table-actions />
		<slot name="pagination" />
		<community-users-table />
		<slot name="pagination" />
	</div>
</template>
