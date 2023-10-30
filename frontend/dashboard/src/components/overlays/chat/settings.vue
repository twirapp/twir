<script setup lang="ts">
import { faker } from '@faker-js/faker';
import {
	ChatBox,
	type Message,
	type Settings as ChatBoxSettings,
	BadgeVersion,
} from '@twir/frontend-chat';
import type {
  Settings,
} from '@twir/grpc/generated/api/api/modules_chat_overlay';
import { useIntervalFn } from '@vueuse/core';
import {
	useNotification,
	NButton,
	NSwitch,
	NSlider,
	NSelect,
} from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import '@twir/frontend-chat/style.css';

import { globalBadges } from './constants.js';

import { useChatOverlayManager } from '@/api/index.js';

const chatManager = useChatOverlayManager();
const { data: settings } = chatManager.getSettings();
const updater = chatManager.updateSettings();

const globalBadgesObject = Object.fromEntries(globalBadges);

const messagesMock = ref<Message[]>([]);
useIntervalFn(() => {
	messagesMock.value.push({
		sender: faker.person.firstName(),
		chunks: [{
			type: 'text',
			value: faker.lorem.words({ min: 1, max: 20 }),
		}],
		createdAt: new Date(),
		internalId: crypto.randomUUID(),
		isAnnounce: faker.datatype.boolean(),
		isItalic: false,
		type: 'message',
		senderColor: faker.color.rgb(),
		announceColor: '',
		badges: {
			[faker.helpers.objectKey(globalBadgesObject)]: '1',
		},
		id: crypto.randomUUID(),
		senderDisplayName: faker.person.lastName(),
	});
}, 1 * 1000);

const formValue = ref<Settings>({
	fontSize: 20,
	hideBots: false,
	hideCommands: false,
	messageHideTimeout: 0,
	messageShowDelay: 0,
	preset: 'clean',
});

const chatBoxSettings = computed<ChatBoxSettings>(() => {
	return {
		channelId: '',
    channelName: '',
    channelDisplayName: '',
    globalBadges,
    channelBadges: new Map<string, BadgeVersion>(),
		...formValue.value,
	};
});

watch(() => settings.value, () => {
	if (!settings.value) return;

	formValue.value = toRaw(settings.value);
}, { immediate: true });

const message = useNotification();
const { t } = useI18n();

async function save() {
	if (!formValue.value) return;

	await updater.mutateAsync(formValue.value);
	message.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	});
}

const sliderMarks = {
	0: '0',
	60: '60',
};

const selectOptions = [
	{ label: 'Clean', value: 'clean' },
	{ label: 'Boxed', value: 'boxed' },
];
</script>

<template>
	<div class="settings">
		<div class="form">
			<div>
				<span>{{ t('overlays.chat.style') }}</span>
				<n-select v-model:value="formValue.preset" :options="selectOptions" />
			</div>

			<div class="switch">
				<n-switch v-model:value="formValue.hideBots" />
				<span>{{ t('overlays.chat.hideBots') }}</span>
			</div>

			<div class="switch">
				<n-switch v-model:value="formValue.hideCommands" />
				<span>{{ t('overlays.chat.hideCommands') }}</span>
			</div>

			<div class="slider">
				<span>{{ t('overlays.chat.fontSize') }}</span>
				<n-slider v-model:value="formValue.fontSize" :min="12" :max="80" :marks="{ 12: '12', 80: '80'}" />
			</div>

			<div class="slider">
				<span>{{ t('overlays.chat.hideTimeout') }}</span>
				<n-slider v-model:value="formValue.messageHideTimeout" :max="60" :marks="sliderMarks" />
			</div>

			<div class="slider">
				<span>{{ t('overlays.chat.showDelay') }}</span>
				<n-slider v-model:value="formValue.messageShowDelay" :max="60" :marks="sliderMarks" />
			</div>

			<n-button secondary type="success" block style="margin-top: 10px" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</div>
		<div class="chatBox">
			<ChatBox
				class="chatBox"
				:messages="messagesMock"
				:settings="chatBoxSettings"
			/>
		</div>
	</div>
</template>

<style scoped>
.chatBox {
	max-height: 40vw;
	overflow: hidden;
}

.settings {
	display: flex;
	gap: 24px;
}

.settings .form {
	display: flex;
	flex-direction: column;
	gap: 12px;
	min-width: 25vw;
}

.input {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.switch {
	display: flex;
	gap: 8px;
	align-items: start;
}

.slider {
	display: flex;
	gap: 4px;
	align-items: center;
	flex-direction: column;
}
</style>
