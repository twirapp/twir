<script setup lang="ts">
import { NowPlaying } from '@twir/frontend-now-playing'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import { useNowPlayingSocket } from '@/composables/now-playing/use-now-playing-socket.ts'
import { useVisibleTrack } from '@/composables/now-playing/use-visible-track.ts'

const route = useRoute()

const { currentTrack, settings } = useNowPlayingSocket({
	apiKey: route.params.apiKey as string,
	overlayId: route.query.id as string,
})

const hideTimeout = computed(() => {
	const timeoutSeconds = settings.value?.hideTimeout
	return timeoutSeconds == null ? timeoutSeconds : timeoutSeconds * 1000
})
const { visibleTrack } = useVisibleTrack(currentTrack, hideTimeout)
</script>

<template>
	<NowPlaying v-if="settings" :settings="settings" :track="visibleTrack" />
</template>
