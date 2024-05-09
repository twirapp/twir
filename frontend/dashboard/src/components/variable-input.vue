<script setup lang='ts'>
import { NMention, NText } from 'naive-ui'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import type { MentionOption } from 'naive-ui'
import type { FunctionalComponent, VNodeChild } from 'vue'

import { useVariablesApi } from '@/api/variables.js'

withDefaults(defineProps<{
	inputType?: 'text' | 'textarea'
	minRows?: number
	maxRows?: number
}>(), {
	inputType: 'text',
})

defineSlots<{
	underSelect: FunctionalComponent
}>()

const text = defineModel<string>({ default: '' })
const { t } = useI18n()

const { allVariables } = useVariablesApi()

const selectVariables = computed<MentionOption[]>(() => {
	return allVariables.value.map((variable) => ({
		label: `(${variable.example})`,
		value: `(${variable.example})`,
		description: variable.description,
	}))
})

function renderVariableSelectLabel(option: {
	type: string
	label: string
	description: string
}): VNodeChild {
	if (!option.description) return `$${option.label}`
	const variable = `$${option.label}`
	const description = h(NText, { depth: 3 }, option.description)
	return h('span', {}, [variable, ' ', description])
}
</script>

<template>
	<NMention
		v-model:value="text"
		:render-label="renderVariableSelectLabel"
		placeholder="Response"
		prefix="$"
		class="w-full"
		:type="inputType"
		:options="selectVariables"
		:autosize="inputType === 'text' ? {} : { minRows, maxRows }"
	>
		<template #empty>
			{{ t('sharedTexts.placeCursorMessage') }}
		</template>
	</NMention>
</template>
