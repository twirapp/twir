<script setup lang="ts">
import type { SelectOption } from 'naive-ui';
import { NSelect } from 'naive-ui';
import { ref, computed, watch, h } from 'vue';

import { useFontSource } from '../composable/use-fontsource.js';

const props = defineProps<{
	initialFontFamily: string
}>();

const emits = defineEmits<{
	'update-font-family': [fontFamily: string]
}>();

const fontSource = useFontSource();
const selectedFont = ref(props.initialFontFamily);

watch(() => selectedFont.value, async (font) => {
	if (!font || fontSource.loading.value) return;
	await fontSource.loadFont(font);
	emits('update-font-family', font);
});

const options = computed((): SelectOption[] => {
	return fontSource.fonts.value.map((font) => ({
		label: font.family,
		value: font.id,
	}));
});

function renderLabel(option: SelectOption) {
	if (!fontSource.loading.value) {
		fontSource.loadFont(option.value as string);
	}

	return h(
		'div',
		{ style: { 'font-family': `"${option.value}"` } },
		{ default: () => option.label },
	);
}
</script>

<template>
	<n-select
		v-model:value="selectedFont"
		:default-value="props.initialFontFamily"
		:render-label="renderLabel"
		filterable
		:options="options"
		:loading="fontSource.loading.value"
		:disabled="fontSource.loading.value"
		check-strategy="child"
	/>
</template>
