<script lang="ts" setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useProfile } from '~~/layers/dashboard/api/auth'
import { useSongRequestsApi } from '~~/layers/dashboard/api/song-requests'
import { useDebouncedVolume } from '~~/layers/dashboard/composables/useDebouncedVolume.js'
import { useGlobalYoutubePlayer } from '~~/layers/dashboard/composables/useGlobalYoutubePlayer.js'
import { useSongRequestGql } from '~~/layers/dashboard/composables/useSongRequestGql.js'
import { convertMillisToTime } from '~~/layers/dashboard/helpers/convertMillisToTime.js'

import type { SongRequestsSettingsOpts } from '~/gql/graphql.js'

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
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardTitle } from '@/components/ui/card'
import { Slider } from '@/components/ui/slider'

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
	updatePosition,
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

const hasPlayableVideo = computed(() => !!playbackState.value || queue.value.length > 0)

const currentQueueItem = computed(() => {
	if (!currentVideo.value) return null
	return queue.value.find((item) => item.title === currentVideo.value?.title) ?? null
})

// Auto-cue first video from queue when nothing is loaded
watch([queue, playerReady], ([items, ready]) => {
	if (!ready || playbackState.value || !items.length) return
	const first = items[0]
	if (first?.songLink) {
		const match = first.songLink.match(/(?:v=|youtu\.be\/)([^&?/]+)/)
		if (match) {
			cueVideoById(match[1]!)
		}
	}
})

const youtubeModuleManager = useSongRequestsApi()
const { data: youtubeSettings } = youtubeModuleManager.useSongRequestQuery()
const youtubeModuleUpdater = youtubeModuleManager.useSongRequestMutation()

const playerContainerRef = ref<HTMLElement | null>(null)

// Update noCookie setting
watch(
	() => props.noCookie,
	(value) => {
		noCookie.value = value
	},
	{ immediate: true }
)

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
	await new Promise((resolve) => setTimeout(resolve, 100))

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
	seekTo(value[0]!)
}

const { localVolume, handleVolumeInput, syncFromServer } = useDebouncedVolume(
	() => channelId.value,
	setVolumeGql
)

function handleVolumeChange(value: number[] | undefined) {
	if (!value) return
	setPlayerVolume(value[0]!)
	handleVolumeInput(value[0]!)
}

function handlePlayPause() {
	if (isPlaying.value) {
		pause()
		pauseVideo()
	} else {
		// Get videoId from queue if no playback state
		const videoId =
			playbackState.value?.videoId ||
			(queue.value[0]?.songLink?.match(/(?:v=|youtu\.be\/)([^&?/]+)/)?.[1] ?? '')
		if (videoId) {
			play(videoId)
			// If same video is already loaded, just resume. Otherwise load new.
			if (videoId === currentVideoId.value) {
				playVideo()
			} else {
				loadVideoById(videoId)
			}
		}
	}
}

function handleSkip() {
	skip()
}

