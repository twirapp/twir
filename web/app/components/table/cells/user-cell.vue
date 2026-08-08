<script setup lang="ts">
import PlatformBadge from '~/components/platform/platform-badge.vue'
import Avatar from '~/components/ui/avatar/Avatar.vue'
import AvatarFallback from '~/components/ui/avatar/AvatarFallback.vue'
import AvatarImage from '~/components/ui/avatar/AvatarImage.vue'
import type { Platform } from '~/gql/graphql.js'

const props = defineProps<{
	name: string
	displayName: string
	avatar: string
	url?: string
	platform?: Platform | null
}>()

const userName = computed(() => {
	return resolveUserName(props.name, props.displayName)
})

const Tag = computed(() => props.url ? 'a' : 'div')

const linkAttrs = computed(() => {
	if (!props.url) return {}

	return {
		href: props.url,
		target: '_blank',
		rel: 'noopener noreferrer',
	}
})

const fallbackLetter = computed(() => userName.value.charAt(0).toUpperCase() || '?')
</script>

<template>
	<component
		:is="Tag"
		:class="['flex items-center gap-4 max-sm:justify-start', url ? 'hover:underline' : '']"
		v-bind="linkAttrs"
	>
		<Avatar class="size-9">
			<AvatarImage :src="avatar" :alt="name" loading="lazy" />
			<AvatarFallback>{{ fallbackLetter }}</AvatarFallback>
		</Avatar>
		<div class="flex min-w-0 flex-col gap-1">
			<span class="truncate">{{ userName }}</span>
			<PlatformBadge
				v-if="platform"
				:platform="platform"
			/>
		</div>
	</component>
</template>
