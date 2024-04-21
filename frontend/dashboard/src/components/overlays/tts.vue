<script setup lang="ts">
import { NModal, NTabs, NTabPane } from 'naive-ui';
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import Card from './card.vue';
import TTSSettings from './tts/settings.vue';
import UsersSettings from './tts/users.vue';

import { useCommandsApi } from '@/api/commands/commands';
import { useTtsOverlayManager } from '@/api/index.js';
import VoiceMessageIcon from '@/assets/overlays/voice-message.svg?use';
import CommandsList from '@/features/commands/components/list.vue';

const commandsManager = useCommandsApi();
const { data: commands } = commandsManager.useQueryCommands();
const ttsCommands = computed(() => {
	return commands.value?.commands.filter((c) => c.module === 'TTS') ?? [];
});

const ttsManager = useTtsOverlayManager();
const { data: settings, isError } = ttsManager.getSettings();

const isModalOpened = ref(false);

const { t } = useI18n();
</script>

<template>
	<card
		title="Text to speech"
		:icon="VoiceMessageIcon"
		:icon-stroke="2"
		:description="t('overlays.tts.description')"
		overlay-path="tts"
		:copy-disabled="!settings || isError"
		@open-settings="isModalOpened = true"
	>
	</card>

	<n-modal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="TTS"
		content-style="padding: 0px; width: 100%"
		style="width: 800px; max-width: calc(100vw - 40px)"
	>
		<n-tabs
			default-value="settings"
			justify-content="space-evenly"
			type="line"
			pane-style="padding-top: 0px"
		>
			<n-tab-pane name="settings" :tab="t('overlays.tts.tabs.general')">
				<TTSSettings />
			</n-tab-pane>
			<n-tab-pane name="users" :tab="t('overlays.tts.tabs.usersSettings')">
				<UsersSettings />
			</n-tab-pane>
			<n-tab-pane name="commands" :tab="t('sidebar.commands.label')">
				<commands-list :commands="ttsCommands" />
			</n-tab-pane>
		</n-tabs>
	</n-modal>
</template>
