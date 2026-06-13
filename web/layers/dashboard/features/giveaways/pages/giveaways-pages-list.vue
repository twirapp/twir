<script setup lang="ts">
import Table from '~~/layers/dashboard/components/table.vue'
import {
	useGiveawaysListTable,
} from '~~/layers/dashboard/features/giveaways/composables/giveaways-use-list-table.js'

const { table, isLoading } = useGiveawaysListTable()
</script>

<template>
	<Table :table :is-loading />
</template>
