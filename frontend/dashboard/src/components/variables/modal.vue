<script setup lang="ts">
import {
	type FormInst,
	type FormRules,
	type FormItemRule,
	type SelectOption,
	NForm,
	NSpace,
	NInput,
	NFormItem,
	NSelect,
	NButton,
} from 'naive-ui';
import { ref, onMounted, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { type EditableCustomVariable, useVariablesApi } from '@/api/variables';
import { VariableType } from '@/gql/graphql';

const props = defineProps<{
	variable?: EditableCustomVariable | null
}>();
const emits = defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableCustomVariable>({
	type: VariableType.Text,
	name: '',
	evalValue: `// semicolons (;) matters, do not forget put them on end of statements.
const request = await fetch('https://jsonplaceholder.typicode.com/todos/1');
const response = await request.json();
// you should return value from your script
return response.title;`,
	description: '',
	response: '',
	id: '',
});

onMounted(() => {
	if (!props.variable) return;
	formValue.value = structuredClone(toRaw(props.variable));
});

const variablesManager = useVariablesApi();
const variablesUpdater = variablesManager.useMutationUpdateVariable();
const variablesCreator = variablesManager.useMutationCreateVariable();

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const opts = {
		type: formValue.value.type,
		name: formValue.value.name,
		evalValue: formValue.value.evalValue,
		description: formValue.value.description,
		response: formValue.value.response,
	};

	if (props.variable?.id) {
		await variablesUpdater.executeMutation({
			id: props.variable.id,
			opts,
		});
	} else {
		await variablesCreator.executeMutation({ opts });
	}

	emits('close');
}

const { t } = useI18n();

const rules: FormRules = {
	name: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('variables.validations.nameRequired'));
			}
			if (value.length > 20) {
				return new Error(t('variables.validations.nameLong'));
			}
			return true;
		},
	},
};
const selectOptions: Array<SelectOption> = Object.entries(VariableType).map(([name, value]) => ({
	label: name,
	value: value as string,
}));
</script>

<template>
	<n-form
		ref="formRef"
		:model="formValue"
		:rules="rules"
		:style="{
			width: formValue.type === VariableType.Script ? '900px' : '400px',
		}"
	>
		<n-space vertical class="w-full">
			<n-form-item :label="t('sharedTexts.name')" path="name" show-require-mark>
				<n-input v-model:value="formValue.name" :placeholder="t('sharedTexts.name')" />
			</n-form-item>

			<n-form-item :label="t('variables.type')" path="type" show-require-mark>
				<n-select
					v-model:value="formValue.type"
					:options="selectOptions"
				/>
			</n-form-item>

			<n-form-item :label="t('sharedTexts.response')" path="response">
				<vue-monaco-editor
					v-if="formValue.type === VariableType.Script"
					v-model:value="formValue.evalValue"
					theme="vs-dark"
					height="500px"
					language="javascript"
				/>
				<n-input
					v-else
					v-model:value="formValue.response"
				/>
			</n-form-item>
		</n-space>
		<n-button secondary type="success" block class="mt-2.5" @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-form>
</template>
