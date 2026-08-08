<script setup lang="ts">
import InputWithIcon from '@/components/ui/InputWithIcon/InputWithIcon.vue'
import { ColorPicker } from '@/components/ui/color-picker'
import { Label } from '@/components/ui/label'

interface Props {
	id: string
	label: string
	class?: string
	description?: string
}

withDefaults(defineProps<Props>(), {
	class: 'flex flex-col gap-2',
	description: undefined,
})

const model = defineModel<string>({ required: true })
</script>

<template>
	<div :class="class">
		<Label :for="id">{{ label }}</Label>
		<InputWithIcon :id="id" v-model="model">
			<ColorPicker v-model="model" />
		</InputWithIcon>
		<p v-if="description" class="text-xs text-muted-foreground">{{ description }}</p>
	</div>
</template>
