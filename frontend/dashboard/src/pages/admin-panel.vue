<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import AdminActions from '@/features/admin-panel/actions/admin-actions.vue'
import AdminAuditLogs from '@/features/admin-panel/audit-logs/audit-logs.vue'
import ChatMessages from '@/features/admin-panel/chat-mesages/chat-messages.vue'
import AdminManageBadges from '@/features/admin-panel/manage-badges/manage-badges.vue'
import AdminManageNotifications
	from '@/features/admin-panel/manage-notifications/manage-notifications.vue'
import AdminManageUsers from '@/features/admin-panel/manage-users/manage-users.vue'
import ShortUrls from '@/features/admin-panel/short-urls/short-urls.vue'
import ToxicMessages from '@/features/admin-panel/toxic-messages/toxic-messages.vue'
import PageLayout, { type PageLayoutTab } from '@/layout/page-layout.vue'

const { t } = useI18n()

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
	{
		name: 'actions',
		title: t('adminPanel.adminActions.title'),
		component: AdminActions,
	},
	{
		name: 'audit-logs',
		title: t('adminPanel.auditLogs.title'),
		component: AdminAuditLogs,
	},
	{
		name: 'chat-messages',
		title: 'Chat logs',
		component: ChatMessages,
	},
	{
		name: 'short-urls',
		title: 'Short Urls',
		component: ShortUrls,
	},
	{
		name: 'toxic-messages',
		title: 'Toxic messages',
		component: ToxicMessages,
	},
]
</script>

<template>
	<PageLayout active-tab="users" :tabs="pageTabs">
		<template #title>
			{{ t('adminPanel.title') }}
		</template>
	</PageLayout>
</template>
