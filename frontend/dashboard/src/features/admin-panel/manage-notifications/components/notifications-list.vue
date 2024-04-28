<script setup lang="ts">
import NotificationsTable from './notifications-table.vue';
import { useNotificationsTable } from '../composables/use-notifications-table';

import Pagination from '@/components/pagination.vue';

const notifications = useNotificationsTable();
</script>

<template>
	<notifications-table>
		<template #pagination>
			<pagination
				:total="notifications.totalNotifications"
				:table="notifications.table"
				:pagination="notifications.pagination"
				@update:page="(page) => notifications.pagination.pageIndex = page"
				@update:page-size="(pageSize) => notifications.pagination.pageSize = pageSize"
			/>
		</template>
	</notifications-table>
</template>
