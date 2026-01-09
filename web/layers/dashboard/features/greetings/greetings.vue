<script setup lang="ts">
import { useGreetingsTable } from './composables/use-greetings-table'
import GreetingsInformation from './ui/greetings-information.vue'

import Table from '#layers/dashboard/components/table.vue'

const greetingsTable = useGreetingsTable()
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<GreetingsInformation />
		<Table :table="greetingsTable.table" :is-loading="greetingsTable.isLoading.value" />
	</div>
</template>
