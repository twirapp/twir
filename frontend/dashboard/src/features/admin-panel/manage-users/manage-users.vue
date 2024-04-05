<script setup lang="ts">

import UsersTablePagination from './components/users-table-pagination.vue';
import UsersTableSearch from './components/users-table-search.vue';
import UsersTable from './components/users-table.vue';

import ShadcnLayout from '@/layout/shadcn-layout.vue';
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<users-table-search />
		<users-table-pagination />
		<shadcn-layout>
			<users-table />
		</shadcn-layout>
		<users-table-pagination />
	</div>
</template>
