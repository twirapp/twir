<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

import { useProfile } from '~~/layers/dashboard/api/auth'
import { useSongRequestsApi } from '~~/layers/dashboard/api/song-requests'
import { useGlobalYoutubePlayer } from '~~/layers/dashboard/composables/useGlobalYoutubePlayer.js'
import { useSongRequestGql } from '~~/layers/dashboard/composables/useSongRequestGql.js'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardTitle } from '@/components/ui/card'
import { Slider } from '@/components/ui/slider'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { convertMillisToTime } from '~~/layers/dashboard/helpers/convertMillisToTime.js'

import type { SongRequestsSettingsOpts } from '~/gql/graphql.js'

const props = defineProps<{
	noCookie: boolean
	openSettingsModal: () => void
}>()

const { data: profile } = useProfile()

const channelId = computed(() => profile.value?.selectedDashboardId ?? '')

const {
	playbackState,
	queue,
	play,
	pause,
	skip,
	setVolume: setVolumeGql,
	deleteFromQueue,
} = useSongRequestGql(channelId)

const {
	playerReady,
	isPlaying,
	sliderTime,
	duration,
	playerVisible,
	isMuted,
	sliderVolume,
	noCookie,
	initPlayer,
	playVideo,
	pauseVideo,
	seekTo,
	setVolume: setPlayerVolume,
	loadVideoById,
	cueVideoById,
	stopVideo,
	toggleMute,
} = useGlobalYoutubePlayer()

const currentVideo = computed(() => playbackState.value)

const currentQueueItem = computed(() => {
	if (!currentVideo.value) return null
	return queue.value.find((item) => item.title === currentVideo.value?.title) ?? null
})

const youtubeModuleManager = useSongRequestsApi()
const { data: youtubeSettings } = youtubeModuleManager.useSongRequestQuery()
const youtubeModuleUpdater = youtubeModuleManager.useSongRequestMutation()

const playerContainerRef = ref<HTMLElement | null>(null)

// Update noCookie setting
watch(() => props.noCookie, (value) => {
	noCookie.value = value
}, { immediate: true })

function updatePlayerPosition() {
	if (!playerContainerRef.value) return

	const rect = playerContainerRef.value.getBoundingClientRect()
	const globalContainer = document.getElementById('global-yt-player-container')

	if (globalContainer) {
		globalContainer.style.left = `${rect.left}px`
		globalContainer.style.top = `${rect.top}px`
		globalContainer.style.width = `${rect.width}px`
		globalContainer.style.height = `${rect.height}px`
		console.log('[Player] Updated position:', rect)
	}
}

onMounted(async () => {
	console.log('[Player] Mounted, initializing player if needed')

	// Initialize player if not already initialized
	if (!playerReady.value) {
		console.log('[Player] Initializing player...')
		await initPlayer()
	}

	// Wait a bit for DOM to settle
	await new Promise(resolve => setTimeout(resolve, 100))

	// Position the global player over our container
	updatePlayerPosition()

	// Update position on window resize
	window.addEventListener('resize', updatePlayerPosition)
	window.addEventListener('scroll', updatePlayerPosition, true)
})

onUnmounted(() => {
	console.log('[Player] Unmounted, moving iframe off-screen')
	window.removeEventListener('resize', updatePlayerPosition)
	window.removeEventListener('scroll', updatePlayerPosition, true)

	// Reset iframe position to off-screen
	const globalContainer = document.getElementById('global-yt-player-container')
	if (globalContainer) {
		globalContainer.style.left = '-9999px'
		globalContainer.style.top = '-9999px'
		globalContainer.style.width = '640px'
		globalContainer.style.height = '360px'
		console.log('[Player] Iframe moved off-screen')
	}
})

const formattedTime = computed(() => {
	return `${convertMillisToTime(sliderTime.value * 1000)} / ${convertMillisToTime(duration.value * 1000)}`
})

function handleSeek(value: number[] | undefined) {
	if (!value) return
	seekTo(value[0])
}

function handleVolumeChange(value: number[] | undefined) {
	if (!value) return
	setPlayerVolume(value[0])
	setVolumeGql(value[0])
}

function handlePlayPause() {
	if (isPlaying.value) {
		pause()
		pauseVideo()
	} else {
		play()
		playVideo()
	}
}

