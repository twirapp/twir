<script setup lang="ts">

import { IconReload } from '@tabler/icons-vue';
import {
	ChatBox,
	type Message,
	type Settings as ChatBoxSettings,
	BadgeVersion,
} from '@twir/frontend-chat';
import type {
	Settings,
} from '@twir/grpc/generated/api/api/overlays_chat';
import { useIntervalFn } from '@vueuse/core';
import {
	useNotification,
	NButton,
	NSwitch,
	NSlider,
	NSelect,
	NTreeSelect,
	type TreeSelectOption,
	useThemeVars,
} from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import '@twir/frontend-chat/style.css';

import { globalBadges } from './constants.js';
import * as faker from './faker.js';

import { useChatOverlayManager, useGoogleFontsList } from '@/api/index.js';

const chatManager = useChatOverlayManager();
const { data: settings } = chatManager.getSettings();
const updater = chatManager.updateSettings();
const {
	data: googleFonts,
	isError: isGoogleFontsError,
	isLoading: isGoogleFontsLoading,
} = useGoogleFontsList();
const themeVars = useThemeVars();

const globalBadgesObject = Object.fromEntries(globalBadges);

const messagesMock = ref<Message[]>([]);
useIntervalFn(() => {
	const internalId = crypto.randomUUID();

	messagesMock.value.push({
		sender: faker.firstName(),
		chunks: [{
			type: 'text',
			value: faker.lorem(),
		}],
		createdAt: new Date(),
		internalId,
		isAnnounce: faker.boolean(),
		isItalic: false,
		type: 'message',
		senderColor: faker.rgb(),
		announceColor: '',
		badges: {
			[faker.randomObjectKey(globalBadgesObject)]: '1',
		},
		id: crypto.randomUUID(),
		senderDisplayName: faker.firstName(),
	});

	if (formValue.value.messageHideTimeout != 0) {
		setTimeout(() => {
			messagesMock.value = messagesMock.value.filter(m => m.internalId != internalId);
		}, formValue.value.messageHideTimeout * 1000);
	}
}, 1 * 1000);


const defaultFont = 'Roboto:700italic';
const formValue = ref<Settings>({
	fontSize: 20,
	hideBots: false,
	hideCommands: false,
	messageHideTimeout: 0,
	messageShowDelay: 0,
	preset: 'clean',
	fontFamily: defaultFont,
	showBadges: true,
	showAnnounceBadge: true,
	reverseMessages: false,
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

const styleSelectOptions = [
	{ label: 'Clean', value: 'clean' },
	{ label: 'Boxed', value: 'boxed' },
];

const fontSelectOptions = computed<TreeSelectOption[]>(() => {
	return googleFonts?.value?.fonts
		.map((f) => {
			const option: TreeSelectOption = {
				label: f.family,
				children: f.files.map((c) => ({
					label: `${f.family}:${c.name}`,
					key: `${f.family}:${c.name}`,
				})),
				key: f.family,
			};

			return option;
		}) ?? [];
});
</script>

<template>
	<div class="settings">
		<div class="form">
			<div>
				<span>{{ t('overlays.chat.style') }}</span>
				<n-select v-model:value="formValue.preset" :options="styleSelectOptions" />
			</div>

			<div class="switch">
				<n-switch v-model:value="formValue.hideBots" />
				<span>{{ t('overlays.chat.hideBots') }}</span>
			</div>

			<div class="switch">
				<n-switch v-model:value="formValue.hideCommands" />
				<span>{{ t('overlays.chat.hideCommands') }}</span>
			</div>

			<div class="switch">
				<n-switch v-model:value="formValue.reverseMessages" />
				<span>{{ t('overlays.chat.reverseMessages') }}</span>
			</div>

			<div class="switch">
				<n-switch v-model:value="formValue.showBadges" />
				<span>{{ t('overlays.chat.showBadges') }}</span>
			</div>

			<div class="switch">
				<n-switch v-model:value="formValue.showAnnounceBadge" :disabled="formValue.preset === 'clean' || !formValue.showBadges" />
				<span>{{ t('overlays.chat.showAnnounceBadge') }}</span>
			</div>

			<div style="display: flex; flex-direction: column; gap: 4px;">
				<div style="display: flex; justify-content: space-between;">
					<span>{{ t('overlays.chat.fontFamily') }}</span>
					<n-button size="tiny" secondary type="success" @click="formValue.fontFamily = defaultFont">
						<IconReload style="height: 15px;" /> {{ t('overlays.chat.revertFont') }}
					</n-button>
				</div>
				<n-tree-select
					v-model:value="formValue.fontFamily"
					filterable
					:options="fontSelectOptions"
					:loading="isGoogleFontsLoading"
					:disabled="isGoogleFontsError"
					check-strategy="child"
				>
					<template #action>
						{{ t('overlays.chat.fontFamilyDescription') }}
						<a
							class="action-link"
							href="https://fonts.google.com/"
							target="_blank"
							:style="{ color: themeVars.successColor }"
						>
							Preview Google Fonts
						</a>
					</template>
				</n-tree-select>
			</div>

			<div class="slider">
				<span>{{ t('overlays.chat.fontSize') }}</span>
				<n-slider
					v-model:value="formValue.fontSize" :min="12" :max="80"
					:marks="{ 12: '12', 80: '80'}"
				/>
			</div>

			<div class="slider">
				<span>{{ t('overlays.chat.hideTimeout') }}</span>
				<n-slider v-model:value="formValue.messageHideTimeout" :max="60" :marks="sliderMarks" />
			</div>

			<div class="slider">
				<span>{{ t('overlays.chat.showDelay') }}</span>
				<n-slider v-model:value="formValue.messageShowDelay" :max="60" :marks="sliderMarks" />
			</div>

			<n-button secondary type="success" block style="margin-top: auto" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</div>
		<div class="chatBox">
			<ChatBox
				class="chatBox"
				:messages="messagesMock"
				:settings="chatBoxSettings"
				:fonts="googleFonts?.fonts"
			/>
		</div>
	</div>
</template>

<style scoped>
.chatBox {
	height: 40dvh;
	overflow: hidden;
	width: 100%;
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

.action-link {
	text-decoration: none;
}
</style>
