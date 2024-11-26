<script setup lang="ts">
import Table from '@/components/table.vue'
import { useTimersTable } from '@/features/timers/composables/use-timers-table'

const table = useTimersTable()
</script>

<template>
	<Table :table="table.table" :is-loading="table.isLoading.value" />
</template>
