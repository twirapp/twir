<script setup lang="ts">
import { computed } from 'vue'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Platform } from '~/gql/graphql.js'
import { resolveUserName } from '~~/layers/dashboard/app/helpers/resolveUserName.js'

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

const platformMeta = computed(() => {
	if (props.platform === Platform.Kick) {
		return {
			label: 'Kick',
			className: 'border-[#53FC18]/30 bg-[#53FC18]/10 text-[#3CB30F]',
		}
	}

	if (props.platform === Platform.Twitch) {
		return {
			label: 'Twitch',
			className: 'border-[#9146FF]/30 bg-[#9146FF]/10 text-[#7C3AED]',
		}
	}

	return null
})
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
			<Badge
				v-if="platformMeta"
				variant="outline"
				class="w-fit text-[10px] uppercase tracking-wide"
				:class="platformMeta.className"
			>
				{{ platformMeta.label }}
			</Badge>
		</div>
	</component>
</template>
