<script setup lang="ts">
import Table from '~~/layers/dashboard/components/table.vue'

import { TooltipProvider } from '@/components/ui/tooltip'

import { useNotificationsTable } from '../composables/use-notifications-table'
import NotificationsTableSearch from './notifications-table-search.vue'

const { t } = useI18n()
const notificationsTable = useNotificationsTable()
</script>

<template>
	<NotificationsTableSearch />
	<slot name="pagination" />

	<TooltipProvider :delay-duration="100">
		<Table
			:table="notificationsTable.table"
			:is-loading="notificationsTable.isLoading.value"
		>
			<template #empty-message>
				{{ t('adminPanel.notifications.emptyNotifications') }}
			</template>
		</Table>
	</TooltipProvider>

	<slot name="pagination" />
</template>
