<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import { useNotificationsTable } from '../composables/use-notifications-table.js';

import SearchBar from '@/components/search-bar.vue';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

const { t } = useI18n();
const notificationsTable = useNotificationsTable();
</script>

<template>
	<div class="flex gap-2 max-lg:w-full">
		<search-bar v-model="notificationsTable.notificationsSearchInput" />
		<Select v-model="notificationsTable.notificationsFilter">
			<SelectTrigger class="h-9 w-[120px]">
				<SelectValue />
			</SelectTrigger>
			<SelectContent>
				<SelectGroup>
					<SelectItem value="globals">
						{{ t('adminPanel.notifications.globals') }}
					</SelectItem>
					<SelectItem value="users">
						{{ t('adminPanel.notifications.users') }}
					</SelectItem>
				</SelectGroup>
			</SelectContent>
		</Select>
	</div>
</template>
