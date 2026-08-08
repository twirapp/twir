<script setup lang="ts">
import type { Platform } from '~/gql/graphql.js'
import { PLATFORM_META } from '~/utils/platforms.js'

withDefaults(
	defineProps<{
		platform: Platform
		size?: number
	}>(),
	{ size: 12 },
)
</script>

<template>
	<Icon
		:name="PLATFORM_META[platform].icon"
		class="platform-icon"
		:style="{ width: size + 'px', height: size + 'px' }"
		aria-hidden="true"
	/>
</template>

<style scoped>
.platform-icon {
	display: inline-block;
	flex-shrink: 0;
}
</style>
