<script setup lang="ts">
import { computed, ref } from 'vue'

import { useSongRequestsApi } from '@/api/song-requests.js'
import { useYoutubeSocket } from '@/components/songRequests/hook.js'
import Player from '@/components/songRequests/player.vue'
import VideosQueue from '@/components/songRequests/queue.vue'
import SettingsModal from '@/components/songRequests/settings.vue'

const isSettingsModalOpened = ref(false)
const openSettingsModal = () => (isSettingsModalOpened.value = true)

useYoutubeSocket()

const youtubeModuleManager = useSongRequestsApi()
const youtubeModuleData = youtubeModuleManager.useSongRequestQuery()

const noCookie = computed(() => {
	return youtubeModuleData.data.value?.songRequests?.playerNoCookieMode ?? false
})
</script>

<template>
	<div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
		<div class="lg:col-span-1">
			<Player
				v-if="!youtubeModuleData.fetching.value"
				:no-cookie="noCookie"
				:open-settings-modal="openSettingsModal"
			/>
		</div>

		<div class="lg:col-span-2">
			<VideosQueue />
		</div>
	</div>

	<SettingsModal v-model:open="isSettingsModalOpened" />
</template>
