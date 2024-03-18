<script setup lang='ts'>
import { IconVariable } from '@tabler/icons-vue';
import { NText, NInput, NInputGroup, NButton, NPopselect } from 'naive-ui';
import { type SelectMixedOption } from 'naive-ui/es/select/src/interface';
import { computed, VNodeChild, h, FunctionalComponent, ref } from 'vue';

import { useAllVariables } from '@/api/index.js';

// eslint-disable-next-line no-undef
const text = defineModel({ default: '' });

withDefaults(defineProps<{
	inputType?: 'text' | 'textarea',
	minRows?: number,
	maxRows?: number,
}>(), {
	inputType: 'text',
});

defineSlots<{
	underSelect: FunctionalComponent
}>();

const allVariables = useAllVariables();
const search = ref('');

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
			})).filter(v => v.value.includes(search.value)),
		},
		{
			type: 'group',
			label: 'Built in',
			key: 'Built in',
			children: builtIn.map(v => ({
				label: v.name,
				value: v.example || v.name,
				description: v.description,
			})).filter(v => v.value.includes(search.value)),
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
			trigger="click"
			:value="null"
			:render-label="renderVariableSelectLabel as any"
			:on-update-value="(_, option) => appendOptionToText(option)"
		>
			<n-button style="height:auto">
				<IconVariable />
			</n-button>
			<template #action>
				<n-input v-model:value="search" placeholder="Search..." size="small" />
			</template>
		</n-popselect>
	</n-input-group>
</template>

<style scoped lang='postcss'>

</style>
