<script setup lang="ts">
import { useAdminToxicMessagesApi } from './composables/use-admin-toxic-mesasges-api.ts'
import { useAdminToxicMessagesTable } from './composables/use-admin-toxic-messages-table.ts'

import Pagination from '@/components/pagination.vue'
import Table from '@/components/table.vue'
import ToxicMessagesPage from '@/features/admin-panel/toxic-messages/ui/toxic-messages-page.vue'

const api = useAdminToxicMessagesApi()
const table = useAdminToxicMessagesTable()
</script>

<template>
	<ToxicMessagesPage>
		<template #pagination>
			<Pagination
				:total="api.totalItems.value"
				:table="table.table"
				:pagination="api.pagination.value"
				@update:page="(page) => api.pagination.value.pageIndex = page"
				@update:page-size="(pageSize) => api.pagination.value.pageSize = pageSize"
			/>
		</template>

		<template #table>
			<Table :table="table.table" :is-loading="api.isLoading.value" />
		</template>
	</ToxicMessagesPage>
</template>
