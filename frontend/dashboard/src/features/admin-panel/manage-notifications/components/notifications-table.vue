<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import NotificationsTableSearch from './notifications-table-search.vue'
import { useNotificationsTable } from '../composables/use-notifications-table'

import Table from '@/components/table.vue'
import { TooltipProvider } from '@/components/ui/tooltip'

const { t } = useI18n()
const notificationsTable = useNotificationsTable()
</script>

<template>
	<NotificationsTableSearch />
	<slot name="pagination" />

	<TooltipProvider :delay-duration="100">
		<Table :table="notificationsTable.table" :is-loading="notificationsTable.isLoading.value">
			<template #empty-message>
				{{ t('adminPanel.notifications.emptyNotifications') }}
			</template>
		</Table>
	</TooltipProvider>

	<slot name="pagination" />
</template>
