<script setup lang="ts">
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

interface Props {
	id: string
	label: string
	min?: number
	max?: number
	step?: number
	placeholder?: string
	class?: string
	inputClass?: string
}

withDefaults(defineProps<Props>(), {
	min: undefined,
	max: undefined,
	step: undefined,
	placeholder: undefined,
	class: 'flex flex-col gap-2',
	inputClass: undefined,
})

const model = defineModel<number>({ required: true })
</script>

<template>
	<div :class="class">
		<Label :for="id">{{ label }}</Label>
		<Input
			:id="id"
			v-model.number="model"
			type="number"
			:min="min"
			:max="max"
			:step="step"
			:placeholder="placeholder"
			:class="inputClass"
			@keydown.stop
		/>
	</div>
</template>
