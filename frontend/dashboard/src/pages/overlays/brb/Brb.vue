<script setup lang="ts">
import type { Settings } from '@twir/api/messages/overlays_be_right_back/overlays_be_right_back';
import { Font, FontSelector } from '@twir/fontsource';
import {
	useThemeVars,
	NButton,
	NColorPicker,
	NDivider,
	NInputNumber,
	NInput,
	NSwitch,
	useNotification,
	NAlert,
} from 'naive-ui';
import { ref, computed, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';


import { useBeRightBackOverlayManager, useProfile, useUserAccessFlagChecker } from '@/api';
import commandButton from '@/components/commandButton.vue';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js';

const themeVars = useThemeVars();
const { t } = useI18n();

const { data: profile } = useProfile();

const defaultSettings = {
	backgroundColor: 'rgba(9, 8, 8, 0.50)',
	fontColor: '#fff',
	fontFamily: 'inter',
	fontSize: 100,
	text: 'AFK FOR',
	late: {
		text: 'LATE FOR',
		displayBrbTime: true,
		enabled: true,
	},
	opacity: 50,
};

const formValue = ref<Settings>(structuredClone(defaultSettings));

const manager = useBeRightBackOverlayManager();
const {
	data: settings,
	isError: isSettingsError,
	isLoading: isSettingsLoading,
} = manager.getSettings();
const updater = manager.updateSettings();

watch(settings, (v) => {
	if (!v) return;
	formValue.value = toRaw(v);
}, { immediate: true });

const brbIframeRef = ref<HTMLIFrameElement | null>(null);
const brbIframeUrl = computed(() => {
	if (!profile.value) return null;

	return `${window.location.origin}/overlays/${profile.value.apiKey}/brb`;
});

const sendIframeMessage = (key: string, data?: any) => {
	if (!brbIframeRef.value) return;
	const win = brbIframeRef.value;

	win.contentWindow?.postMessage(JSON.stringify({
		key,
		data: toRaw(data),
	}));
};

const sendSettings = () => {
	sendIframeMessage('settings', {
		...toRaw(formValue.value),
		channelName: profile.value?.login,
		channelId: profile.value?.id,
	});
};

watch(brbIframeRef, (v) => {
	if (!v) return;

	v.contentWindow?.addEventListener('message', (e) => {
		const parsed = JSON.parse(e.data);
		if (parsed.key !== 'getSettings') return;

		sendSettings();
	});
});

watch(() => formValue, () => {
	if (!brbIframeRef.value) return;

	sendSettings();
}, { deep: true });

const { copyOverlayLink } = useCopyOverlayLink('brb');

const message = useNotification();

async function save() {
	await updater.mutateAsync(formValue.value);

	message.success({
		title: t('sharedTexts.saved'),
		duration: 5000,
	});
}

const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays;
});

function setDefaultSettings() {
	formValue.value = structuredClone(defaultSettings);
}

// TODO: fontWeight and fontStyle should be a select
const fontData = ref<Font | null>(null);
watch(() => fontData.value, (font) => {
	if (!font) return;
	formValue.value.fontFamily = font.id;
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
					@click="copyOverlayLink()"
				>
					{{ t('overlays.copyOverlayLink') }}
				</n-button>
				<n-button
					secondary
					type="success"
					@click="save"
				>
					{{ t('sharedButtons.save') }}
				</n-button>
			</div>

			<div class="card-body">
				<div class="card-body-column">
					<n-divider class="!my-0">
						{{ t('overlays.brb.settings.main.label') }}
					</n-divider>

					<div class="item">
						<div class="flex flex-col gap-1">
							<command-button
								name="brb"
								:title="t('overlays.brb.settings.main.startCommand.description')"
							/>
							<n-alert type="info" :show-icon="false">
								<span v-html="t('overlays.brb.settings.main.startCommand.example')" />
							</n-alert>
						</div>

						<command-button
							name="brbstop"
							:title="t('overlays.brb.settings.main.stopCommand.description')"
						/>
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.text') }}</span>
						<n-input v-model:value="formValue.text" :maxlength="500" />
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.background') }}</span>
						<n-color-picker
							v-model:value="formValue.backgroundColor" :modes="['rgb']"
							show-preview
						/>
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.font.color') }}</span>
						<n-color-picker
							v-model:value="formValue.fontColor" :modes="['hex', 'rgb']"
							:show-alpha="false"
						/>
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.font.family') }}</span>
						<font-selector
							v-model:font="fontData"
							:font-family="formValue.fontFamily"
							font-style="normal"
							:font-weight="400"
						/>
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.font.size') }}</span>
						<n-input-number v-model:value="formValue.fontSize" :min="1" :max="500" />
					</div>
				</div>

				<div class="card-body-column">
					<n-divider class="!my-0">
						{{ t('overlays.brb.settings.late.label') }}
					</n-divider>

					<div class="item">
						<span>{{ t('overlays.brb.settings.late.text') }}</span>
						<n-input v-model:value="formValue.late!.text" :maxlength="500" />
					</div>

					<div class="flex gap-2">
						<n-switch v-model:value="formValue.late!.enabled" />
						<span>{{ t('sharedTexts.enabled') }}</span>
					</div>

					<div class="flex gap-2">
						<n-switch v-model:value="formValue.late!.displayBrbTime" />
						<span>{{ t('overlays.brb.settings.late.displayBrb') }}</span>
					</div>
				</div>
			</div>
		</div>

		<div>
			<iframe
				v-if="brbIframeUrl"
				ref="brbIframeRef"
				:src="brbIframeUrl"
				class="iframe"
				border="0"
			/>
			<div class="absolute top-9 right-10 font-medium">
				<div class="flex gap-2">
					<n-button secondary size="small" type="warning" @click="sendIframeMessage('stop')">
						{{ t('overlays.brb.preview.stop') }}
					</n-button>
					<n-button
						secondary
						size="small"
						type="success"
						@click="() => {
							sendSettings();
							sendIframeMessage('start', { minutes: 0.1 })
						}"
					>
						{{ t('overlays.brb.preview.start') }}
					</n-button>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
@import '../styles.css';

.card {
	background-color: v-bind('themeVars.cardColor');
}

.iframe {
	border: 1px solid v-bind('themeVars.borderColor');
}
</style>
