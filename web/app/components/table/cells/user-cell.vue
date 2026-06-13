<script setup lang="ts">
import Avatar from '~/components/ui/avatar/Avatar.vue'
import AvatarFallback from '~/components/ui/avatar/AvatarFallback.vue'
import AvatarImage from '~/components/ui/avatar/AvatarImage.vue'
import Badge from '~/components/ui/badge/Badge.vue'
import { Platform } from '~/gql/graphql.js'

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
		v-bind="linkAttrs"
	>
		<Avatar class="size-9">
			<AvatarImage :src="avatar" :alt="name" loading="lazy" />
			<AvatarFallback>{{ fallbackLetter }}</AvatarFallback>
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
