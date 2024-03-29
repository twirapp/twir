<script setup lang='ts'>
import { IconSearch } from '@tabler/icons-vue';
import { NInput, NDropdown, NTooltip } from 'naive-ui';
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

const header = {
	key: 'header',
	type: 'render',
	render: () => {
		return h(
			NInput,
			{
				value: search.value,
				onInput: (v: string) => search.value = v, placeholder: 'Search for variable...',
				size: 'small',
			},
			{
				suffix: () => h(IconSearch),
			},
		);
	},
};

const selectVariables = computed<SelectMixedOption[]>(() => {
	const variables = allVariables.data?.value;
	if (!variables) return [];

	const builtIn = variables.filter((variable) => variable.isBuiltIn);
	const custom = variables.filter((variable) => !variable.isBuiltIn);

	return [
		header,
		{
			type: 'group',
			label: 'Custom',
			key: 'Custom',
			children: custom.map(v => ({
				label: v.name,
				value: v.example || v.name,
				description: v.description,
			})).filter(v => v.value.includes(search.value) || v.description?.includes(search.value)),
		},
		{
			type: 'group',
			label: 'Built in',
			key: 'Built in',
			children: builtIn.map(v => ({
				label: v.name,
				value: v.example || v.name,
				description: v.description,
			})).filter(v => v.value.includes(search.value) || v.description?.includes(search.value)),
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
		NTooltip,
		{ placement: 'left' },
		{
			default: () => option.description,
			trigger: () => h('span', { style: { width: '100%', display: 'flex' } }, `$(${option.label})`),
		},
	);
}

function appendOptionToText(_: unknown, option: SelectMixedOption) {
	const newText = ` $(${option.value})`;
	text.value += newText;
}

const showDropdown = ref(false);
</script>

<template>
	<n-dropdown
		trigger="hover"
		:options="selectVariables"
		style="max-width: 100%; max-height: 300px;"
		:render-label="renderVariableSelectLabel as any"
		scrollable
		@select="appendOptionToText"
	>
		<n-input
			v-model:value="text"
			style="width: 100%"
			:type="inputType"
			:autosize="inputType === 'text' ? {} : { minRows, maxRows }"
			placeholder="Response text"
			@click="showDropdown = true"
		/>
	</n-dropdown>
</template>

<style scoped>
:deep(.n-base-select-menu) {
	max-width: 400px;
}
</style>
