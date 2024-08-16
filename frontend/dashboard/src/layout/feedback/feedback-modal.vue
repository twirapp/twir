<script setup lang="ts">
import { NAlert, NForm, NFormItem, NInput } from 'naive-ui'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useFeedbackForm } from './feedback.js'

import type { FormInst, FormRules } from 'naive-ui'

const messageMaxLength = 1000

const { form, error } = useFeedbackForm()
const { t } = useI18n()

const formRef = ref<FormInst | null>(null)
const rules: FormRules = {
	message: {
		required: true,
		message: t('feedback.validation.emptyMessage'),
		trigger: ['input', 'blur'],
	},
}
</script>

<template>
	<NForm ref="formRef" :rules="rules" :model="form">
		<NFormItem :label="t('feedback.messageLabel')" path="message">
			<NInput
				v-model:value="form.message"
				:maxlength="messageMaxLength"
				type="textarea"
				show-count
				:count-graphemes="(v) => v.length"
				:autosize="{
					minRows: 3,
				}"
			/>
		</NFormItem>
	</NForm>

	<NAlert v-if="error" type="error">
		{{ error }}
	</NAlert>
</template>
