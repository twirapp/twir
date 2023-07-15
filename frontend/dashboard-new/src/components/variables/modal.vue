<script setup lang='ts'>
import { useMonaco } from '@guolao/vue-monaco-editor';
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
	NInputNumber,
} from 'naive-ui';
import { ref, onMounted, toRaw, onUnmounted, watch } from 'vue';

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
	response: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('Response is required');
			}
			if (value.length > 1000) {
				return new Error('Response is too long');
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

const monacoContainerRef = ref();
const { monacoRef, unload: unloadMonaco } = useMonaco();
onUnmounted(() => !monacoRef.value && unloadMonaco());
</script>

<template>
  <n-form
    ref="formRef"
    :model="formValue"
    :rules="rules"
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
        <div v-if="formValue.type === VariableType.SCRIPT" ref="monacoContainerRef"></div>
        <n-input
          v-else
          v-model:value="formValue.evalValue"
        />
      </n-form-item>
    </n-space>
  </n-form>
</template>

<style scoped lang='postcss'>

</style>
