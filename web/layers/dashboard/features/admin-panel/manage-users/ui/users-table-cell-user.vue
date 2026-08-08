<script setup lang="ts">
import { computed } from 'vue'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import PlatformBadge from '@/components/platform/platform-badge.vue'
import type { Platform } from '~/gql/graphql.js'
import { resolveUserName } from '~~/layers/dashboard/helpers/resolveUserName.js'

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
		<div class="flex min-w-0 flex-col gap-1">
			<span class="truncate">{{ userName }}</span>
			<PlatformBadge
				v-if="platform"
				:platform="platform"
			/>
		</div>
	</component>
</template>
