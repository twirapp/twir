<script setup lang="ts">
import type {
	GetResponse as OBSSettings,
} from '@twir/api/messages/modules_obs_websocket/modules_obs_websocket';
import {
	type FormInst,
	type FormRules,
	type FormItemRule,
	NForm,
	NFormItem,
	NInput,
	NInputNumber,
	NButton,
	NAlert,
	NSpace,
	useMessage,
} from 'naive-ui';
import { ref, toRaw, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';

import { useObsOverlayManager, useProfile } from '@/api/index.js';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink';

const obsSettingsManager = useObsOverlayManager();
const { refetch, data: settings } = obsSettingsManager.getSettings();
const obsSettingsUpdater = obsSettingsManager.updateSettings();

onMounted(async () => {
	const settings = await refetch();
	if (!settings.data) return;
	formValue.value = toRaw(settings.data);
});

const { t } = useI18n();

const formRef = ref<FormInst | null>(null);
const formValue = ref<Omit<OBSSettings, 'isConnected'>>({
	audioSources: [],
	scenes: [],
	sources: [],
	serverAddress: 'localhost',
	serverPassword: '',
	serverPort: 4455,
});
const rules: FormRules = {
	serverAddress: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('Address is required');
			}

			return true;
		},
	},
	serverPort: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: number) => {
			if (!value) {
				return new Error('Port is required');
			}

			return true;
		},
	},
	serverPassword: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('Password is required');
			}

			return true;
		},
	},
};

const message = useMessage();

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	await obsSettingsUpdater.mutateAsync(formValue.value);
	message.success('Settings updated, now you can paste overlay link into obs', {
		duration: 2500,
	});
}

const { copyOverlayLink } = useCopyOverlayLink('obs');
const { data: profile } = useProfile();
</script>

<template>
	<n-alert type="info">
		{{ t('overlays.obs.description') }}
	</n-alert>

	<n-form
		ref="formRef"
		:model="formValue"
		:rules="rules"
		style="margin-top:15px"
	>
		<n-form-item
			:label="t('overlays.obs.address')"
			required
			path="serverAddress"
		>
			<n-input v-model:value="formValue.serverAddress" placeholder="Usually it's localhost" />
		</n-form-item>

		<n-form-item :label="t('overlays.obs.port')" required path="serverPort">
			<n-input-number
				v-model:value="formValue.serverPort"
				:min="1"
				:max="66000"
				placeholder="Socket port"
			/>
		</n-form-item>

		<n-form-item label="Password" required path="serverPassword">
			<n-input
				v-model:value="formValue.serverPassword"
				type="password"
				show-password-on="click"
				placeholder="Socket password"
			/>
		</n-form-item>

		<n-alert :type="settings?.isConnected ? 'success' : 'error'" :bordered="false">
			{{
				settings?.isConnected ? t('overlays.obs.connected') : t('overlays.obs.notConnected')
			}}
		</n-alert>

		<n-space vertical class="mt-2.5">
			<n-button :disabled="profile?.id !== profile?.selectedDashboardId" block secondary type="info" @click="copyOverlayLink()">
				{{ t('overlays.copyOverlayLink') }}
			</n-button>

			<n-button block secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</n-space>
	</n-form>
</template>
