<script setup lang="ts">
import type { PageLayoutTab } from '~~/layers/dashboard/layout/page-layout.vue'

import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

import GiveawaysPagesHistory from '@/features/giveaways/pages/giveaways-pages-history.vue'
import GiveawaysPagesList from '@/features/giveaways/pages/giveaways-pages-list.vue'
import GiveawaysSettingsDialog from '@/features/giveaways/ui/giveaways-settings-dialog.vue'

const { t } = useI18n()

// Define tabs for the main giveaways page
const tabs: PageLayoutTab[] = [
	{
		name: 'list',
		component: GiveawaysPagesList,
		title: t('giveaways.pages.list.title'),
	},
	{
		name: 'history',
		component: GiveawaysPagesHistory,
		title: t('giveaways.pages.history.title'),
	},
]
</script>

<template>
	<!-- Otherwise show the main giveaways page with tabs -->
	<PageLayout
		:tabs="tabs"
		active-tab="list"
	>
		<template #action>
			<GiveawaysSettingsDialog />
		</template>

		<template #title>
			{{ t('giveaways.title') }}
		</template>
	</PageLayout>
</template>
