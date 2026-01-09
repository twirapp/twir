<script setup lang="ts">
import { computed, ref } from 'vue'

import { useSongRequestsApi } from '#layers/dashboard/api/song-requests'
import { useYoutubeSocket } from '#layers/dashboard/components/songRequests/hook.js'
import Player from '#layers/dashboard/components/songRequests/player.vue'
import VideosQueue from '#layers/dashboard/components/songRequests/queue.vue'
import SettingsModal from '#layers/dashboard/components/songRequests/settings.vue'

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
