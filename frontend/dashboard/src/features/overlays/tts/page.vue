<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { PageLayoutTab } from '@/layout/page-layout.vue'

import PageLayout from '@/layout/page-layout.vue'
import TtsGeneralSettings from '@/features/overlays/tts/ui/general-settings.vue'
import TtsUsersSettings from '@/features/overlays/tts/ui/users-settings.vue'
import TtsCommands from '@/features/overlays/tts/ui/commands.vue'

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
	<PageLayout :tabs="tabs" active-tab="settings">
		<template #title>
			{{ t('overlays.tts.title') }}
		</template>
	</PageLayout>
</template>

