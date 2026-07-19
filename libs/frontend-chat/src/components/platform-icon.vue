<script setup lang="ts">
import { computed } from 'vue'

import type { MessagePlatform } from '../types.js'

const props = defineProps<{
	platform: MessagePlatform
}>()

const color = computed(() => {
	switch (props.platform) {
		case 'twitch':
			return '#9146FF'
		case 'kick':
			return '#53FC18'
		default:
			return 'currentColor'
	}
})
</script>

<template>
	<svg
		v-if="platform === 'twitch'"
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 24 24"
		width="1em"
		height="1em"
		:fill="color"
		class="platform-icon"
	>
		<path
			d="M11.571 4.714h1.715v5.143H11.57zm4.715 0H18v5.143h-1.714zM6 0L1.714 4.286v15.428h5.143V24l4.286-4.286h3.428L22.286 12V0zm14.571 11.143l-3.428 3.428h-3.429l-3 3v-3H6.857V1.714h13.714Z"
		/>
	</svg>
	<svg
		v-else-if="platform === 'kick'"
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 24 24"
		width="1em"
		height="1em"
		:fill="color"
		class="platform-icon"
	>
		<path d="M3 2h4v8l6-8h5l-7 9 7 11h-5l-6-9v9H3Z" />
	</svg>
</template>

<style scoped>
.platform-icon {
	height: 1em;
	width: 1em;
	display: inline-block;
	vertical-align: middle;
}
</style>
