<script setup lang="ts">
import { NForm, NFormItem, NInput, NAlert } from 'naive-ui';
import type { FormInst, FormRules } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useFeedbackForm } from './feedback.js';

const messageMaxLength = 1000;

const feedbackFormStore = useFeedbackForm();
const { t } = useI18n();

const { form, error } = storeToRefs(feedbackFormStore);

const formRef = ref<FormInst | null>(null);
const rules: FormRules = {
	message: {
		required: true,
		message: t('feedback.validation.emptyMessage'),
		trigger: ['input', 'blur'],
	},
};
</script>

<template>
	<n-form ref="formRef" :rules="rules" :model="form">
		<n-form-item :label="t('feedback.messageLabel')" path="message">
			<n-input
				v-model:value="form.message"
				:maxlength="messageMaxLength"
				type="textarea"
				show-count
				:count-graphemes="(v) => v.length"
				:autosize="{
					minRows: 3,
				}"
			/>
		</n-form-item>
	</n-form>
	<n-alert v-if="error" type="error">
		{{ error }}
	</n-alert>
</template>
