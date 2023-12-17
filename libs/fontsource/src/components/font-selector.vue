<script setup lang="ts">
import type { SelectOption } from 'naive-ui';
import { NSelect } from 'naive-ui';
import { ref, computed, watch, h } from 'vue';

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
const selectedFont = ref(props.fontFamily);

watch(() => props.fontFamily, () => {
	selectedFont.value = props.fontFamily;
});

watch([fontSource.fonts.value, selectedFont], async () => {
	const font = fontSource.getFont(selectedFont.value);
	if (!font) return;
	emits('update-font', font);
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
		:default-value="props.fontFamily"
		:render-label="renderLabel"
		filterable
		:options="options"
		:loading="fontSource.loading.value"
		:disabled="fontSource.loading.value"
		check-strategy="child"
	/>
</template>
