<script setup lang="ts">


defineProps<{
	imageUrl?: string
	name: string
}>()
</script>

<template>
	<div class="flex items-center gap-4 max-sm:justify-start">
		<UiAvatar class="size-9">
			<UiAvatarImage v-if="imageUrl" :src="imageUrl" :alt="name" loading="lazy" />
			<UiAvatarFallback>{{ name }}</UiAvatarFallback>
		</UiAvatar>
		{{ name }}
	</div>
</template>
