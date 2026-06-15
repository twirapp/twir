<script setup lang="ts">
import Table from '~~/layers/dashboard/app/components/table.vue'
import { useTimersTable } from '~~/layers/dashboard/app/features/timers/composables/use-timers-table'

const table = useTimersTable()
</script>

<template>
	<Table :table="table.table" :is-loading="table.isLoading.value" />
</template>
