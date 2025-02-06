<script setup lang="ts">
import { resolveUserName } from '@/lib/users'

const props = defineProps<{
	name: string
	displayName: string
	avatar: string
}>()

const userName = computed(() => {
	return resolveUserName(props.name, props.displayName)
})
</script>

<template>
	<div class="flex items-center gap-4 max-sm:justify-start">
		<img :src="avatar" :alt="name" class="size-9 rounded-full" />
		<!-- <UiAvatar class="size-9">
			<UiAvatarImage :src="avatar" :alt="name" />
			<UiAvatarFallback>{{ name.charAt(0).toUpperCase() }}</UiAvatarFallback>
		</UiAvatar> -->
		{{ userName }}
	</div>
</template>
