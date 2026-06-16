<script setup lang="ts">
import type { PageLayoutTab } from '~~/layers/dashboard/layout/page-layout.vue'

import { computed } from 'vue'
import TtsCommands from '~~/layers/dashboard/features/overlays/tts/ui/commands.vue'
import TtsGeneralSettings from '~~/layers/dashboard/features/overlays/tts/ui/general-settings.vue'
import TtsUsersSettings from '~~/layers/dashboard/features/overlays/tts/ui/users-settings.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

const { t } = useI18n()

const tabs = computed<PageLayoutTab[]>(() => [
	{
		title: t('overlays.tts.tabs.settings'),
		component: TtsGeneralSettings,
		name: 'settings',
	},
	{
		title: t('overlays.tts.tabs.users'),
		component: TtsUsersSettings,
		name: 'users',
	},
	{
		title: t('overlays.tts.tabs.commands'),
		component: TtsCommands,
		name: 'commands',
	},
])
</script>

<template>
	<PageLayout
		:tabs="tabs"
		active-tab="settings"
	>
		<template #title>
			{{ t('overlays.tts.title') }}
		</template>
	</PageLayout>
</template>
