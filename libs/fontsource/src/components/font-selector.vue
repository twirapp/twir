<script setup lang="ts">
import type { SelectOption } from 'naive-ui';
import { NSelect } from 'naive-ui';
import { computed, watch, h } from 'vue';

import { generateFontKey } from '../api.js';
import { useFontSource } from '../composable/use-fontsource.js';
import type { Font } from '../types.js';

const props = defineProps<{
	fontFamily: string
	fontWeight: number
	fontStyle: string
}>();

const emits = defineEmits<{
	'update-font': [font: Font]
}>();

const fontSource = useFontSource();

// eslint-disable-next-line no-undef
const selectedFont = defineModel<string>('selectedFont', { default: '' });

watch(selectedFont, (v) => {
	if (!v) return;
	const font = fontSource.getFont(v);
	if (!font) return;
	emits('update-font', font);
});

watch(() => props.fontFamily, (v) => {
	selectedFont.value = v;
});

const options = computed((): SelectOption[] => {
	return fontSource.fontList.value.map((font) => ({
		label: font.family,
		value: font.id,
	}));
});

function renderLabel(option: SelectOption) {
	if (!fontSource.loading.value) {
		fontSource.loadFont(option.value as string, props.fontWeight, props.fontStyle);
	}

	const fontFamily = generateFontKey(option.value as string, props.fontWeight, props.fontStyle);
	return h(
		'div',
		{ style: { 'font-family': `"${fontFamily}"` } },
		{ default: () => option.label },
	);
}
</script>

<template>
	<n-select
		v-model:value="selectedFont"
		:render-label="renderLabel"
		filterable
		:options="options"
		:loading="fontSource.loading.value"
		:disabled="fontSource.loading.value"
		check-strategy="child"
	/>
</template>
