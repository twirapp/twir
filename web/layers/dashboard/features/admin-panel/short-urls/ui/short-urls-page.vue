<script setup lang="ts">
import Table from '~~/layers/dashboard/components/table.vue'
import {
	useAdminShortUrlsApi,
} from '~~/layers/dashboard/features/admin-panel/short-urls/composables/use-admin-short-urls-api.js'
import {
	useAdminShortUrlsTable,
} from '~~/layers/dashboard/features/admin-panel/short-urls/composables/use-admin-short-urls-table.js'
import ShortUrlsForm from '~~/layers/dashboard/features/admin-panel/short-urls/ui/short-urls-form.vue'

const table = useAdminShortUrlsTable()
const api = useAdminShortUrlsApi()
</script>

<template>
	<div class="flex flex-col gap-4">
		<ShortUrlsForm />
		<slot name="pagination" />
		<Table :table="table.table" :is-loading="api.isLoading.value" />
		<slot name="pagination" />
	</div>
</template>
