<script setup lang='ts'>
import { IconMessageCircle } from '@tabler/icons-vue';
import type { GetResponse as TTSSettings } from '@twir/grpc/generated/api/api/modules_tts';
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
	NForm,
	NFormItem,
	NSwitch,
	NAlert,
	NSelect,
	NDivider,
	NSlider,
} from 'naive-ui';
import { ref, computed } from 'vue';


import { useTtsOverlayManager, useCommandsManager, useProfile } from '@/api/index.js';
import CommandsList from '@/components/commands/list.vue';
import UsersSettings from '@/components/overlays/tts/users.vue';

const ttsManager = useTtsOverlayManager();
const ttsSettings = ttsManager.getSettings();
const ttsUpdater = ttsManager.updateSettings();
const ttsInfo = ttsManager.getInfo();

type Voices = Array<{ label: string, value: string }>
const voicesOptions = computed<Voices>(() => {
	if (!ttsInfo.data.value?.voicesInfo) return [];

	const voicesArray = Object.entries(ttsInfo.data.value.voicesInfo)
		.sort(([, valueA], [, valueB]) => valueA.country.localeCompare(valueB.country))
		.reduce((result, [key, value]) => {
			result.push({
				label: `[${value.country}] ${value.name}`,
				value: key,
			});
			return result;
		}, [] as Voices);


	return voicesArray;
});

const formValue = ref<TTSSettings['data']>({
	enabled: true,
	voice: 'alan',
	disallowedVoices: [],
	pitch: 50,
	rate: 50,
	volume: 30,
	doNotReadTwitchEmotes: true,
	doNotReadEmoji: true,
	doNotReadLinks: true,
	allowUsersChooseVoiceInMainCommand: false,
	maxSymbols: 0,
	readChatMessages: false,
	readChatMessagesNicknames: false,
});

async function save() {
	await ttsUpdater.mutateAsync({ data: formValue.value });
}

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
        <n-space vertical style="padding: 20px;">
          <n-alert type="info">
            Hint: you can use events system to trigger tts on reward.
          </n-alert>

          <n-skeleton v-if="!formValue || ttsSettings.isLoading.value" :sharp="false" size="large" />

          <n-form v-else>
            <n-space justify="space-between">
              <n-text>Enabled</n-text>
              <n-switch v-model:value="formValue.enabled" />
            </n-space>

            <n-space justify="space-between">
              <n-text>Allow users use different voices in main (!tts) command</n-text>
              <n-switch v-model:value="formValue.allowUsersChooseVoiceInMainCommand" />
            </n-space>

            <n-space justify="space-between">
              <n-text>Do not read emoji</n-text>
              <n-switch v-model:value="formValue.doNotReadEmoji" />
            </n-space>

            <n-space justify="space-between">
              <n-text>Do not read twitch emotes. Including 7tv, ffz, bttv</n-text>
              <n-switch v-model:value="formValue.doNotReadTwitchEmotes" />
            </n-space>

            <n-space justify="space-between">
              <n-text>Do not read links</n-text>
              <n-switch v-model:value="formValue.doNotReadLinks" />
            </n-space>

            <n-space justify="space-between">
              <n-text>Read all chat messages in tts</n-text>
              <n-switch v-model:value="formValue.readChatMessages" />
            </n-space>

            <n-space justify="space-between">
              <n-text>Read nicknames when reading tts</n-text>
              <n-switch v-model:value="formValue.readChatMessagesNicknames" />
            </n-space>

            <n-divider />

            <n-form-item label="Voice" show-require-mark>
              <n-select
                v-model:value="formValue.voice"
                remote
                :loading="ttsInfo.isLoading.value"
                :options="voicesOptions"
              />
            </n-form-item>

            <n-form-item label="Disallowed for usage voices">
              <n-select
                v-model:value="formValue.disallowedVoices"
                remote
                clearable
                :loading="ttsInfo.isLoading.value"
                :options="voicesOptions"
                multiple
              />
            </n-form-item>

            <n-space style="width:100%" vertical>
              <n-form-item label="Volume">
                <n-slider v-model:value="formValue.volume" :step="1" />
              </n-form-item>
              <n-form-item label="Pitch">
                <n-slider v-model:value="formValue.pitch" :step="1" />
              </n-form-item>
              <n-form-item label="Rate">
                <n-slider v-model:value="formValue.rate" :step="1" />
              </n-form-item>
            </n-space>
          </n-form>

          <n-button secondary type="success" block style="margin-top: 10px" @click="save">
            Save
          </n-button>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="users" tab="Users settings">
        <users-settings />
      </n-tab-pane>
      <n-tab-pane name="commands" tab="Commands">
        <commands-list :commands="ttsCommands" :show-header="false" />
      </n-tab-pane>
    </n-tabs>
  </n-modal>
</template>
