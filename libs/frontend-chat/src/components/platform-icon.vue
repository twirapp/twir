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
		case 'vk_video_live':
			return '#0077FF'
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
	<svg
		v-else-if="platform === 'vk_video_live'"
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 24 24"
		width="1em"
		height="1em"
		:fill="color"
		class="platform-icon"
	>
		<path
			d="m9.489.004l.729-.003h3.564l.73.003l.914.01l.433.007l.418.011l.403.014l.388.016l.374.021l.36.025l.345.03l.333.033c1.74.196 2.933.616 3.833 1.516s1.32 2.092 1.516 3.833l.034.333l.029.346l.025.36l.02.373l.025.588l.012.41l.013.644l.009.915l.004.98l-.001 3.313l-.003.73l-.01.914l-.007.433l-.011.418l-.014.403l-.016.388l-.021.374l-.025.36l-.03.345l-.033.333c-.196 1.74-.616 2.933-1.516 3.833s-2.092 1.32-3.833 1.516l-.333.034l-.346.029l-.36.025l-.373.02l-.588.025l-.41.012l-.644.013l-.915.009l-.98.004l-3.313-.001l-.73-.003l-.914-.01l-.433-.007l-.418-.011l-.403-.014l-.388-.016l-.374-.021l-.36-.025l-.345-.03l-.333-.033c-1.74-.196-2.933-.616-3.833-1.516S.367 20.221.171 18.48l-.034-.333l-.029-.346l-.025-.36l-.02-.373l-.025-.588l-.012-.41l-.013-.644l-.009-.915l-.004-.98l.001-3.313l.003-.73l.01-.914l.007-.433l.011-.418l.014-.403l.016-.388l.021-.374l.025-.36l.03-.345l.033-.333c.196-1.74.616-2.933 1.516-3.833S3.779.367 5.52.171l.333-.034l.346-.029l.36-.025l.373-.02l.588-.025l.41-.012l.644-.013zM6.79 7.3H4.05c.13 6.24 3.25 9.99 8.72 9.99h.31v-3.57c2.01.2 3.53 1.67 4.14 3.57h2.84c-.78-2.84-2.83-4.41-4.11-5.01c1.28-.74 3.08-2.54 3.51-4.98h-2.58c-.56 1.98-2.22 3.78-3.8 3.95V7.3H10.5v6.92c-1.6-.4-3.62-2.34-3.71-6.92"
		/>
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
