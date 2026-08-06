<script setup lang="ts">
import { toast } from 'vue-sonner'
import { useSongRequestsApi } from '~~/layers/dashboard/api/song-requests'

import { Button } from '@/components/ui/button'

const { t } = useI18n()

const songRequestsApi = useSongRequestsApi()
const settingsQuery = songRequestsApi.useSongRequestQuery()

const capabilities = computed(() => settingsQuery.data.value?.songRequests?.spotifyCapabilities)
const device = computed(() => capabilities.value?.selectedDevice ?? capabilities.value?.activeDevice)

const refreshDeviceMutation = songRequestsApi.useSpotifyRefreshDeviceMutation()
const refreshing = ref(false)

async function refreshDevice() {
	refreshing.value = true
	try {
		const result = await refreshDeviceMutation.executeMutation({})
		if (result.error) {
			toast.error(result.error.message, { duration: 5000 })
			return
		}
		toast.success(t('songRequests.spotify.deviceRefreshed'), { duration: 2500 })
	} finally {
		refreshing.value = false
	}
}
</script>

<template>
	<div
		v-if="capabilities?.canUseSpotify"
		class="border-border bg-card flex h-9 items-center gap-2 rounded-md border px-3 text-sm"
	>
		<Icon
			name="simple-icons:spotify"
			class="size-4 shrink-0 text-[#1DB954]"
		/>
		<span class="max-w-48 truncate">
			{{
				device
					? t('songRequests.spotify.device', { name: device.name, type: device.type })
					: t('songRequests.spotify.noDevice')
			}}
		</span>
		<Button
			size="icon"
			variant="ghost"
			class="size-6"
			:disabled="refreshing"
			:title="t('songRequests.spotify.refreshDevice')"
			@click="refreshDevice"
		>
			<Icon
				name="lucide:refresh-cw"
				class="size-3.5"
				:class="{ 'animate-spin': refreshing }"
			/>
		</Button>
	</div>
</template>
