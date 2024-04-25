<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import { useUsersTable } from '../composables/use-users-table.js';

import Table from '@/components/table.vue';

const { t } = useI18n();
const usersTable = useUsersTable();
</script>

<template>
	<Table
		:table="usersTable.table"
		:is-loading="usersTable.isLoading"
	>
		<template #empty-message>
			{{ t('adminPanel.manageUsers.noUsers') }}
		</template>
	</Table>
</template>
