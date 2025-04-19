<script setup lang="ts">
import { h } from 'vue'
import { useI18n } from 'vue-i18n'

import type { PageLayoutTab } from '@/layout/page-layout.vue'

import GiveawaysPagesCurrentGiveaway
	from '@/features/giveaways/pages/giveaways-pages-current-giveaway.vue'
import GiveawaysPagesHistory from '@/features/giveaways/pages/giveaways-pages-history.vue'
import PageLayout from '@/layout/page-layout.vue'

const { t } = useI18n()

const tabs: PageLayoutTab[] = [
	{
		name: 'current',
		component: () => h(GiveawaysPagesCurrentGiveaway),
		title: t('giveaways.pages.current.title'),
	},
	{
		name: 'history',
		component: () => h(GiveawaysPagesHistory),
		title: t('giveaways.pages.history.title'),
	},
]
</script>

<template>
	<PageLayout :tabs="tabs" active-tab="current">
		<template #title>
			Giveaways
		</template>
	</PageLayout>
</template>
