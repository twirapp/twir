<script setup lang="ts">
import Table from '@/components/table.vue'
import { useAuditTable } from '@/features/admin-panel/audit-logs/composables/use-audit-table'

const auditTable = useAuditTable()
</script>

<template>
	<Table
		:table="auditTable.table"
		:is-loading="auditTable.isLoading.value"
	/>
</template>
