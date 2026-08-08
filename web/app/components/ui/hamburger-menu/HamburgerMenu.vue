<script setup lang="ts">
defineProps<{
	isOpen: boolean
}>();
</script>

<template>
	<div class="flex cursor-pointer h-9 w-9 items-center justify-center">
		<Icon
			v-if="!isOpen"
			name="lucide:menu"
		/>
		<Icon
			v-else
			name="lucide:x"
		/>
	</div>
</template>
