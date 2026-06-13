<script setup lang="ts">
import Table from '@/components/table.vue'
import { useGiveawaysHistoryTable } from '@/features/giveaways/composables/giveaways-use-history-table.js'

const { table, isLoading } = useGiveawaysHistoryTable()
</script>

<template>
	<Table :table="table" :is-loading="isLoading" />
</template>
