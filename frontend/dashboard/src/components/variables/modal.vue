<script setup lang='ts'>
import { VariableType } from '@twir/grpc/generated/api/api/variables';
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

import { useVariablesManager } from '@/api/index.js';
import { EditableVariable } from '@/components/variables/types.js';

const props = defineProps<{
	variable?: EditableVariable | null
}>();
const emits = defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableVariable>({
	type: VariableType.TEXT,
	name: '',
	evalValue: '',
	description: '',
	response: '',
	id: '',
});

onMounted(() => {
	if (!props.variable) return;
	formValue.value = structuredClone(toRaw(props.variable));
});

const variablesManager = useVariablesManager();
const variablesUpdater = variablesManager.update;
const variablesCreator = variablesManager.create;

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = {
		...formValue.value,
	};

	if (data.id) {
		await variablesUpdater.mutateAsync({
			id: data.id,
			variable: data,
		});
	} else {
		await variablesCreator.mutateAsync(data);
	}

	emits('close');
}

const rules: FormRules = {
	name: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('Name is required');
			}
			if (value.length > 20) {
				return new Error('Name is too long');
			}
			return true;
		},
	},
};
const selectOptions: Array<SelectOption> = [
	{
		label: 'Text',
		value: VariableType.TEXT,
	},
	{
		label: 'Script',
		value: VariableType.SCRIPT,
	},
	{
		label: 'Number',
		value: VariableType.NUMBER,
	},
];
</script>

<template>
  <n-form
    ref="formRef"
    :model="formValue"
    :rules="rules"
    :style="{
      width: formValue.type === VariableType.SCRIPT ? '900px' : '400px',
    }"
  >
    <n-space vertical style="width: 100%">
      <n-form-item label="Name" path="name">
        <n-input v-model:value="formValue.name" />
      </n-form-item>

      <n-form-item label="Type" path="type">
        <n-select
          v-model:value="formValue.type"
          :options="selectOptions"
        />
      </n-form-item>

      <n-form-item label="Eval value" path="response">
        <vue-monaco-editor
          v-if="formValue.type === VariableType.SCRIPT"
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
    <n-button secondary type="success" block style="margin-top: 10px" @click="save">
      Save
    </n-button>
  </n-form>
</template>

<style scoped lang='postcss'>

</style>
