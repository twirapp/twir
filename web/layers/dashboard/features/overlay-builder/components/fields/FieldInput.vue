<script setup lang="ts">
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

interface Props {
	id: string
	label: string
	type?: string
	placeholder?: string
	list?: string
	class?: string
	inputClass?: string
	description?: string
}

withDefaults(defineProps<Props>(), {
	type: 'text',
	placeholder: undefined,
	class: 'flex flex-col gap-2',
	inputClass: undefined,
	list: undefined,
	description: undefined,
})

const model = defineModel<string | number>({ required: true })
const emit = defineEmits<{
	blur: [event: FocusEvent]
}>()
</script>

<template>
	<div :class="class">
		<Label :for="id">{{ label }}</Label>
		<Input
			:id="id"
			v-model="model"
			:type="type"
			:placeholder="placeholder"
			:list="list"
			:class="inputClass"
			@keydown.stop
			@blur="emit('blur', $event)"
		/>
		<p v-if="description" class="text-xs text-muted-foreground">{{ description }}</p>
	</div>
</template>
