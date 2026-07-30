<script setup lang="ts">
import type { Platform } from '~/gql/graphql.js'

import { getPlatformMeta } from '~/utils/platforms.js'

interface Props {
	platform: Platform
	live?: boolean
	size?: 'sm' | 'md'
}

const props = withDefaults(defineProps<Props>(), {
	live: false,
	size: 'md',
})

const meta = computed(() => getPlatformMeta(props.platform))
</script>

<template>
	<span class="relative inline-flex flex-none items-center">
		<Icon
			:name="meta.icon"
			:title="meta.label"
			:class="[meta.colorClass, size === 'sm' ? 'h-3 w-3' : 'h-3.5 w-3.5']"
		/>
		<span
			v-if="live"
			class="absolute -top-0.5 -right-0.5 h-1.5 w-1.5 rounded-full bg-red-500"
		/>
	</span>
</template>
