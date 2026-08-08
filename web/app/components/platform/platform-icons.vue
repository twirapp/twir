<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import type { Platform } from '~/gql/graphql.js'

import { PLATFORM_OPTIONS } from '~/utils/platforms.js'

const props = defineProps<{
	platforms: Platform[]
	class?: HTMLAttributes['class']
}>()

const visiblePlatforms = computed(() =>
	PLATFORM_OPTIONS.filter((option) => props.platforms.includes(option.platform)),
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
			:key="platform.platform"
			:name="platform.icon"
			:title="platform.label"
			class="size-3.5"
			:class="platform.colorClass"
		/>
	</div>
</template>
