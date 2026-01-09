<script setup lang="ts">
import { LoaderCircleIcon } from 'lucide-vue-next'
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'

import { useSpotifyIntegration } from '#layers/dashboard/api/integrations/spotify'

const spotifyIntegration = useSpotifyIntegration()
const route = useRoute()

onMounted(async () => {
	const { code } = route.query
	if (typeof code !== 'string') {
		return window.close()
	}

	await spotifyIntegration.postCode.executeMutation({ input: { code } })
	spotifyIntegration.broadcastRefresh()
	window.close()
})
</script>

<template>
	<div class="flex justify-center items-center h-full">
		<LoaderCircleIcon class="animate-spin size-12" />
	</div>
</template>
