<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import type { BrbOnStartFn, BrbOnStopFn } from '@/types.js'

import BrbTimer, { type BrbTimerMethods } from '@/components/brb-timer.vue'
import { useBrbIframe } from '@/composables/brb/use-brb-iframe.js'
import { useBeRightBackOverlaySocket } from '@/composables/brb/use-brb-socket.js'

const route = useRoute()
const brbTimerRef = ref<BrbTimerMethods | null>(null)

const onStart: BrbOnStartFn = (minutes, text) => {
	brbTimerRef.value?.start(minutes, text)
}

const onStop: BrbOnStopFn = () => {
	brbTimerRef.value?.stop()
}

const iframe = useBrbIframe({
	onStart,
	onStop,
})

const apiKey = route.params.apiKey as string

const socket = useBeRightBackOverlaySocket({
	apiKey,
	onStart,
	onStop,
})

onMounted(() => {
	if (window.frameElement) {
		iframe.create()
	} else {
		socket.create()
	}
})

onUnmounted(() => {
	iframe.destroy()
	socket.destroy()
})
</script>

<template>
	<div id="brb-container" class="container">
		<BrbTimer ref="brbTimerRef" />
	</div>
</template>

<style scoped>
.container {
	overflow: hidden;
}
</style>
