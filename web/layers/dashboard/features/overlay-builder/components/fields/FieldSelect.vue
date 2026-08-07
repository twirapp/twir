<script setup lang="ts">
import type { AcceptableValue } from 'reka-ui'

import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

interface Option {
	readonly value: string
	readonly label: string
}

interface Props {
	id: string
	label: string
	options: readonly Option[]
	placeholder?: string
	class?: string
}

withDefaults(defineProps<Props>(), {
	placeholder: undefined,
	class: 'flex flex-col gap-2',
})

const model = defineModel<string>({ required: true })

function updateModel(value: AcceptableValue) {
	if (typeof value === 'string') model.value = value
}
</script>

<template>
	<div :class="class">
		<Label :for="id">{{ label }}</Label>
		<Select :model-value="model" @update:model-value="updateModel">
			<SelectTrigger :id="id" class="w-full">
				<SelectValue :placeholder="placeholder" />
			</SelectTrigger>
			<SelectContent>
				<SelectItem v-for="option in options" :key="option.value" :value="option.value">
					{{ option.label }}
				</SelectItem>
			</SelectContent>
		</Select>
	</div>
</template>
