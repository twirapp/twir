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

import TwitchUserSearch from '@/components/twitchUsers/single.vue';
import VariableInput from '@/components/variable-input.vue';
import { useGreetingsApi, type Greetings } from '@/api/greetings';
import type { GreetingsCreateInput } from '@/gql/graphql';

const props = defineProps<{
	greeting?: Greetings | null
}>();
const emits = defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<Omit<Greetings, 'twitchProfile'>>({
	id: '',
	text: '',
	userId: '',
	enabled: true,
	isReply: true,
});

onMounted(() => {
	if (!props.greeting) return;
	formValue.value = structuredClone(toRaw(props.greeting));
});

const greetingsApi = useGreetingsApi();
const greetingsUpdate = greetingsApi.useMutationUpdateGreetings();
const greetingsCreate = greetingsApi.useMutationCreateGreetings();

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;
	const opts: GreetingsCreateInput = {
		enabled: data.enabled,
		isReply: data.isReply,
		text: data.text,
		userId: data.userId,
	}

	if (data.id) {
		await greetingsUpdate.executeMutation({
			id: data.id,
			opts,
		});
	} else {
		await greetingsCreate.executeMutation({ opts });
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
		<n-space vertical class="w-full">
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
