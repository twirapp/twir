<script setup lang="ts">
import { useQuotesTable } from './composables/use-quotes-table.js'

import Pagination from '~~/layers/dashboard/components/pagination.vue'
import SearchBar from '~~/layers/dashboard/components/search-bar.vue'
import Table from '~~/layers/dashboard/components/table.vue'

const quotesTable = useQuotesTable()
</script>

<template>
	<div class="flex w-full flex-col gap-4">
		<SearchBar v-model="quotesTable.searchInput.value" />
		<Pagination
			:total="quotesTable.filteredQuotes.value.length"
			:table="quotesTable.table"
			:pagination="quotesTable.pagination.value"
			@update:page="(page) => quotesTable.pagination.value.pageIndex = page"
			@update:page-size="(pageSize) => quotesTable.pagination.value.pageSize = pageSize"
		/>
		<Table
			:table="quotesTable.table"
			:is-loading="quotesTable.isLoading.value"
		/>
		<Pagination
			:total="quotesTable.filteredQuotes.value.length"
			:table="quotesTable.table"
			:pagination="quotesTable.pagination.value"
			@update:page="(page) => quotesTable.pagination.value.pageIndex = page"
			@update:page-size="(pageSize) => quotesTable.pagination.value.pageSize = pageSize"
		/>
	</div>
</template>
