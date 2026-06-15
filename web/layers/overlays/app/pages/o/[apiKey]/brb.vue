<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

import type { BrbOnStartFn, BrbOnStopFn } from '~~/layers/overlays/app/types'

import BrbTimer, { type BrbTimerMethods } from '~~/layers/overlays/app/components/brb-timer.vue'
import { useBrbEmotes } from '~~/layers/overlays/app/composables/brb/use-brb-emotes'
import { useBeRightBackOverlayGraphQL } from '~~/layers/overlays/app/composables/brb/use-brb-graphql'
import { useBrbIframe } from '~~/layers/overlays/app/composables/brb/use-brb-iframe'

definePageMeta({ layout: false })

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

const emotes = useBrbEmotes()

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
	emotes.destroy()
})
</script>

<template>
	<div id="brb-container mx-auto" class="container">
		<BrbTimer ref="brbTimerRef" />
	</div>
</template>

<style scoped>
.container {
	overflow: hidden;
}
</style>