<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import AdminManageBadges from '@/features/admin-panel/manage-badges/manage-badges.vue';
import AdminManageNotifications from '@/features/admin-panel/manage-notifications/manage-notifications.vue';
import AdminManageUsers from '@/features/admin-panel/manage-users/manage-users.vue';
import PageLayout, { type PageLayoutTab } from '@/layout/page-layout.vue';

const { t } = useI18n();

const pageTabs: PageLayoutTab[] = [
	{
		name: 'users',
		title: t('adminPanel.manageUsers.title'),
		component: AdminManageUsers,
	},
	{
		name: 'notifications',
		title: t('adminPanel.notifications.title'),
		component: AdminManageNotifications,
	},
	{
		name: 'badges',
		title: t('adminPanel.manageBadges.title'),
		component: AdminManageBadges,
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
