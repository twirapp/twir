<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';
import { useThemeVars, NButton, NColorPicker, NDivider, NInputNumber, NInput, NSwitch, NModal, useNotification } from 'naive-ui';
import { ref, computed, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCopyOverlayLink } from '../copyOverlayLink';

import { useBeRightBackOverlayManager, useProfile } from '@/api';
import commandButton from '@/components/commandButton.vue';
import FontSelector from '@/components/fontSelector.vue';

defineProps<{
	showSettings: boolean
}>();
defineEmits<{
	close: []
}>();

const themeVars = useThemeVars();
const { t } = useI18n();

const { data: profile } = useProfile();

const defaultSettings = {
	backgroundColor: 'rgba(9, 8, 8, 0.49)',
	fontColor: '#fff',
	fontFamily: '',
	fontSize: 100,
	text: 'AFK FOR',
	late: {
		text: 'LATE FOR',
		displayBrbTime: true,
		enabled: true,
	},
};

const formValue = ref<Settings>(defaultSettings);

const manager = useBeRightBackOverlayManager();
const { data: settings } = manager.getSettings();
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
</script>

<template>
	<n-modal
		:show="showSettings"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Afk overlay"
		content-style="padding: 10px; width: 100%"
		style="width: 50dvw;"
		footer-style="padding: 8px;"
		@close="$emit('close')"
	>
		<div class="settings">
			<div class="form">
				<div style="display: flex; flex-direction: column; gap: 12px;">
					<n-divider style="margin: 0">
						{{ t('overlays.brb.settings.main.label') }}
					</n-divider>

					<div class="item">
						<div>
							<command-button name="brb" :title="t('overlays.brb.settings.main.startCommand.description')" />
							<span v-html="t('overlays.brb.settings.main.startCommand.example')" />
						</div>
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.text') }}</span>
						<n-input v-model:value="formValue.text" :maxlength="500" />
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.background') }}</span>
						<n-color-picker v-model:value="formValue.backgroundColor" :modes="['rgb']" />
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.font.color') }}</span>
						<n-color-picker v-model:value="formValue.fontColor" :modes="['rgb']" />
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.font.family') }}</span>
						<font-selector v-model="formValue.fontFamily" :clearable="true" />
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.font.size') }}</span>
						<n-input-number v-model:value="formValue.fontSize" :min="1" :max="500" />
					</div>
				</div>

				<div style="display: flex; flex-direction: column; gap: 12px;">
					<n-divider style="margin: 0">
						{{ t('overlays.brb.settings.late.label') }}
					</n-divider>

					<div class="item">
						<span>{{ t('overlays.brb.settings.late.text') }}</span>
						<n-input v-model:value="formValue.late!.text" :maxlength="500" />
					</div>

					<div style="display: flex; gap: 8px">
						<n-switch v-model:value="formValue.late!.enabled" />
						<span>{{ t('sharedTexts.enabled') }}</span>
					</div>

					<div style="display: flex; gap: 8px">
						<n-switch v-model:value="formValue.late!.displayBrbTime" />
						<span>{{ t('overlays.brb.settings.late.displayBrb') }}</span>
					</div>
				</div>
			</div>
			<div>
				<div style="position: absolute; top: 85px; right: 20px; font-weight: 500;">
					<div style="display: flex; gap: 8px">
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
				<iframe
					v-if="brbIframeUrl"
					ref="brbIframeRef"
					:src="brbIframeUrl"
					class="iframe"
				/>
			</div>
		</div>
		<template #footer>
			<div class="footer">
				<n-button
					secondary
					type="error"
					@click="formValue = defaultSettings"
				>
					{{ t('sharedButtons.setDefaultSettings') }}
				</n-button>

				<div style="display: flex; gap: 8px;">
					<n-button
						secondary
						type="info"
						@click="copyOverlayLink"
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
			</div>
		</template>
	</n-modal>
</template>

<style scoped>
.settings {
	display: flex;
	gap: 16px;
	width: 100%;
}

.form {
	padding: 8px;
	border-radius: 8px;
	background-color: v-bind('themeVars.cardColor');
	display: flex;
	gap: 8px;
}

.form > div {
	width: 50%;
}

.form .item {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.settings > div {
	width: 50%;
	min-height: 50dvh;
}

.footer {
	display: flex;
	justify-content: space-between;
	gap: 8px;
}

.iframe {
	height: 100%;
	width: 100%;
	aspect-ratio: 16/9;
	border: 0;
	border: 1px solid v-bind('themeVars.borderColor');
	border-radius: 8px;
}
</style>
