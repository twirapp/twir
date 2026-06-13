<script setup lang="ts">
import { onMounted } from 'vue'

import { useSpotifyIntegration } from '~~/layers/dashboard/api/integrations/spotify.js'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

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
		<Icon name="lucide:loader-circle" class="animate-spin size-12" />
	</div>
</template>
