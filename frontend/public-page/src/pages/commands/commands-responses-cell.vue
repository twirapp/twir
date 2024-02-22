<script setup lang="ts">
defineProps<{
	responses: string[],
	description?: string
}>();
</script>

<template>
	<div v-if="!description" class="flex flex-col break-words">
		<span v-for="(response, responseIndex) of responses" :key="responseIndex">
			{{ response }}
		</span>
	</div>
	<div v-else>
		{{ description }}
	</div>
</template>
