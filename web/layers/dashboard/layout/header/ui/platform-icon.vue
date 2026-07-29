<script setup lang="ts">
import { computed } from 'vue'

import { Platform } from '~/gql/graphql.js'

interface Props {
	platform: Platform
	live?: boolean
	size?: 'sm' | 'md'
}

const props = withDefaults(defineProps<Props>(), {
	live: false,
	size: 'md',
})

const iconName = computed(() => {
	switch (props.platform) {
		case Platform.Twitch:
			return 'simple-icons:twitch'
		case Platform.Kick:
			return 'simple-icons:kick'
		case Platform.VkVideoLive:
			return 'simple-icons:vk'
		case Platform.Youtube:
			return 'simple-icons:youtube'
	}
})

const iconColor = computed(() => {
	switch (props.platform) {
		case Platform.Twitch:
			return 'text-[#9146FF]'
		case Platform.Kick:
			return 'text-[#53FC18]'
		case Platform.VkVideoLive:
			return 'text-[#0077FF]'
		case Platform.Youtube:
			return 'text-[#FF0000]'
	}
})
</script>

<template>
	<span class="relative inline-flex flex-none items-center">
		<Icon
			:name="iconName"
			:class="[iconColor, size === 'sm' ? 'h-3 w-3' : 'h-3.5 w-3.5']"
		/>
		<span
			v-if="live"
			class="absolute -top-0.5 -right-0.5 h-1.5 w-1.5 rounded-full bg-red-500"
		/>
	</span>
</template>
