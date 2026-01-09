<script setup lang="ts">
import { useKeywordsTable } from './composables/use-keywords-table'
import KeywordsInformation from './ui/keywords-information.vue'

import Table from '#layers/dashboard/components/table.vue'

const table = useKeywordsTable()
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<KeywordsInformation />
		<Table :table="table.table" :is-loading="table.isLoading.value" />
	</div>
</template>
