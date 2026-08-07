<script setup lang="ts">
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'

interface Props {
	id: string
	label: string
	description?: string
}

defineProps<Props>()

const model = defineModel<boolean>({ required: true })
</script>

<template>
	<div class="flex items-center justify-between">
		<div class="space-y-0.5">
			<Label :for="id">{{ label }}</Label>
			<p v-if="description" class="text-xs text-muted-foreground">{{ description }}</p>
		</div>
		<Switch :id="id" v-model="model" />
	</div>
</template>
