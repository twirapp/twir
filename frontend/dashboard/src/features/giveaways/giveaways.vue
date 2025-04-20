<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import type { PageLayoutTab } from '@/layout/page-layout.vue'

import GiveawaysPagesHistory from '@/features/giveaways/pages/giveaways-pages-history.vue'
import GiveawaysPagesList from '@/features/giveaways/pages/giveaways-pages-list.vue'
import PageLayout from '@/layout/page-layout.vue'

const { t } = useI18n()

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
	<PageLayout :tabs="tabs" active-tab="list">
		<template #title>
			Giveaways
		</template>
	</PageLayout>
</template>
