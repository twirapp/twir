<script setup lang='ts'>
import { IconMessageCircle } from '@tabler/icons-vue';
import {
  NCard,
	NSpace,
	NText,
	NSkeleton,
	NModal,
	NButton,
	NTabs,
	NTabPane,
	useMessage,
} from 'naive-ui';
import { ref, computed } from 'vue';

import TTSSettings from './tts/settings.vue';
import UsersSettings from './tts/users.vue';

import { useTtsOverlayManager, useCommandsManager, useProfile } from '@/api/index.js';
import CommandsList from '@/components/commands/list.vue';

const ttsManager = useTtsOverlayManager();
const ttsSettings = ttsManager.getSettings();

const commandsManager = useCommandsManager();
const allCommands = commandsManager.getAll({});
const ttsCommands = computed(() => {
	return allCommands.data.value?.commands.filter(c => c.module === 'TTS') ?? [];
});

const userProfile = useProfile();
const overlayLink = computed(() => {
	return `${window.location.origin}/overlays/${userProfile.data?.value?.apiKey}/tts`;
});

const messages = useMessage();
const copyOverlayLink = () => {
	navigator.clipboard.writeText(overlayLink.value);
	messages.success('Copied link url, paste it in obs as browser source');
	return overlayLink;
};

const isModalOpened = ref(false);
</script>

<template>
  <n-card
    class="overlay-item"
    content-style="padding: 0px" @click="isModalOpened = true"
  >
    <n-skeleton v-if="ttsSettings.isLoading.value" size="large" :repeat="4" />
    <n-space v-else vertical align="center">
      <IconMessageCircle style="width: 112px; height: 112px" />
      <n-text strong style="font-size: 50px">
        TTS
      </n-text>
    </n-space>
  </n-card>

  <n-modal
    v-model:show="isModalOpened"
    :mask-closable="false"
    :segmented="true"
    preset="card"
    title="TTS"
    content-style="padding: 0px; width: 100%"
    style="width: 800px; max-width: calc(100vw - 40px);"
  >
    <template #header-extra>
      <n-button secondary type="success" @click="copyOverlayLink">
        Copy link url
      </n-button>
    </template>

    <n-tabs default-value="settings" justify-content="space-evenly" type="line" pane-style="padding-top: 0px">
      <n-tab-pane name="settings" tab="Settings">
        <TTSSettings />
      </n-tab-pane>
      <n-tab-pane name="users" tab="Users settings">
        <UsersSettings />
      </n-tab-pane>
      <n-tab-pane name="commands" tab="Commands">
        <commands-list :commands="ttsCommands" :show-header="false" />
      </n-tab-pane>
    </n-tabs>
  </n-modal>
</template>
