<script setup lang="ts">
import UserAccountSettings from '~/layers/dashboard/pages/settings/user-account-settings.vue'
import UserPublicSettings from '~/layers/dashboard/pages/settings/user-public-settings.vue'
import UserLinkedAccounts from '~/layers/dashboard/pages/settings/user-linked-accounts.vue'

import PageLayout, { type PageLayoutTab } from '~/layers/dashboard/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { t } = useI18n()

const pageTabs: PageLayoutTab[] = [
	{
		name: 'account',
		title: t('userSettings.account.title'),
		component: UserAccountSettings
	},
	{
		name: 'public',
		title: t('userSettings.public.title'),
		component: UserPublicSettings
	},
	{
		name: 'linked-accounts',
		title: 'Linked Accounts',
		component: UserLinkedAccounts
	}
]
</script>

<template>
	<PageLayout active-tab="account" :tabs="pageTabs">
		<template #title>
			{{ t('userSettings.title') }}
		</template>
	</PageLayout>
</template>
