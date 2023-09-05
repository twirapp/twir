<script setup lang='ts'>
import type {
	GetResponse as OBSSettings,
} from '@twir/grpc/generated/api/api/modules_obs_websocket';
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
import { ref, watch, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { useObsOverlayManager } from '@/api/index.js';

const obsSettingsManager = useObsOverlayManager();
const obsSettings = obsSettingsManager.getSettings();
const obsSettingsUpdater = obsSettingsManager.updateSettings();

const { t } = useI18n();

const formRef = ref<FormInst | null>(null);
const formValue = ref<Omit<OBSSettings, 'isConnected'>>({
	audioSources: [],
	scenes: [],
	sources: [],
	serverAddress: '',
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

watch(obsSettings.data, (v) => {
	if (!v) return;
	formValue.value = toRaw(v);
}, { immediate: true });

const message = useMessage();

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	await obsSettingsUpdater.mutateAsync(formValue.value);
	message.success('Settings updated, now you can paste overlay link into obs', {
		duration: 5000,
	});
}

async function checkConnection() {
	await obsSettings.refetch();
}
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

		<n-alert :type="obsSettings.data.value?.isConnected ? 'success' : 'error'" :bordered="false">
			{{
				obsSettings.data.value?.isConnected ? t('overlays.obs.connected') : t('overlays.obs.notConnected')
			}}
		</n-alert>

		<n-space vertical style="margin-top: 10px">
			<n-button block secondary type="info" @click="checkConnection">
				{{ t('overlays.obs.checkConnection') }}
			</n-button>

			<n-button block secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</n-space>
	</n-form>
</template>
