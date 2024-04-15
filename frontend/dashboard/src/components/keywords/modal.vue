<script setup lang="ts">
import {
	type FormInst,
	type FormRules,
	type FormItemRule,
	NForm,
	NFormItem,
	NInput,
	NInputNumber,
	NSpace,
	NCheckbox,
	NDivider,
	NGrid,
	NGridItem,
	NSwitch,
	NCard,
	NText,
	NButton,
	NAlert,
	NA,
} from 'naive-ui';
import { ref, onMounted, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { useKeywordsApi, type Keyword } from '@/api/keywords';
import VariableInput from '@/components/variable-input.vue';

const props = defineProps<{
	keyword?: Keyword | null
}>();

const emits = defineEmits<{
	close: []
}>();

const { t } = useI18n();
const formRef = ref<FormInst | null>(null);
const formValue = ref<Keyword>({
	text: '',
	response: '',
	cooldown: 0,
	enabled: true,
	isReply: true,
	usageCount: 0,
	isRegularExpression: false,
});

onMounted(() => {
	if (!props.keyword) return;
	formValue.value = structuredClone(toRaw(props.keyword));
});

const keywordsApi = useKeywordsApi();
const keywordsUpdate = keywordsApi.useMutationUpdateKeyword();
const keywordsCreate = keywordsApi.useMutationCreateKeyword();

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;
	if (data.id) {
		await keywordsUpdate.executeMutation({
			id: data.id,
			opts: {
				text: data.text,
				response: data.response,
				cooldown: data.cooldown,
				isReply: data.isReply,
				enabled: data.enabled,
				isRegularExpression: data.isRegularExpression,
				usageCount: data.usageCount,
			},
		});
	} else {
		await keywordsCreate.executeMutation({ opts: data });
	}

	emits('close');
}

const rules: FormRules = {
	text: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('keywords.validations.triggerRequired'));
			}

			if (value.length > 500) return new Error(t('keywords.validations.triggerLong'));

			return true;
		},
	},
	response: {
		trigger: ['input', 'blue'],
		validator: (_: FormItemRule, value: string) => {
			if (value?.length > 500) {
				return new Error(t('keywords.validations.responseLong'));
			}
		},
	},
};
</script>

<template>
	<n-form ref="formRef" :model="formValue" :rules="rules">
		<n-space vertical class="w-full">
			<n-space vertical class="gap-0">
				<n-form-item :label="t('keywords.triggerText')" path="text" show-require-mark>
					<n-input v-model:value="formValue.text" />
				</n-form-item>
				<n-checkbox v-model:checked="formValue.isRegularExpression">
					{{ t('keywords.isRegular') }}
				</n-checkbox>
				<n-alert
					v-if="formValue.isRegularExpression"
					type="info"
				>
					<i18n-t
						keypath="keywords.regularDescription"
					>
						<n-a
							href="https://yourbasic.org/golang/regexp-cheat-sheet/#cheat-sheet"
							target="_blank"
						>
							{{ t('keywords.regularDescriptionCheatSheet') }}
						</n-a>
					</i18n-t>
				</n-alert>
			</n-space>

			<n-form-item :label="t('sharedTexts.response')" path="response">
				<variable-input
					v-model="formValue.response"
					:min-rows="1"
					:max-rows="6"
					inputType="textarea"
				/>
			</n-form-item>

			<n-divider>{{ t('keywords.settings') }}</n-divider>

			<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
				<n-grid-item :span="1">
					<n-form-item :label="t('keywords.cooldown')" path="cooldown">
						<n-input-number v-model:value="formValue.cooldown" />
					</n-form-item>
				</n-grid-item>

				<n-grid-item :span="1">
					<n-form-item :label="t('keywords.usages')" path="usages">
						<n-input-number v-model:value="formValue.usageCount" />
					</n-form-item>
				</n-grid-item>

				<n-grid-item :span="1">
					<n-card>
						<div class="flex flex-row justify-between">
							<n-space vertical>
								<n-text>{{ t('sharedTexts.reply.label') }}</n-text>
								<n-text>{{ t('sharedTexts.reply.text') }}</n-text>
							</n-space>
							<n-switch v-model:value="formValue.isReply" />
						</div>
					</n-card>
				</n-grid-item>
			</n-grid>

			<n-button secondary type="success" block @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</n-space>
	</n-form>
</template>
