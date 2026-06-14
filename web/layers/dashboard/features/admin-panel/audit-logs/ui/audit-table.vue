<script setup lang="ts">
import Table from '~~/layers/dashboard/components/table.vue'
import { useAuditTable } from '~~/layers/dashboard/features/admin-panel/audit-logs/composables/use-audit-table'

const auditTable = useAuditTable()
</script>

<template>
	<Table
		:table="auditTable.table"
		:is-loading="auditTable.isLoading.value"
	/>
</template>
