<script setup lang="ts">
import AuditFilters from './audit-filters.vue'

import AuditTable from '~/features/admin-panel/audit-logs/ui/audit-table.vue'
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<AuditFilters />
		<slot name="pagination" />
		<AuditTable />
		<slot name="pagination" />
	</div>
</template>
