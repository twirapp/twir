<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import AdminNotifications from './admin-notifications.vue';

import PageLayout, { type PageLayoutTab } from '@/layout/page-layout.vue';

const { t } = useI18n();

const pageTabs: PageLayoutTab[] = [
	{
		name: 'notifications',
		title: t('adminPanel.notifications.title'),
		component: AdminNotifications,
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
