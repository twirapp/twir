<script setup lang="ts">
import { IconPlayerPlay } from '@tabler/icons-vue';
import type { GetResponse as TTSSettings } from '@twir/api/messages/modules_tts/modules_tts';
import {
	NAlert,
	NButton,
	NDivider,
	NForm,
	NFormItem,
	NGrid,
	NGridItem,
	NInput,
	NRow,
	NSelect,
	NSkeleton,
	NSlider,
	NSpace,
	NSwitch,
	NText,
	useMessage,
} from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useTtsOverlayManager } from '@/api/index.js';

const ttsManager = useTtsOverlayManager();
const ttsSettings = ttsManager.getSettings();
const ttsUpdater = ttsManager.updateSettings();
const ttsInfo = ttsManager.getInfo();
const ttsSay = ttsManager.useSay();

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

	for (const [voiceKey, voice] of Object.entries(ttsInfo.data.value.voicesInfo)) {
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
			value: voiceKey,
			label: `${voice.name} (${voice.gender})`,
		});
	}

	return Object.entries(voices).map(([, group]) => group);
});

const formValue = ref<TTSSettings['data']>({
	enabled: false,
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
const { t } = useI18n();


async function save() {
	await ttsUpdater.mutateAsync({ data: formValue.value });
	message.success(t('sharedTexts.saved'));
}

const previewText = ref('');

async function previewVoice() {
	if (!previewText.value || !formValue.value) return;

	await ttsSay.mutateAsync({
		voice: formValue.value.voice,
		text: previewText.value,
		volume: formValue.value.volume,
		pitch: formValue.value.pitch,
		rate: formValue.value.rate,
	});
}
</script>

<template>
	<n-space vertical class="p-5">
		<n-alert type="info">
			{{ t('overlays.tts.eventsHint') }}
		</n-alert>

		<n-skeleton v-if="!formValue || ttsSettings.isLoading.value" :sharp="false" size="large" />

		<n-form v-else class="mt-4">
			<n-grid cols="1 s:1 m:2 l:2" responsive="screen" :x-gap="20" :y-gap="20">
				<n-grid-item :span="1">
					<n-space justify="space-between">
						<n-text>{{ t('sharedTexts.enabled') }}</n-text>
						<n-switch v-model:value="formValue.enabled" />
					</n-space>
				</n-grid-item>

				<n-grid-item :span="1">
					<n-row justify-content="space-between" align-items="flex-start" class="flex-nowrap">
						<n-text>{{ t('overlays.tts.allowUsersChooseVoice') }}</n-text>
						<n-switch v-model:value="formValue.allowUsersChooseVoiceInMainCommand" />
					</n-row>
				</n-grid-item>

				<n-grid-item :span="1">
					<n-space justify="space-between">
						<n-text>{{ t('overlays.tts.doNotReadEmoji') }}</n-text>
						<n-switch v-model:value="formValue.doNotReadEmoji" />
					</n-space>
				</n-grid-item>

				<n-grid-item :span="1">
					<n-space justify="space-between">
						<n-text>{{ t('overlays.tts.doNotReadTwitchEmotes') }}</n-text>
						<n-switch v-model:value="formValue.doNotReadTwitchEmotes" />
					</n-space>
				</n-grid-item>

				<n-grid-item :span="1">
					<n-space justify="space-between">
						<n-text>{{ t('overlays.tts.doNotReadLinks') }}</n-text>
						<n-switch v-model:value="formValue.doNotReadLinks" />
					</n-space>
				</n-grid-item>

				<n-grid-item>
					<n-space justify="space-between">
						<n-text>{{ t('overlays.tts.readChatMessages') }}</n-text>
						<n-switch v-model:value="formValue.readChatMessages" />
					</n-space>
				</n-grid-item>

				<n-grid-item :span="1">
					<n-space justify="space-between">
						<n-text>{{ t('overlays.tts.readChatMessagesNicknames') }}</n-text>
						<n-switch v-model:value="formValue.readChatMessagesNicknames" />
					</n-space>
				</n-grid-item>
			</n-grid>

			<n-divider />

			<n-form-item :label="t('overlays.tts.voice')" show-require-mark>
				<n-select
					v-model:value="formValue.voice"
					remote
					:loading="ttsInfo.isLoading.value"
					:options="voicesOptions"
				/>
			</n-form-item>

			<n-form-item :label="t('overlays.tts.disallowedVoices')">
				<n-select
					v-model:value="formValue.disallowedVoices"
					remote
					clearable
					:loading="ttsInfo.isLoading.value"
					:options="voicesOptions"
					multiple
				/>
			</n-form-item>

			<n-space class="w-full" vertical size="small">
				<n-form-item :label="t('overlays.tts.volume')" size="small">
					<n-slider v-model:value="formValue.volume" :step="1" />
				</n-form-item>
				<n-form-item :label="t('overlays.tts.pitch')" size="small">
					<n-slider v-model:value="formValue.pitch" :step="1" />
				</n-form-item>
				<n-form-item :label="t('overlays.tts.rate')" size="small">
					<n-slider v-model:value="formValue.rate" :step="1" />
				</n-form-item>
			</n-space>

			<n-divider class="m-0 mb-2.5" />

			<n-form-item :label="`🎤 ${t('overlays.tts.previewText')}`">
				<div class="flex gap-1 w-full">
					<n-input
						v-model:value="previewText" :placeholder="t('overlays.tts.previewText')"
						class="w-1/2"
					/>
					<n-button text @click="previewVoice">
						<IconPlayerPlay />
					</n-button>
				</div>
			</n-form-item>
		</n-form>

		<n-button secondary type="success" block class="mt-2.5" @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-space>
</template>
