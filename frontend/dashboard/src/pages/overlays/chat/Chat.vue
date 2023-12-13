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
	NDivider,
	NColorPicker,
} from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { globalBadges } from './constants.js';
import * as faker from './faker.js';

import {
	useChatOverlayManager,
	useGoogleFontsList,
	useProfile,
	useUserAccessFlagChecker,
} from '@/api/index.js';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink';

const chatManager = useChatOverlayManager();
const {
	data: settings,
	isLoading: isSettingsLoading,
	isError: isSettingsError,
} = chatManager.getSettings();
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

const defaultSettings: Settings = {
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
	textShadowColor: 'rgba(0,0,0,1)',
	textShadowSize: 0,
	chatBackgroundColor: 'rgba(0, 0, 0, 0)',
};

const formValue = ref<Settings>(structuredClone(defaultSettings));

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

const setDefaultSettings = () => {
	formValue.value = structuredClone(defaultSettings);
};

const { copyOverlayLink } = useCopyOverlayLink('chat');
const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const { data: profile } = useProfile();
const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays;
});
</script>

<template>
	<div class="page">
		<div class="card">
			<div class="card-header">
				<n-button
					secondary
					type="error"
					@click="setDefaultSettings"
				>
					{{ t('sharedButtons.setDefaultSettings') }}
				</n-button>
				<n-button
					secondary
					type="info"
					:disabled="isSettingsError || isSettingsLoading || !canCopyLink"
					@click="copyOverlayLink"
				>
					{{ t('overlays.copyOverlayLink') }}
				</n-button>
				<n-button secondary type="success" @click="save">
					{{ t('sharedButtons.save') }}
				</n-button>
			</div>

			<div class="card-body">
				<div class="card-body-column">
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
						<n-switch
							v-model:value="formValue.showAnnounceBadge"
							:disabled="formValue.preset === 'clean' || !formValue.showBadges"
						/>
						<span>{{ t('overlays.chat.showAnnounceBadge') }}</span>
					</div>

					<n-divider />

					<div style="display: flex; flex-direction: column; gap: 4px;">
						<div style="display: flex; justify-content: space-between;">
							<span>{{ t('overlays.chat.fontFamily') }}</span>
							<n-button
								size="tiny" secondary type="success"
								@click="formValue.fontFamily = defaultFont"
							>
								<IconReload style="height: 15px;" />
								{{ t('overlays.chat.resetToDefault') }}
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
						<span>{{ t('overlays.chat.fontSize') }}({{ formValue.fontSize }}px)</span>
						<n-slider
							v-model:value="formValue.fontSize" :min="12" :max="80"
							:marks="{ 12: '12', 80: '80'}"
						/>
					</div>

					<div class="slider">
						<div style="display: flex; justify-content: space-between; margin-bottom: 4px;">
							<span>{{ t('overlays.chat.backgroundColor') }}</span>
							<n-button
								size="tiny" secondary type="success"
								@click="formValue.chatBackgroundColor = defaultSettings.chatBackgroundColor"
							>
								<IconReload style="height: 15px;" />
								{{ t('overlays.chat.resetToDefault') }}
							</n-button>
						</div>
						<n-color-picker
							v-model:value="formValue.chatBackgroundColor"
							default-value="rgba(16, 16, 20, 1)"
						/>
					</div>


					<div class="slider">
						<span>{{ t('overlays.chat.textShadow') }}({{ formValue.textShadowSize }}px)</span>
						<n-color-picker
							v-model:value="formValue.textShadowColor"
							default-value="rgba(0,0,0,1)"
						/>
						<n-slider v-model:value="formValue.textShadowSize" :min="0" :max="30" />
					</div>

					<n-divider />

					<div class="slider">
						<span>{{ t('overlays.chat.hideTimeout') }}({{ formValue.messageHideTimeout }}s)</span>
						<n-slider v-model:value="formValue.messageHideTimeout" :max="60" :marks="sliderMarks" />
					</div>

					<div class="slider">
						<span>{{ t('overlays.chat.showDelay') }}({{ formValue.messageShowDelay }}s)</span>
						<n-slider v-model:value="formValue.messageShowDelay" :max="60" :marks="sliderMarks" />
					</div>
				</div>
			</div>
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
@import '../styles.css';

.card-body-column {
	width: 100%;
}

.card {
	background-color: v-bind('themeVars.cardColor');
}

:deep(.chat) {
	height: 80dvh;
}
</style>
