<script setup lang="ts">
import Table from '@/components/table.vue'
import {
	useGiveawaysListTable,
} from '@/features/giveaways/composables/giveaways-use-list-table.ts'

const { table, isLoading } = useGiveawaysListTable()
</script>

<template>
	<Table :table :is-loading />
</template>
