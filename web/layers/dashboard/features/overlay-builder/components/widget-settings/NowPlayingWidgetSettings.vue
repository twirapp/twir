<script setup lang="ts">
import { watch } from 'vue'

import { useNowPlayingOverlayApi } from '~~/layers/dashboard/api/overlays/now-playing'
import NowPlayingForm from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/now-playing-form.vue'
import { useNowPlayingForm } from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/use-now-playing-form'

const { data: entities } = useNowPlayingOverlayApi().useNowPlayingQuery()
const { setData } = useNowPlayingForm()

watch(
	() => entities.value?.nowPlayingOverlays,
	(overlays) => {
		const firstOverlay = overlays?.[0]
		if (firstOverlay) setData(firstOverlay)
	},
	{ immediate: true },
)
</script>

<template>
	<NowPlayingForm embedded />
</template>
