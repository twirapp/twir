<script setup lang="ts">
import { NSelect } from 'naive-ui'
import { computed, h, ref, watch } from 'vue'
import { type VNodeChild, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

import { generateFontKey } from '../api.js'
import { useFontSource } from '../composable/use-fontsource.js'

import type { Font } from '../types.js'
import type { SelectMixedOption } from 'naive-ui/es/select/src/interface.js'

const props = defineProps<{
	fontFamily: string
	fontWeight: number
	fontStyle: string
	subsets?: string[]
}>()

const { t } = useI18n()
const fontSource = useFontSource()
const font = defineModel<Font | null>('font')

const selectedFont = ref<string>('')
const availableSubsets = ref<Set<string>>(new Set())
const filteredSubsets = ref<string[]>(props.subsets ?? [])

watch(() => fontSource.fontList.value, (fonts) => {
	if (!fonts) return

	for (const font of fonts) {
		for (const subset of font.subsets) {
			availableSubsets.value.add(subset)
		}
	}
})

watch(() => selectedFont.value, (selectedFont) => {
	font.value = fontSource.getFont(selectedFont)
})

interface FontOption {
	label: string
	value: string
	fontWeight: number
	fontStyle: string
}

const fontOptions = computed((): Array<SelectMixedOption & FontOption> => {
	return fontSource.fontList.value
		.filter((font) => {
			if (!filteredSubsets.value.length) return true
			return filteredSubsets.value.every((subset) => font.subsets.includes(subset))
		})
		.map((font) => ({
			label: font.family,
			value: font.id,
			fontWeight: font.weights.includes(400) ? 400 : font.weights[0],
			fontStyle: font.styles.includes('normal') ? 'normal' : font.styles[0],
		}))
})

const availableSubsetsOptions = computed(() => {
	return [...availableSubsets.value.values()]
		.map(subset => ({ label: subset, value: subset }))
})

function renderLabel(option: FontOption): VNodeChild {
	if (!fontSource.loading.value) {
		fontSource.loadFont(option.value, option.fontWeight, option.fontStyle)
	}

	const fontFamily = generateFontKey(option.value!, option.fontWeight, option.fontStyle)
	return h(
		'div',
		{ style: { 'font-family': `"${fontFamily}"` } },
		{ default: () => option.label },
	)
}

onMounted(async () => {
	const loadedFont = await fontSource
		.loadFont(props.fontFamily, props.fontWeight, props.fontStyle)

	if (loadedFont) {
		font.value = loadedFont
		selectedFont.value = loadedFont.id
	}
})
</script>

<template>
	<NSelect
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
			<NSelect
				v-model:value="filteredSubsets"
				clearable
				multiple
				size="tiny"
				:options="availableSubsetsOptions"
				:placeholder="t('overlays.chat.selectSubsetPlaceholder')"
			/>
		</template>
	</NSelect>
</template>
