<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import type { BrbOnStartFn, BrbOnStopFn } from '@/types.js'

import BrbTimer, { type BrbTimerMethods } from '@/components/brb-timer.vue'
import { useBeRightBackOverlayGraphQL } from '@/composables/brb/use-brb-graphql.js'
import { useBrbIframe } from '@/composables/brb/use-brb-iframe.js'

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

const graphql = useBeRightBackOverlayGraphQL({
	onStart,
	onStop,
})

onMounted(() => {
	if (window.frameElement) {
		iframe.create()
	} else {
		const apiKey = route.params.apiKey as string
		if (!apiKey) {
			console.error('API key is required for Be Right Back overlay')
			return
		}
		graphql.connect(apiKey)
	}
})

onUnmounted(() => {
	iframe.destroy()
	graphql.destroy()
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
