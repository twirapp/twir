<script setup lang="ts">
import { computed, ref } from 'vue'
import { toast } from 'vue-sonner'

import { useSongRequestsApi } from '~~/layers/dashboard/api/song-requests.js'
import Player from '~~/layers/dashboard/components/songRequests/player.vue'
import VideosQueue from '~~/layers/dashboard/components/songRequests/queue.vue'
import SettingsModal from '~~/layers/dashboard/components/songRequests/settings.vue'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { t } = useI18n()
const isSettingsModalOpened = ref(false)
const openSettingsModal = () => (isSettingsModalOpened.value = true)

const youtubeModuleManager = useSongRequestsApi()
const youtubeModuleData = youtubeModuleManager.useSongRequestQuery()

const noCookie = computed(() => {
	return youtubeModuleData.data.value?.songRequests?.playerNoCookieMode ?? false
})

const channelApiKey = computed(() => {
	return youtubeModuleData.data.value?.songRequests?.channelApiKey ?? ''
})

const requestUrl = useRequestURL()

const widgetLink = computed(() => {
	if (!channelApiKey.value) return ''
	return `${requestUrl.origin}/w/${channelApiKey.value}/song-requests`
})

const overlayLink = computed(() => {
	if (!channelApiKey.value) return ''
	return `${requestUrl.origin}/o/${channelApiKey.value}/song-requests`
})

function copyLink(link: string, label: string) {
	if (!link) {
		toast.error('Failed to copy link to clipboard', { duration: 2500 })
		return
	}

	navigator.clipboard.writeText(link).then(() => {
		toast.success(`${label} link copied!`, { duration: 3000 })
	}).catch(() => {
		toast.error('Failed to copy link to clipboard', { duration: 2500 })
	})
}
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

	<Card v-if="channelApiKey" class="mt-4">
		<CardHeader>
			<CardTitle>Widget &amp; Overlay Links</CardTitle>
		</CardHeader>
		<CardContent class="space-y-3">
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium min-w-20">Widget:</span>
				<code class="flex-1 text-sm bg-muted px-2 py-1 rounded">{{ widgetLink }}</code>
				<Button variant="outline" size="sm" @click="copyLink(widgetLink, 'Widget')">
					Copy
				</Button>
			</div>
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium min-w-20">OBS Overlay:</span>
				<code class="flex-1 text-sm bg-muted px-2 py-1 rounded">{{ overlayLink }}</code>
				<Button variant="outline" size="sm" @click="copyLink(overlayLink, 'Overlay')">
					Copy
				</Button>
			</div>
		</CardContent>
	</Card>

	<SettingsModal v-model:open="isSettingsModalOpened" />
</template>
