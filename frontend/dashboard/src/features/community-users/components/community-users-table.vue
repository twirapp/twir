<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import { useCommunityUsersTable } from '../composables/use-community-users-table';

import Table from '@/components/table.vue';

const { t } = useI18n();
const communityUsersTable = useCommunityUsersTable();
</script>

<template>
	<Table :table="communityUsersTable.table" :is-loading="communityUsersTable.isLoading">
		<template #empty-message>
			{{ t('community.users.table.empty') }}
		</template>
	</Table>
</template>
