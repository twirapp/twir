<script setup lang="ts">
import PageLayout, { type PageLayoutTab } from '~~/layers/dashboard/layout/page-layout.vue'
import UserAccountSettings from './user-account-settings.vue'
import UserLinkedAccounts from './user-linked-accounts.vue'
import UserPublicSettings from './user-public-settings.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const { t } = useI18n()

const pageTabs: PageLayoutTab[] = [
	{
		name: 'account',
		title: t('userSettings.account.title'),
		component: UserAccountSettings,
	},
	{
		name: 'public',
		title: t('userSettings.public.title'),
		component: UserPublicSettings,
	},
	{
		name: 'linked-accounts',
		title: 'Linked Accounts',
		component: UserLinkedAccounts,
	},
]
</script>

<template>
	<PageLayout
		active-tab="account"
		:tabs="pageTabs"
	>
		<template #title>
			{{ t('userSettings.title') }}
		</template>
	</PageLayout>
</template>
