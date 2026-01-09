<script setup lang="ts">
import Pagination from '#layers/dashboard/components/pagination.vue'
import { useAuditTable } from '~/features/admin-panel/audit-logs/composables/use-audit-table'
import AuditPage from '~/features/admin-panel/audit-logs/ui/audit-page.vue'

const auditTable = useAuditTable()
</script>

<template>
	<AuditPage>
		<template #pagination>
			<Pagination
				:total="auditTable.total.value"
				:table="auditTable.table"
				:pagination="auditTable.pagination.value"
				@update:page="(page) => auditTable.pagination.value.pageIndex = page"
				@update:page-size="(pageSize) => auditTable.pagination.value.pageSize = pageSize"
			/>
		</template>
	</AuditPage>
</template>