function toSettingsOpts(settings: Record<string, unknown>): SongRequestsSettingsOpts {
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

const currentVideoId = ref('')

// Watch playback state changes from GQL subscription
watch(
	playbackState,
	(state) => {
		if (!playerReady.value) return

		isUpdatingFromSubscription.value = true

		if (!state) {
			stopVideo()
			currentVideoId.value = ''
			sliderTime.value = 0
			duration.value = 0
			isUpdatingFromSubscription.value = false
			return
		}

		// Only load new video if videoId changed
		if (state.videoId && state.videoId !== currentVideoId.value) {
			currentVideoId.value = state.videoId
			loadVideoById(state.videoId)
		} else if (state.position > 0) {
			seekTo(state.position)
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
			syncFromServer(state.volume)
		}

		lastSyncedPosition.value = state.position
		lastKnownIsPlaying.value = state.isPlaying

		// Reset flag after a tick to let watchers see the state change
		nextTick(() => {
			isUpdatingFromSubscription.value = false
		})
	},
	{ deep: true }
)

// Detect user-initiated play/pause from YouTube player controls
const isUpdatingFromSubscription = ref(false)
const lastKnownIsPlaying = ref(false)

watch(isPlaying, (nowPlaying, wasPlaying) => {
	if (isUpdatingFromSubscription.value) return
	if (!playbackState.value) return
	if (nowPlaying === wasPlaying) return

	if (nowPlaying) {
		const videoId =
			playbackState.value.videoId ||
			(queue.value[0]?.songLink?.match(/(?:v=|youtu\.be\/)([^&?/]+)/)?.[1] ?? '')
		if (videoId) play(videoId)
	} else {
		pause()
	}
	lastKnownIsPlaying.value = nowPlaying
})

// Detect user-initiated seek from YouTube player (only when paused)
const lastSyncedPosition = ref(0)
let seekDebounce: ReturnType<typeof setTimeout> | null = null

watch(sliderTime, (newTime) => {
	if (isUpdatingFromSubscription.value) return
	if (!playbackState.value || playbackState.value.isPlaying) return

	const diff = Math.abs(newTime - lastSyncedPosition.value)
	if (diff < 2) return

	if (seekDebounce) clearTimeout(seekDebounce)
	seekDebounce = setTimeout(() => {
		updatePosition(newTime)
		lastSyncedPosition.value = newTime
	}, 300)
})

// Update lastSyncedPosition when subscription delivers new state
watch(
	() => playbackState.value?.position,
	(pos) => {
		if (pos !== undefined) {
			lastSyncedPosition.value = pos
		}
	}
)

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
			<div class="flex flex-row items-center justify-between border-b px-2 py-2">
				<CardTitle class="text-base">{{ t('songRequests.player.title') }}</CardTitle>
				<div class="flex gap-1">
					<Button
						variant="outline"
						size="icon"
						class="size-8"
						:disabled="!canUsePlayer"
						@click="playerVisible = !playerVisible"
					>
						<Icon
							name="lucide:eye-off"
							v-if="playerVisible"
							class="size-4"
						/>
						<Icon
							name="lucide:eye"
							v-else
							class="size-4"
						/>
					</Button>
					<Button
						variant="outline"
						size="icon"
						class="size-8"
						:disabled="!canUsePlayer"
						@click="openSettingsModal"
					>
						<Icon
							name="lucide:settings"
							class="size-4"
						/>
					</Button>
				</div>
			</div>
			<div
				v-if="!canUsePlayer"
				class="p-6 text-center"
			>
				<p class="text-muted-foreground">{{ t('songRequests.player.noAccess') }}</p>
			</div>

			<div v-else>
				<div
					v-show="playerVisible"
					ref="playerContainerRef"
					class="relative h-[300px] overflow-hidden bg-black"
				>
					<!-- Global YouTube iframe will be positioned over this container -->
				</div>

				<div class="flex flex-col gap-4 px-6 py-5">
					<div class="flex items-center gap-2">
						<Button
							size="icon"
							class="flex size-8 min-w-8"
							variant="secondary"
							:disabled="!hasPlayableVideo"
							@click="handlePlayPause"
						>
							<Icon
								name="lucide:play"
								v-if="!isPlaying"
								class="size-4"
							/>
							<Icon
								name="lucide:pause"
								v-else
								class="size-4"
							/>
						</Button>

						<Button
							size="icon"
							class="flex size-8 min-w-8"
							variant="secondary"
							:disabled="!hasPlayableVideo"
							@click="handleSkip"
						>
							<Icon
								name="lucide:skip-forward"
								class="size-4"
							/>
						</Button>

						<Slider
							:model-value="[sliderTime]"
							:step="1"
							:max="duration || 1"
							:disabled="!currentVideo"
							@update:model-value="handleSeek"
						/>
						<span class="text-muted-foreground text-xs whitespace-nowrap">{{ formattedTime }}</span>
					</div>

					<div class="flex items-center gap-2">
						<Button
							size="icon"
							variant="secondary"
							class="size-8 min-w-8"
							@click="toggleMute"
						>
							<Icon
								name="lucide:volume-2"
								v-if="!isMuted"
								class="size-4"
							/>
							<Icon
								name="lucide:volume-x"
								v-else
								class="size-4"
							/>
						</Button>
						<Slider
							:model-value="[localVolume]"
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
				<div class="flex w-full flex-col gap-2">
					<div class="flex items-center gap-2">
						<Icon
							name="lucide:list-music"
							class="size-4 shrink-0"
						/>
						<span class="truncate">{{ currentVideo?.title }}</span>
					</div>

					<div
						v-if="currentQueueItem"
						class="flex items-center gap-2"
					>
						<Icon
							name="lucide:user"
							class="size-4 shrink-0"
						/>
						<span>{{
							currentQueueItem?.orderedByDisplayName || currentQueueItem?.orderedByName
						}}</span>
					</div>

					<div
						v-if="currentQueueItem"
						class="flex items-center gap-2"
					>
						<Icon
							name="lucide:link"
							class="size-4 shrink-0"
						/>
						<a
							class="truncate text-sm underline"
							:href="currentQueueItem.songLink ?? `https://youtu.be/${currentVideo?.videoId}`"
							target="_blank"
						>
							{{ currentQueueItem.songLink || `youtu.be/${currentVideo?.videoId}` }}
						</a>
					</div>
				</div>

				<div
					v-if="currentQueueItem"
					class="mb-2 flex w-full justify-end gap-2"
				>
					<AlertDialog>
						<AlertDialogTrigger as-child>
							<Button
								size="sm"
								variant="destructive"
							>
								<Icon
									name="lucide:ban"
									class="mr-1 size-4"
								/>
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
							<Button
								size="sm"
								variant="destructive"
							>
								<Icon
									name="lucide:ban"
									class="mr-1 size-4"
								/>
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

			<div
				v-else
				class="flex w-full flex-col items-center justify-center gap-2 py-4"
			>
				<Icon
					name="lucide:loader2"
					class="text-muted-foreground size-6 animate-spin"
				/>
				<p class="text-muted-foreground text-sm">{{ t('songRequests.waiting') }}</p>
			</div>
		</CardFooter>
	</Card>
</template>
