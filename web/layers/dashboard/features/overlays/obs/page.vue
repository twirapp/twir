<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { PageLayoutTab } from '@/layout/page-layout.vue'

import PageLayout from '@/layout/page-layout.vue'
import ObsSettings from '@/features/overlays/obs/settings.vue'

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
	<PageLayout :tabs="tabs" active-tab="settings">
		<template #title>
			{{ t('overlays.obs.title') }}
		</template>
	</PageLayout>
</template>
