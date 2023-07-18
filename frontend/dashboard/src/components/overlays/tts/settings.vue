<script setup lang='ts'>
import type { GetResponse as TTSSettings } from '@twir/grpc/generated/api/api/modules_tts';
import {
	NSlider,
	NSpace,
	NButton,
	NSkeleton,
	NSwitch,
	NAlert,
	NForm,
	NDivider,
	NSelect,
	NFormItem,
	NText,
	NGrid,
	NGridItem,
	NRow,
	useMessage,
} from 'naive-ui';
import { computed, ref, watch, toRaw } from 'vue';

import { useTtsOverlayManager } from '@/api/index.js';

const ttsManager = useTtsOverlayManager();
const ttsSettings = ttsManager.getSettings();
const ttsUpdater = ttsManager.updateSettings();
const ttsInfo = ttsManager.getInfo();

const countriesMapping: Record<string, string> = {
	'ru': '🇷🇺 Russian',
	'mk': '🇲🇰 Macedonian',
	'uk': '🇺🇦 Ukrainian',
	'ka': '🇬🇪 Georgian',
	'ky': '🇰🇬 Kyrgyz',
	'en': '🇺🇸 English',
	'pt': '🇵🇹 Portuguese',
	'eo': '🇺🇳 Esperanto',
	'sq': '🇦🇱 Albanian',
	'cs': '🇨🇿 Czech',
	'pl': '🇵🇱 Polish',
	'br': '🇧🇷 Brazilian',
};

type Voice = { label: string, value: string, key: string }
type VoiceGroup = Omit<Voice, 'value' | 'gender'> & { children: Voice[], type: 'group' }
const voicesOptions = computed<VoiceGroup[]>(() => {
	if (!ttsInfo.data.value?.voicesInfo) return [];

	const voices: Record<string, VoiceGroup> = {};

	for (const [, voice] of Object.entries(ttsInfo.data.value.voicesInfo)) {
		let lang = voice.lang;

		if (voice.lang === 'tt') {
			lang = 'ru';
		}

		if (!voices[lang]) {
			voices[lang] = {
				key: lang,
				label: `${countriesMapping[lang] ?? ''}`,
				type: 'group',
				children: [],
			};
		}

		voices[lang].children.push({
			key: lang,
			value: voice.name,
			label: `${voice.name} (${voice.gender})`,
		});
	}

	return Object.entries(voices).map(([, group]) => group);
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

watch(ttsSettings.data, (v) => {
	if (!v?.data) return;
	formValue.value = toRaw(v.data);
}, { immediate: true });

const message = useMessage();

async function save() {
	await ttsUpdater.mutateAsync({ data: formValue.value });
	message.success('Settings updated');
}
</script>

<template>
  <n-space vertical style="padding: 20px;">
    <n-alert type="info">
      Hint: you can use events system to trigger tts on reward.
    </n-alert>

    <n-skeleton v-if="!formValue || ttsSettings.isLoading.value" :sharp="false" size="large" />

    <n-form v-else style="margin-top: 15px">
      <n-grid cols="1 s:1 m:2 l:2" responsive="screen" :x-gap="20" :y-gap="20">
        <n-grid-item :span="1">
          <n-space justify="space-between">
            <n-text>Enabled</n-text>
            <n-switch v-model:value="formValue.enabled" />
          </n-space>
        </n-grid-item>

        <n-grid-item :span="1">
          <n-row justify-content="space-between" align-items="flex-start" style="flex-wrap: nowrap">
            <n-text>Allow users use different voices in main (!tts) command</n-text>
            <n-switch v-model:value="formValue.allowUsersChooseVoiceInMainCommand" />
          </n-row>
        </n-grid-item>

        <n-grid-item :span="1">
          <n-space justify="space-between">
            <n-text>Do not read emoji</n-text>
            <n-switch v-model:value="formValue.doNotReadEmoji" />
          </n-space>
        </n-grid-item>

        <n-grid-item :span="1">
          <n-space justify="space-between">
            <n-text>Do not read twitch emotes. Including 7tv, ffz, bttv</n-text>
            <n-switch v-model:value="formValue.doNotReadTwitchEmotes" />
          </n-space>
        </n-grid-item>

        <n-grid-item :span="1">
          <n-space justify="space-between">
            <n-text>Do not read links</n-text>
            <n-switch v-model:value="formValue.doNotReadLinks" />
          </n-space>
        </n-grid-item>

        <n-grid-item>
          <n-space justify="space-between">
            <n-text>Read all chat messages in tts</n-text>
            <n-switch v-model:value="formValue.readChatMessages" />
          </n-space>
        </n-grid-item>

        <n-grid-item :span="1">
          <n-space justify="space-between">
            <n-text>Read nicknames when reading tts</n-text>
            <n-switch v-model:value="formValue.readChatMessagesNicknames" />
          </n-space>
        </n-grid-item>
      </n-grid>

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

      <n-space style="width:100%" vertical size="small">
        <n-form-item label="Volume" size="small">
          <n-slider v-model:value="formValue.volume" :step="1" />
        </n-form-item>
        <n-form-item label="Pitch" size="small">
          <n-slider v-model:value="formValue.pitch" :step="1" />
        </n-form-item>
        <n-form-item label="Rate" size="small">
          <n-slider v-model:value="formValue.rate" :step="1" />
        </n-form-item>
      </n-space>
    </n-form>

    <n-button secondary type="success" block style="margin-top: 10px" @click="save">
      Save
    </n-button>
  </n-space>
</template>

<style scoped lang='postcss'>

</style>
