<script setup lang='ts'>
import {
	type FormInst,
	type FormRules,
	type FormItemRule,
	NForm,
	NSpace,
	NFormItem,
	NButton,
	NSwitch,
} from 'naive-ui';
import { ref, onMounted, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { useGreetingsManager } from '@/api/index.js';
import type { EditableGreeting } from '@/components/greetings/types.js';
import TwitchUserSearch from '@/components/twitchUsers/single.vue';
import VariableInput from '@/components/variable-input.vue';

const props = defineProps<{
	greeting?: EditableGreeting | null
}>();
const emits = defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableGreeting>({
	enabled: true,
	isReply: true,
	text: '',
	userId: '',
	id: '',
});

onMounted(() => {
	if (!props.greeting) return;
	formValue.value = structuredClone(toRaw(props.greeting));
});

const greetingsManager = useGreetingsManager();
const greetingsUpdater = greetingsManager.update;
const greetingsCreator = greetingsManager.create;

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;

	if (data.id) {
		await greetingsUpdater.mutateAsync({
			id: data.id,
			greeting: data,
		});
	} else {
		await greetingsCreator.mutateAsync(data);
	}

	emits('close');
}

const { t } = useI18n();

const rules: FormRules = {
	userId: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('greetings.validations.userName'));
			}

			return true;
		},
	},
	text: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('greetings.validations.textRequired'));
			}

			return true;
		},
	},
};
</script>

<template>
	<n-form
		ref="formRef"
		:model="formValue"
		:rules="rules"
	>
		<n-space vertical style="width: 100%">
			<n-form-item :label="t('sharedTexts.userName')" path="userId" show-require-mark>
				<twitch-user-search
					v-model="formValue.userId"
					:initial-user-id="props.greeting?.userId"
				/>
			</n-form-item>
			<n-form-item :label="t('sharedTexts.response')" path="text" show-require-mark>
				<variable-input v-model="formValue.text" input-type="textarea" />
			</n-form-item>

			<n-form-item :label="t('sharedTexts.reply.text')" path="text">
				<n-switch v-model:value="formValue.isReply" />
			</n-form-item>
		</n-space>

		<n-button secondary type="success" block @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-form>
</template>
