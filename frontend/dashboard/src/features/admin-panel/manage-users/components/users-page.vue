<script setup lang="ts">
import UsersTableSearch from './users-table-search.vue';
import UsersTable from './users-table.vue';

</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<users-table-search />
		<slot name="pagination" />
		<users-table />
		<slot name="pagination" />
	</div>
</template>
