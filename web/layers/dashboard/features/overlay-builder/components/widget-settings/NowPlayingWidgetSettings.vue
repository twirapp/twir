<script setup lang="ts">
import { NowPlaying, Preset } from '@twir/frontend-now-playing'
import { watch } from 'vue'

import { useNowPlayingOverlayApi } from '~~/layers/dashboard/api/overlays/now-playing'
import NowPlayingForm from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/now-playing-form.vue'
import { useNowPlayingForm } from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/use-now-playing-form'
import { useNowPlayingPreviewTrack } from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/use-now-playing-track'

const { data: entities } = useNowPlayingOverlayApi().useNowPlayingQuery()
const { data: formData, setData } = useNowPlayingForm()
const { track } = useNowPlayingPreviewTrack()

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
	<div class="flex flex-col gap-4">
		<NowPlaying
			:settings="formData ?? { preset: Preset.TRANSPARENT }"
			:track="track"
		/>
		<NowPlayingForm embedded />
	</div>
</template>
