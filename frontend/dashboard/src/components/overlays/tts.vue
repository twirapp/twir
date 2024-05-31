<script setup lang="ts">
import { NModal, NTabPane, NTabs } from 'naive-ui'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import Card from './card.vue'
import TTSSettings from './tts/settings.vue'
import UsersSettings from './tts/users.vue'

import { useCommandsApi } from '@/api/commands/commands'
import { useTtsOverlayManager } from '@/api/index.js'
import VoiceMessageIcon from '@/assets/overlays/voice-message.svg?use'
import CommandsList from '@/features/commands/ui/list.vue'

const commandsManager = useCommandsApi()
const { data: commands } = commandsManager.useQueryCommands()
const ttsCommands = computed(() => {
	return commands.value?.commands.filter((c) => c.module === 'TTS') ?? []
})

const ttsManager = useTtsOverlayManager()
const { data: settings, isError } = ttsManager.getSettings()

const isModalOpened = ref(false)

const { t } = useI18n()
</script>

<template>
	<Card
		title="Text to speech"
		:icon="VoiceMessageIcon"
		:icon-stroke="2"
		:description="t('overlays.tts.description')"
		overlay-path="tts"
		:copy-disabled="!settings || isError"
		@open-settings="isModalOpened = true"
	>
	</Card>

	<NModal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="TTS"
		content-style="padding: 0px; width: 100%"
		style="width: 800px; max-width: calc(100vw - 40px)"
	>
		<NTabs
			default-value="settings"
			justify-content="space-evenly"
			type="line"
			pane-style="padding-top: 0px"
		>
			<NTabPane name="settings" :tab="t('overlays.tts.tabs.general')">
				<TTSSettings />
			</NTabPane>
			<NTabPane name="users" :tab="t('overlays.tts.tabs.usersSettings')">
				<UsersSettings />
			</NTabPane>
			<NTabPane name="commands" :tab="t('sidebar.commands.label')">
				<CommandsList :commands="ttsCommands" />
			</NTabPane>
		</NTabs>
	</NModal>
</template>
