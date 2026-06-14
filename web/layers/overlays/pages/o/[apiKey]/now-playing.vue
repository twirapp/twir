<script setup lang="ts">
import { NowPlaying, type Track } from '~~/lib/frontend-now-playing'
import { ref, watch } from 'vue'

import { useNowPlayingSocket } from '~~/layers/overlays/composables/now-playing/use-now-playing-socket

definePageMeta({ layout: false })

const route = useRoute()

const { currentTrack, settings } = useNowPlayingSocket({
	apiKey: route.params.apiKey as string,
	overlayId: route.query.id as string,
})

const showedTrack = ref<Track | null | undefined>(null)

const timerId = ref<number | null>(null)
watch(currentTrack, (track) => {
	if (timerId.value != null) {
		clearTimeout(timerId.value)
	}

	showedTrack.value = track

	if (settings.value?.hideTimeout) {
		timerId.value = setTimeout(() => {
			showedTrack.value = null
		}, settings.value.hideTimeout * 1000) as unknown as number
	}
})
</script>

<template>
	<NowPlaying v-if="settings" :settings="settings" :track="showedTrack" />
</template>