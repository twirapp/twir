<script setup lang="ts">
import { useNow } from '@vueuse/core'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import CustomTemplate from '@/components/stream-stats/custom-template.vue'
import DesignBar from '@/components/stream-stats/design-bar.vue'
import DesignCards from '@/components/stream-stats/design-cards.vue'
import DesignMinimal from '@/components/stream-stats/design-minimal.vue'
import {
	type StreamStatsCounters,
	type StreamStatsSettings,
	useStreamStatsCounters,
} from '@/composables/stream-stats/use-stream-stats-counters.js'
import { useStreamStatsGraphQL } from '@/composables/stream-stats/use-stream-stats-graphql.js'
import { Platform, StreamStatsOverlayDesign } from '@/gql/graphql'

const route = useRoute()
const isPreviewMode = route.query.preview === '1'

const graphql = useStreamStatsGraphQL()
const settingsOverrides = ref<Partial<StreamStatsSettings>>({})

const previewCounters: StreamStatsCounters = {
	live: true,
	viewers: 1342,
	messages: 5489,
	startedAt: new Date(Date.now() - 83 * 60 * 1000).toISOString(),
	subscribers: 156,
	followers: 12043,
	platformViewers: [
		{ platform: Platform.Twitch, viewers: 1000 },
		{ platform: Platform.Kick, viewers: 300 },
		{ platform: Platform.Youtube, viewers: 42 },
	],
}

const settings = computed<StreamStatsSettings | null>(() => {
	const base = graphql.settings.value
	if (!base) return null
	return { ...base, ...settingsOverrides.value }
})

const counters = computed<StreamStatsCounters | null>(() => {
	if (isPreviewMode) return previewCounters
	return graphql.counters.value
})

const now = useNow({ interval: 1000 })
const { items, placeholders } = useStreamStatsCounters(settings, counters, now)

const designComponent = computed(() => {
	switch (settings.value?.design) {
		case StreamStatsOverlayDesign.Bar:
			return DesignBar
		case StreamStatsOverlayDesign.Minimal:
			return DesignMinimal
		case StreamStatsOverlayDesign.Cards:
		default:
			return DesignCards
	}
})

const isVisible = computed(() => Boolean(settings.value && counters.value?.live))

function onWindowMessage(msg: MessageEvent<string>) {
	if (!isPreviewMode) return

	let parsed: { key?: string, data?: Partial<StreamStatsSettings> }
	try {
		parsed = JSON.parse(msg.data)
	} catch {
		return
	}

	if (parsed?.key !== 'settings' || !parsed.data) return
	settingsOverrides.value = { ...settingsOverrides.value, ...parsed.data }
}

onMounted(() => {
	if (isPreviewMode) {
		window.addEventListener('message', onWindowMessage)
		window.parent.postMessage(JSON.stringify({ key: 'getSettings' }), '*')
	}

	const apiKey = route.params.apiKey as string
	if (!apiKey) {
		console.error('API key is required for Stream Stats overlay')
		return
	}

	graphql.connect(apiKey)
})

onUnmounted(() => {
	window.removeEventListener('message', onWindowMessage)
	graphql.destroy()
})
</script>

<template>
	<template v-if="settings && isVisible">
		<CustomTemplate
			v-if="settings.customHtmlEnabled"
			:html="settings.customHtml"
			:css="settings.customCss"
			:values="placeholders"
		/>
		<component :is="designComponent" v-else :items="items" />
	</template>
</template>
