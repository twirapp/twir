<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import type { BrbOnStartFn, BrbOnStopFn } from '@/types.js'

import BrbTimer, { type BrbTimerMethods } from '@/components/brb-timer.vue'
import { useBrbEmotes } from '@/composables/brb/use-brb-emotes.js'
import { useBeRightBackOverlayGraphQL } from '@/composables/brb/use-brb-graphql.js'
import { useBrbIframe } from '@/composables/brb/use-brb-iframe.js'
import { useBrbSettings } from '@/composables/brb/use-brb-settings.js'

const route = useRoute()
const brbTimerRef = ref<BrbTimerMethods | null>(null)
const isPreviewMode = route.query.preview === '1'

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

const { settings } = useBrbSettings()

function startPreviewTimer() {
	brbTimerRef.value?.start(5, '')
}

watch(settings, (value) => {
	if (!isPreviewMode || !value) return
	startPreviewTimer()
})

watch(
	() => brbTimerRef.value?.isActive(),
	(active, previous) => {
		if (!isPreviewMode || active !== false || previous !== true) return
		startPreviewTimer()
	}
)

onMounted(() => {
	if (route.query.embed === 'settings') {
		iframe.create()
		return
	}

	const apiKey = route.params.apiKey as string
	if (!apiKey) {
		console.error('API key is required for Be Right Back overlay')
		return
	}
	graphql.connect(apiKey)
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