function handleSkip() {
	skip()
}

function toSettingsOpts(
	settings: Record<string, unknown>,
): SongRequestsSettingsOpts {
	const { channelApiKey, __typename, ...rest } = settings
	return rest as unknown as SongRequestsSettingsOpts
}

async function banUser(userId: string) {
	if (!youtubeSettings.value?.songRequests) return
	const settings = toSettingsOpts(youtubeSettings.value.songRequests as Record<string, unknown>)

	await youtubeModuleUpdater.executeMutation({
		opts: {
			...settings,
			denyList: {
				...settings.denyList,
				users: [...settings.denyList!.users, userId],
			},
		},
	})

	// Delete all videos from this user
	const userVideos = queue.value.filter((video) => video.orderedByName === userId)
	for (const video of userVideos) {
		deleteFromQueue(video.id)
	}
}

async function banSong(videoId: string) {
	if (!youtubeSettings.value?.songRequests) return
	const settings = toSettingsOpts(youtubeSettings.value.songRequests as Record<string, unknown>)

	await youtubeModuleUpdater.executeMutation({
		opts: {
			...settings,
			denyList: {
				...settings.denyList,
				songs: [...settings.denyList!.songs, videoId],
			},
		},
	})

	deleteFromQueue(videoId)
}

// Watch playback state changes from GQL subscription
watch(playbackState, (state) => {
	if (!playerReady.value) return

	if (!state) {
		stopVideo()
		sliderTime.value = 0
		duration.value = 0
		return
	}

	// Handle video changes
	if (state.videoId) {
		loadVideoById(state.videoId)
		if (state.position > 0) {
			seekTo(state.position)
		}
	}

	// Handle play/pause
	if (state.isPlaying) {
		playVideo()
	} else {
		pauseVideo()
	}

	// Handle volume
	if (state.volume !== undefined) {
		setPlayerVolume(state.volume)
	}
}, { deep: true })

const { t } = useI18n()

const selectedDashboard = computed(() => {
	return profile.value?.availableDashboards.find(
		(dashboard) => dashboard.id === profile.value?.selectedDashboardId
	)
})

const canUsePlayer = computed(() => {
	if (!profile.value || !selectedDashboard.value) {
		return false
	}

	return profile.value.linkedAccounts.some((account) => {
		if (account.platform !== selectedDashboard.value?.platform) {
			return false
		}

		if (account.platform === 'twitch') {
			return account.platformLogin === selectedDashboard.value.twitchProfile?.login
		}

		if (account.platform === 'kick') {
			return account.platformUserId === String(selectedDashboard.value.kickProfile?.id)
		}

		return false
	})
})
</script>

