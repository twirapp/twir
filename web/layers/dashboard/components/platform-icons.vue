<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { Platform } from '~/gql/graphql.js'

const props = defineProps<{
	platforms: Platform[]
	class?: HTMLAttributes['class']
}>()

const platformOptions = [
	{ id: Platform.Twitch, label: 'Twitch', icon: 'simple-icons:twitch', colorClass: 'text-[#9146FF]' },
	{ id: Platform.Kick, label: 'Kick', icon: 'simple-icons:kick', colorClass: 'text-[#53FC18]' },
	{ id: Platform.VkVideoLive, label: 'VK Video Live', icon: 'simple-icons:vk', colorClass: 'text-[#0077FF]' },
	{ id: Platform.Youtube, label: 'YouTube', icon: 'simple-icons:youtube', colorClass: 'text-[#FF0000]' },
] as const

const visiblePlatforms = computed(() =>
	platformOptions.filter((option) => props.platforms.includes(option.id))
)
</script>

<template>
	<div
		v-if="visiblePlatforms.length"
		class="flex items-center gap-1"
		:class="props.class"
	>
		<Icon
			v-for="platform in visiblePlatforms"
			:key="platform.id"
			:name="platform.icon"
			:title="platform.label"
			class="size-3.5"
			:class="platform.colorClass"
		/>
	</div>
</template>
