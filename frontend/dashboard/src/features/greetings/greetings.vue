<script setup lang="ts">
import { useGreetingsTable } from './composables/use-greetings-table'

import Table from '@/components/table.vue'

const greetingsTable = useGreetingsTable()
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<Table :table="greetingsTable.table" :is-loading="greetingsTable.isLoading" />
	</div>
</template>