<template>
	<Card class="p-0">
		<CardContent class="p-0">
			<div class="flex flex-row justify-between items-center px-2 py-2 border-b">
				<CardTitle class="text-base">{{ t('songRequests.player.title') }}</CardTitle>
				<div class="flex gap-1">
					<Button
						variant="outline"
						size="icon"
						class="size-8"
						:disabled="!canUsePlayer"
						@click="playerVisible = !playerVisible"
					>
						<Icon name="lucide:eye-off" v-if="playerVisible" class="size-4" />
						<Icon name="lucide:eye" v-else class="size-4" />
					</Button>
					<Button variant="outline" size="icon" class="size-8" :disabled="!canUsePlayer" @click="openSettingsModal">
						<Icon name="lucide:settings" class="size-4" />
					</Button>
				</div>
			</div>
			<div v-if="!canUsePlayer" class="p-6 text-center">
				<p class="text-muted-foreground">{{ t('songRequests.player.noAccess') }}</p>
			</div>

			<div v-else>
				<div
					v-show="playerVisible"
					ref="playerContainerRef"
					class="h-[300px] overflow-hidden relative bg-black"
				>
					<!-- Global YouTube iframe will be positioned over this container -->
				</div>

				<div class="flex flex-col gap-4 py-5 px-6">
					<div class="flex gap-2 items-center">
						<Button
							size="icon"
							class="flex size-8 min-w-8"
							variant="secondary"
							:disabled="currentVideo == null"
							@click="handlePlayPause"
						>
							<Icon name="lucide:play" v-if="!isPlaying" class="size-4" />
							<Icon name="lucide:pause" v-else class="size-4" />
						</Button>

						<Button
							size="icon"
							class="flex size-8 min-w-8"
							variant="secondary"
							:disabled="currentVideo == null"
							@click="handleSkip"
						>
							<Icon name="lucide:skip-forward" class="size-4" />
						</Button>

						<Slider
							:model-value="[sliderTime]"
							:step="1"
							:max="duration || 1"
							:disabled="!currentVideo"
							@update:model-value="handleSeek"
						/>
						<span class="text-xs text-muted-foreground whitespace-nowrap">{{ formattedTime }}</span>
					</div>

					<div class="flex items-center gap-2">
						<Button
							size="icon"
							variant="secondary"
							class="size-8 min-w-8"
							@click="toggleMute"
						>
							<Icon name="lucide:volume2" v-if="!isMuted" class="size-4" />
							<Icon name="lucide:volume-x" v-else class="size-4" />
						</Button>
						<Slider
							:model-value="[sliderVolume]"
							:step="1"
							:max="100"
							@update:model-value="handleVolumeChange"
						/>
					</div>
				</div>
			</div>
		</CardContent>

		<CardFooter class="flex-col items-start gap-4 border-t pt-2">
			<template v-if="currentVideo">
				<div class="flex flex-col gap-2 w-full">
					<div class="flex items-center gap-2">
						<Icon name="lucide:list-music" class="size-4 shrink-0" />
						<span class="truncate">{{ currentVideo?.title }}</span>
					</div>

					<div v-if="currentQueueItem" class="flex items-center gap-2">
						<Icon name="lucide:user" class="size-4 shrink-0" />
						<span>{{ currentQueueItem?.orderedByDisplayName || currentQueueItem?.orderedByName }}</span>
					</div>

					<div v-if="currentQueueItem" class="flex items-center gap-2">
						<Icon name="lucide:link" class="size-4 shrink-0" />
						<a
							class="underline text-sm truncate"
							:href="currentQueueItem.songLink ?? `https://youtu.be/${currentVideo?.videoId}`"
							target="_blank"
						>
							{{ currentQueueItem.songLink || `youtu.be/${currentVideo?.videoId}` }}
						</a>
					</div>
				</div>

				<div v-if="currentQueueItem" class="flex gap-2 mb-2 justify-end w-full">
					<AlertDialog>
						<AlertDialogTrigger as-child>
							<Button size="sm" variant="destructive">
								<Icon name="lucide:ban" class="size-4 mr-1" />
								{{ t('songRequests.ban.song') }}
							</Button>
						</AlertDialogTrigger>
						<AlertDialogContent>
							<AlertDialogHeader>
								<AlertDialogTitle>{{ t('songRequests.ban.songConfirm') }}</AlertDialogTitle>
							</AlertDialogHeader>
							<AlertDialogFooter>
								<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
								<AlertDialogAction @click="banSong(currentVideo.videoId)">
									{{ t('deleteConfirmation.confirm') }}
								</AlertDialogAction>
							</AlertDialogFooter>
						</AlertDialogContent>
					</AlertDialog>

					<AlertDialog>
						<AlertDialogTrigger as-child>
							<Button size="sm" variant="destructive">
								<Icon name="lucide:ban" class="size-4 mr-1" />
								{{ t('songRequests.ban.user') }}
							</Button>
						</AlertDialogTrigger>
						<AlertDialogContent>
							<AlertDialogHeader>
								<AlertDialogTitle>{{ t('songRequests.ban.userConfirm') }}</AlertDialogTitle>
							</AlertDialogHeader>
							<AlertDialogFooter>
								<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
								<AlertDialogAction @click="banUser(currentQueueItem.orderedByName)">
									{{ t('deleteConfirmation.confirm') }}
								</AlertDialogAction>
							</AlertDialogFooter>
						</AlertDialogContent>
					</AlertDialog>
				</div>
			</template>

			<div v-else class="flex flex-col items-center justify-center w-full py-4 gap-2">
				<Icon name="lucide:loader2" class="size-6 animate-spin text-muted-foreground" />
				<p class="text-muted-foreground text-sm">{{ t('songRequests.waiting') }}</p>
			</div>
		</CardFooter>
	</Card>
</template>
