<script setup lang="ts">


import UserAccountSettings from './user-account-settings.vue'
import UserPublicSettings from './user-public-settings.vue'

import PageLayout, { type PageLayoutTab } from '~/layout/page-layout.vue'

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
