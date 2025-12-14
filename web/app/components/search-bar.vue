<script setup lang="ts">
defineProps<{ placeholder?: string }>()

const model = defineModel<string>({ default: '' })
</script>

<template>
	<div class="relative w-full items-center">
		<UiInput
			v-model="model"
			type="text"
			:placeholder="placeholder ?? 'Search...'"
			class="h-9 pl-10"
		/>
		<span class="absolute start-2 inset-y-0 flex items-center justify-center px-2">
			<Icon name="lucide:search" class="w-4 h-4 text-muted-foreground" />
		</span>
	</div>
</template>
