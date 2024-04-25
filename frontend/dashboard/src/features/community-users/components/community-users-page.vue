<script setup lang="ts">

import CommunityUsersTableSearch from './community-users-table-search.vue';
import CommunityUsersTable from './community-users-table.vue';
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<community-users-table-search />
		<slot name="pagination" />
		<community-users-table />
		<slot name="pagination" />
	</div>
</template>
