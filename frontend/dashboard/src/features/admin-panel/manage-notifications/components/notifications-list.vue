<script setup lang="ts">
import NotificationsTable from './notifications-table.vue'
import { useNotificationsTable } from '../composables/use-notifications-table'

import Pagination from '@/components/pagination.vue'

const notifications = useNotificationsTable()
</script>

<template>
	<NotificationsTable>
		<template #pagination>
			<Pagination
				:total="notifications.totalNotifications.value"
				:table="notifications.table"
				:pagination="notifications.pagination.value"
				@update:page="(page) => notifications.pagination.value.pageIndex = page"
				@update:page-size="(pageSize) => notifications.pagination.value.pageSize = pageSize"
			/>
		</template>
	</NotificationsTable>
</template>
