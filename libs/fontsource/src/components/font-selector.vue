<script setup lang="ts">
import type { SelectOption } from 'naive-ui';
import { NSelect } from 'naive-ui';
import { computed, watch, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { generateFontKey } from '../api.js';
import { useFontSource } from '../composable/use-fontsource.js';
import type { Font } from '../types.js';

const props = defineProps<{
	fontFamily: string
	fontWeight: number
	fontStyle: string
	subsets?: string[]
}>();

const emits = defineEmits<{
	'update-font': [font: Font]
}>();

const { t } = useI18n();
const fontSource = useFontSource();

const availableSubsets = ref<Set<string>>(new Set());
const filteredSubsets = ref<string[]>([]);

const unsubscripteFontList = watch(fontSource.fontList, (v) => {
	if (!v) return;

	for (const font of v) {
		for (const subset of font.subsets) {
			availableSubsets.value.add(subset);
		}
	}
});

if (props.subsets) {
	unsubscripteFontList();
	for (const subset of props.subsets) {
		availableSubsets.value.add(subset);
	}
}

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

const fontOptions = computed((): SelectOption[] => {
	return fontSource.fontList.value
		.filter((font) => {
			if (!filteredSubsets.value.length) return true;
			return filteredSubsets.value.every((subset) => font.subsets.includes(subset));
		})
		.map((font) => ({
			label: font.family,
			value: font.id,
		}));
});

const availableSubsetsOptions = computed(() => {
	return [...availableSubsets.value.values()]
		.map(subset => ({ label: subset, value: subset }));
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
		:options="fontOptions"
		:loading="fontSource.loading.value"
		:disabled="fontSource.loading.value"
		check-strategy="child"
	>
		<template v-if="!props.subsets" #action>
			{{ t('overlays.chat.availabeFonts') }}: {{ fontOptions.length }}
			<n-select
				v-model:value="filteredSubsets"
				clearable
				multiple
				size="tiny"
				:options="availableSubsetsOptions"
				:placeholder="t('overlays.chat.selectSubsetPlaceholder')"
			/>
		</template>
	</n-select>
</template>
