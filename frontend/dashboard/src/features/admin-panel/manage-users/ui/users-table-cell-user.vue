<script setup lang="ts">
import { computed } from 'vue'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { resolveUserName } from '@/helpers/resolveUserName.js'

const props = defineProps<{
	name: string
	displayName: string
	avatar: string
	url?: string
}>()

const userName = computed(() => {
	return resolveUserName(props.name, props.displayName)
})

const Tag = computed(() => props.url ? 'a' : 'div')
</script>

<template>
	<component
		:is="Tag"
		:class="['flex items-center gap-4 max-sm:justify-start', url ? 'hover:underline' : '']"
		:href="url"
		target="_blank"
		rel="noopener noreferrer"
	>
		<Avatar class="size-9">
			<AvatarImage :src="avatar" :alt="name" loading="lazy" />
			<AvatarFallback>{{ name.charAt(0).toUpperCase() }}</AvatarFallback>
		</Avatar>
		{{ userName }}
	</component>
</template>
