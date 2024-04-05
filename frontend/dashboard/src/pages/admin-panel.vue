<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import AdminManageUsers from '@/features/admin-panel/manage-users/manage-users.vue';
import AdminNotifications from '@/features/admin-panel/notifications/notifications.vue';
import PageLayout, { type PageLayoutTab } from '@/layout/page-layout.vue';

const { t } = useI18n();

const pageTabs: PageLayoutTab[] = [
	{
		name: 'notifications',
		title: t('adminPanel.notifications.title'),
		component: AdminNotifications,
	},
	{
		name: 'user-roles',
		title: t('adminPanel.manageUsers.title'),
		component: AdminManageUsers,
	},
];
</script>

<template>
	<page-layout active-tab="notifications" :tabs="pageTabs">
		<template #title>
			{{ t('adminPanel.title') }}
		</template>
	</page-layout>
</template>
