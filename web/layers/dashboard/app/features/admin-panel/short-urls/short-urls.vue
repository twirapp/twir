<script setup lang="ts">
import Pagination from '~~/layers/dashboard/app/components/pagination.vue'
import {
	useAdminShortUrlsApi,
} from '~~/layers/dashboard/app/features/admin-panel/short-urls/composables/use-admin-short-urls-api.js'
import {
	useAdminShortUrlsTable,
} from '~~/layers/dashboard/app/features/admin-panel/short-urls/composables/use-admin-short-urls-table.js'
import ShortUrlsPage from '~~/layers/dashboard/app/features/admin-panel/short-urls/ui/short-urls-page.vue'

const table = useAdminShortUrlsTable()
const api = useAdminShortUrlsApi()
</script>

<template>
	<ShortUrlsPage>
		<template #pagination>
			<Pagination
				:total="api.totalItems.value"
				:table="table.table"
				:pagination="api.pagination.value"
				@update:page="(page) => api.pagination.value.pageIndex = page"
				@update:page-size="(pageSize) => api.pagination.value.pageSize = pageSize"
			/>
		</template>
	</ShortUrlsPage>
</template>
