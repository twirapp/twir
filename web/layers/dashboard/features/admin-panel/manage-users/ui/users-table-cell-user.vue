<script setup lang="ts">
import { computed } from 'vue'


import { resolveUserName } from '#layers/dashboard/helpers/resolveUserName'

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
		<UiAvatar class="size-9">
			<UiAvatarImage :src="avatar" :alt="name" loading="lazy" />
			<UiAvatarFallback>{{ name.charAt(0).toUpperCase() }}</UiAvatarFallback>
		</UiAvatar>
		{{ userName }}
	</div>
</template>
