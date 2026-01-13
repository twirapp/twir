<script setup lang="ts">
defineProps<{
	title?: string
}>()
</script>

<template>
	<div>
		<Head>
			<Title>{{ title || 'Hastebin' }}</Title>
		</Head>
		<div class="flex h-full bg-[#1E1E1E]">
			<slot />
		</div>
	</div>
</template>
