<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import { useUsersTable } from '../composables/use-users-table.js';

import { Button } from '@/components/ui/button';

const usersTable = useUsersTable();
const { t } = useI18n();
</script>

<template>
	<div class="flex items-center justify-end space-x-2 py-4">
		<div class="flex-1 text-sm text-muted-foreground">
			{{
				t('sharedTexts.pagination', {
					page: usersTable.table.getState().pagination.pageIndex + 1,
					total: usersTable.table.getPageCount(),
					items: usersTable.totalUsers,
				})
			}}
		</div>
		<div class="space-x-2">
			<Button
				variant="outline"
				size="sm"
				:disabled="!usersTable.table.getCanPreviousPage()"
				@click="usersTable.table.previousPage()"
			>
				{{ t('sharedButtons.previous') }}
			</Button>
			<Button
				variant="outline"
				size="sm"
				:disabled="!usersTable.table.getCanNextPage()"
				@click="usersTable.table.nextPage()"
			>
				{{ t('sharedButtons.next') }}
			</Button>
		</div>
	</div>
</template>
