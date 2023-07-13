<script setup lang='ts'>
import { IconVariable } from '@tabler/icons-vue';
import type { Command, Command_Response } from '@twir/grpc/generated/api/api/commands';
import {
	NForm,
	NFormItem,
	FormInst,
	FormRules,
	FormItemRule,
	NInput,
	NDynamicTags,
	NGrid,
	NGridItem,
	NDivider,
	NButton,
	NSpace,
	NInputGroup,
	NInputGroupLabel,
	NDynamicInput,
	NSelect,
	NBadge,
	NText,
} from 'naive-ui';
import { type SelectMixedOption } from 'naive-ui/es/select/src/interface';
import { ref, computed, VNodeChild, h, reactive } from 'vue';

import { useAllVariables } from '@/api/index.js';

defineProps<{
	command: Command | null
}>();

type FormCommand = Omit<Command, 'responses'> & {
	responses: Omit<Command_Response, 'id' | 'commandId' | 'order'>[]
};

const formRef = ref<FormInst | null>(null);
const formValue = reactive<FormCommand>({
	name:'',
	aliases: [],
	responses: [],
});
const nameValidator = (rule: FormItemRule, value: string) => {
	if (!value) {
		return new Error('Please input a name');
	}
	if (value.startsWith('!')) {
		return new Error('Name cannot start with !');
	}
	return true;
};
const rules: FormRules = {
	name: [{
		trigger: ['input', 'blur'],
		validator: nameValidator,
	}],
	aliases: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string[]) => {
			value.forEach((alias) => nameValidator(rule, alias));

			return true;
		},
	},
};

const allVariables = useAllVariables();
const selectVariables = computed<SelectMixedOption[]>(() => {
	const variables = allVariables.data?.value;
	if (!variables) return [];

	const builtIn = variables.filter((variable) => variable.isBuiltIn);
	const custom = variables.filter((variable) => !variable.isBuiltIn);

	return [
		{
			type: 'group',
			label: 'Custom',
			key: 'Custom',
			children: custom.map(v => ({
				label: v.name,
				value: v.example || v.name,
				description: v.description,
			})),
		},
		{
			type: 'group',
			label: 'Built in',
			key: 'Built in',
			children: builtIn.map(v => ({
				label: v.name,
				value: v.example || v.name,
				description: v.description,
			})),
		},
	];
});

function renderVariableSelectLabel(option: {
	type: string,
	label: string,
	description: string
}): VNodeChild {
	if (option.type === 'group') return option.label;
	if (!option.description) return option.label;

	return h(
		'div',
		{
			style: {
				display: 'flex',
				alignItems: 'center',
			},
		},
		[
			h(
				'div',
				{
					style: {
						padding: '4px 0',
					},
				},
				[
					h('div', null, `$(${option.label})`),
					h(
						NText,
						{ depth: 3, tag: 'div' },
						{
							default: () => option.description,
						},
					),
				],
			),
		],
	);
}

function appendOptionToResponse(responseIndex: number, option: SelectMixedOption) {
	const response = formValue.responses[responseIndex];
	if (!response) return;

	response.text += ` $(${option.value})`;
}
</script>

<template>
  <n-form ref="formRef" :model="formValue" :rules="rules">
    <n-grid :cols="12" :x-gap="10">
      <n-grid-item :span="6">
        <n-form-item label="Name" path="name">
          <n-input-group>
            <n-input-group-label>!</n-input-group-label>
            <n-input v-model:value="formValue.name" placeholder="Input Name" />
          </n-input-group>
        </n-form-item>
      </n-grid-item>
      <n-grid-item :span="6">
        <n-form-item label="Aliases" path="aliases">
          <n-dynamic-tags v-model:value="formValue.aliases" />
        </n-form-item>
      </n-grid-item>
    </n-grid>

    <n-divider>
      Responses
    </n-divider>

    <n-dynamic-input
      v-model:value="formValue.responses"
      :on-create="() => ({ text: '' })"
      placeholder="Response"
      show-sort-button
    >
      <template #default="{ value, index: responseIndex }">
        <n-input-group>
          <n-input
            v-model:value="value.text"
            type="text"
          />
          <n-select
            :style="{ width: '33%' }"
            :options="selectVariables"
            :loading="allVariables.isLoading.value"
            placeholder="Search variable..."
            :filterable="true"
            :value="null"
            :clear-filter-after-select="true"
            :consistent-menu-width="false"
            :render-label="renderVariableSelectLabel as any"
            :on-update-value="(_, option) => appendOptionToResponse(responseIndex, option)"
          >
            <template #arrow>
              <IconVariable />
            </template>
          </n-select>
        </n-input-group>
      </template>
    </n-dynamic-input>
    {{ JSON.stringify(formValue) }}
  </n-form>
</template>
