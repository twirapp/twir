<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import { useNotificationsTable } from '../composables/use-notifications-table.js';

import { Button } from '@/components/ui/button';

const { t } = useI18n();
const notificationsTable = useNotificationsTable();
</script>

<template>
	<div class="flex items-center gap-2 max-lg:w-full">
		<div class="flex-1 text-sm text-muted-foreground">
			{{ t('sharedTexts.pagination', {
				page: notificationsTable.table.getState().pagination.pageIndex + 1,
				total: notificationsTable.table.getPageCount().toLocaleString(),
				items: notificationsTable.totalNotifications,
			}) }}
		</div>
		<Button
			variant="outline"
			size="sm"
			:disabled="!notificationsTable.table.getCanPreviousPage()"
			@click="notificationsTable.table.previousPage()"
		>
			{{ t('sharedButtons.previous') }}
		</Button>
		<Button
			variant="outline"
			size="sm"
			:disabled="!notificationsTable.table.getCanNextPage()"
			@click="notificationsTable.table.nextPage()"
		>
			{{ t('sharedButtons.next') }}
		</Button>
	</div>
</template>
