<script setup lang='ts'>
import { IconVariable } from '@tabler/icons-vue';
import { NText, NInput, NInputGroup, NButton, NPopselect } from 'naive-ui';
import { type SelectMixedOption } from 'naive-ui/es/select/src/interface';
import { computed, VNodeChild, h, defineModel, FunctionalComponent, defineSlots } from 'vue';

import { useAllVariables } from '@/api/index.js';

const text = defineModel({ default: '' });

const props = withDefaults(defineProps<{
	inputType: 'text' | 'textarea',
	minRows: number,
	maxRows: number,
}>(), {
	inputType: 'text',
});

defineSlots<{
	underSelect: FunctionalComponent
}>();

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

function appendOptionToText(option: SelectMixedOption) {
	const newText = ` $(${option.value})`;
	text.value += newText;
}
</script>

<template>
  <n-input-group>
    <n-input
      v-model:value="text"
      style="width: 100%"
      :type="inputType"
      :autosize="inputType === 'text' ? {} : { minRows, maxRows }"
      placeholder="Response text"
    />
    <n-popselect
      :options="selectVariables"
      :loading="allVariables.isLoading.value"
      scrollable
      :value="null"
      :render-label="renderVariableSelectLabel as any"
      :on-update-value="(_, option) => appendOptionToText(option)"
    >
      <n-button style="height:auto">
        <IconVariable />
      </n-button>
    </n-popselect>
  </n-input-group>
</template>

<style scoped lang='postcss'>

</style>
