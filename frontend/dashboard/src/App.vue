<script setup lang="ts">
import { provideClient } from '@urql/vue'
import { onMounted, ref, watch } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { Loader } from 'lucide-vue-next'

import { urqlClient } from './plugins/urql.js'
import { useYoutubeSocket } from './components/songRequests/hook.js'
import { useGlobalYoutubePlayer } from './composables/useGlobalYoutubePlayer.js'

const isRouterReady = ref(false)
const router = useRouter()
router.isReady().finally(() => (isRouterReady.value = true))

provideClient(urqlClient)

const { currentVideo } = useYoutubeSocket()
const { playerReady, initPlayer } = useGlobalYoutubePlayer()

// Initialize player when there's a current video
onMounted(() => {
	watch(currentVideo, async (video) => {
		if (video && !playerReady.value) {
			console.log('[App] Initializing player for video:', video.videoId)
			// Wait a bit for the DOM to be ready
			await new Promise(resolve => setTimeout(resolve, 200))
			await initPlayer()
		}
	}, { immediate: true })
})
</script>

<template>
	<div v-if="!isRouterReady" class="flex justify-center items-center h-full">
		<Loader class="h-10 w-10 sonner-spinner" />
	</div>
	<template v-else>
		<router-view />
		<!-- Global YouTube player container - always in DOM, never destroyed, never moved -->
		<!-- Position controlled by player.vue via inline styles -->
		<div
			id="global-yt-player-container"
			class="fixed pointer-events-none"
			style="left: -9999px; top: -9999px; width: 640px; height: 360px;"
		></div>
	</template>
</template>
