<script setup lang='ts'>
import type { GetResponse as OBSSettings } from '@twir/grpc/generated/api/api/modules_obs_websocket';
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
	useMessage,
} from 'naive-ui';
import { ref, watch, toRaw } from 'vue';

import { useObsOverlayManager } from '@/api/index.js';

const obsSettingsManager = useObsOverlayManager();
const obsSettings = obsSettingsManager.getSettings();
const obsSettingsUpdater = obsSettingsManager.updateSettings();

const formRef = ref<FormInst | null>(null);
const formValue = ref<OBSSettings>({
	audioSources: [],
	scenes: [],
	sources: [],
	serverAddress: '',
	serverPassword: '',
	serverPort: 4455,
});
const rules: FormRules = {
	address: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('Address is required');
			}

			return true;
		},
	},
	port: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: number) => {
			if (!value) {
				return new Error('Port is required');
			}

			return true;
		},
	},
	password: {
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
	message.success('Settings updated');
}
</script>

<template>
  <n-alert type="info">
    This overlay used for connect TwirApp with your obs. It gives opportunity to bot manage your sources, scenes, audio sources on events.
  </n-alert>

  <n-form
    ref="formRef"
    :model="formValue"
    :rules="rules"
    style="margin-top:15px"
  >
    <n-form-item
      label="Address."
      required
      path="address"
    >
      <n-input placeholder="Usually it's localhost" />
    </n-form-item>

    <n-form-item label="Port" required path="port">
      <n-input-number :min="1" :max="66000" placeholder="Socket port" />
    </n-form-item>

    <n-form-item label="Password" required path="password">
      <n-input
        type="password"
        show-password-on="click"
        placeholder="Socket password"
        :maxlength="8"
      />
    </n-form-item>

    <n-button block secondary type="success" @click="save">
      Save
    </n-button>
  </n-form>
</template>
