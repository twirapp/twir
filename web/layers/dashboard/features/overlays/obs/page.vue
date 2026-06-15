<script setup lang="ts">
import type { PageLayoutTab } from '~~/layers/dashboard/layout/page-layout.vue'

import { computed } from 'vue'
import ObsSettings from '~~/layers/dashboard/features/overlays/obs/settings.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

const { t } = useI18n()

const tabs = computed<PageLayoutTab[]>(() => [
	{
		title: t('overlays.obs.tabs.settings'),
		component: ObsSettings,
		name: 'settings',
	},
])
</script>

<template>
	<PageLayout
		:tabs="tabs"
		active-tab="settings"
	>
		<template #title>
			{{ t('overlays.obs.title') }}
		</template>
	</PageLayout>
</template>
