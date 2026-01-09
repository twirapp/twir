<script setup lang="ts">


import NotificationsTableSearch from './notifications-table-search.vue'
import { useNotificationsTable } from '../composables/use-notifications-table'

import Table from '#layers/dashboard/components/table.vue'


const { t } = useI18n()
const notificationsTable = useNotificationsTable()
</script>

<template>
	<NotificationsTableSearch />
	<slot name="pagination" />

	<UiTooltipProvider :delay-duration="100">
		<Table :table="notificationsTable.table" :is-loading="notificationsTable.isLoading.value">
			<template #empty-message>
				{{ t('adminPanel.notifications.emptyNotifications') }}
			</template>
		</Table>
	</UiTooltipProvider>

	<slot name="pagination" />
</template>
