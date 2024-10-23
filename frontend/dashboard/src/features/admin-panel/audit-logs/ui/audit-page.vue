<script setup lang="ts">
import AuditTable from '@/features/admin-panel/audit-logs/ui/audit-table.vue'
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<slot name="pagination" />
		<AuditTable />
		<slot name="pagination" />
	</div>
</template>
