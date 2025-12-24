<script lang="ts" setup>
import {
	Ban,
	Eye,
	EyeOff,
	Link,
	ListMusic,
	Loader2,
	Pause,
	Play,
	Settings,
	SkipForward,
	User,
	Volume2,
	VolumeX,
} from 'lucide-vue-next'
import { useLocalStorage, useScriptTag } from '@vueuse/core'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api/index.js'
import { useYoutubeSocket } from '@/components/songRequests/hook.js'
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
import { convertMillisToTime } from '@/helpers/convertMillisToTime.js'

interface YTPlayer {
	playVideo: () => void
	pauseVideo: () => void
	stopVideo: () => void
	loadVideoById: (videoId: string) => void
	cueVideoById: (videoId: string) => void
	seekTo: (seconds: number, allowSeekAhead: boolean) => void
	getPlayerState: () => number
	getCurrentTime: () => number
	getDuration: () => number
	setVolume: (volume: number) => void
	getVolume: () => number
	mute: () => void
	unMute: () => void
	isMuted: () => boolean
	destroy: () => void
}

const enum PlayerState {
	ENDED = 0,
	PLAYING = 1,
	PAUSED = 2,
	BUFFERING = 3,
}

interface YT {
	Player: new (
		elementId: string,
		options: {
			height?: string
			width?: string
			videoId?: string
			host?: string
			playerVars?: {
				autoplay?: number
				controls?: number
				rel?: number
				playsinline?: number
			}
			events?: {
				onReady?: (event: { target: YTPlayer }) => void
				onStateChange?: (event: { target: YTPlayer; data: number }) => void
				onError?: (event: { target: YTPlayer; data: number }) => void
			}
		}
	) => YTPlayer
	PlayerState: typeof PlayerState
}

declare global {
	interface Window {
		YT?: YT
		onYouTubeIframeAPIReady?: () => void
	}
}

const props = defineProps<{
	noCookie: boolean
	openSettingsModal: () => void
}>()

const { currentVideo, nextVideo, sendPlaying, banSong, banUser } = useYoutubeSocket()

const player = ref<YTPlayer>()
const playerReady = ref(false)
const isPlaying = ref(false)
const sliderTime = ref(0)
const duration = ref(0)
const updateTimeInterval = ref<number>()

const volume = useLocalStorage('twirPlayerVolume', 10)
const playerVisible = useLocalStorage('twirPlayerVisible', true)
const isMuted = useLocalStorage('twirPlayerIsMuted', false)

function playNext() {
	nextVideo()
}

function startTimeUpdate() {
	if (updateTimeInterval.value) {
		clearInterval(updateTimeInterval.value)
	}
	updateTimeInterval.value = window.setInterval(() => {
		if (player.value && playerReady.value) {
			try {
				const time = player.value.getCurrentTime()
				if (time !== undefined && !Number.isNaN(time)) {
					sliderTime.value = time
				}
			} catch (e) {
				console.warn('Could not get current time:', e)
			}
		}
	}, 100)
}

function stopTimeUpdate() {
	if (updateTimeInterval.value) {
		clearInterval(updateTimeInterval.value)
		updateTimeInterval.value = undefined
	}
}

function onPlayerReady(event: { target: YTPlayer }) {
	player.value = event.target
	playerReady.value = true

	const vol = volume.value
	event.target.setVolume(vol)
	if (isMuted.value) {
		event.target.mute()
	}
}

function onPlayerStateChange(event: { target: YTPlayer; data: number }) {
	if (event.data === PlayerState.PLAYING) {
		isPlaying.value = true
		if (player.value && playerReady.value) {
			try {
				const dur = player.value.getDuration()
				if (dur && !Number.isNaN(dur)) {
					duration.value = dur
				}
			} catch (e) {
				console.warn('Could not get duration:', e)
			}
		}
		startTimeUpdate()
		sendPlaying()
	} else if (event.data === PlayerState.PAUSED) {
		isPlaying.value = false
		stopTimeUpdate()
	} else if (event.data === PlayerState.ENDED) {
		isPlaying.value = false
		stopTimeUpdate()
		playNext()
	} else if (event.data === PlayerState.BUFFERING) {
		if (player.value && playerReady.value) {
			try {
				const dur = player.value.getDuration()
				if (dur && !Number.isNaN(dur)) {
					duration.value = dur
				}
			} catch (e) {
				console.warn('Could not get duration while buffering:', e)
			}
		}
	}
}

function onPlayerError(event: { target: YTPlayer; data: number }) {
	console.error('YouTube Player Error:', event.data)
	isPlaying.value = false
	stopTimeUpdate()
}

const { load: loadYTScript } = useScriptTag('https://www.youtube.com/iframe_api', () => {}, {
	manual: true,
})

onMounted(async () => {
	const initPlayer = () => {
		if (!window.YT || !window.YT.Player) return

		const playerConfig: any = {
			height: '300',
			width: '100%',
			host: props.noCookie ? 'https://www.youtube-nocookie.com' : undefined,
			playerVars: {
				autoplay: 0,
				controls: 1,
				rel: 0,
				playsinline: 1,
			},
			events: {
				onReady: onPlayerReady,
				onStateChange: onPlayerStateChange,
				onError: onPlayerError,
			},
		}

		if (currentVideo.value) {
			playerConfig.videoId = currentVideo.value.videoId
		}

		// oxlint-disable-next-line no-new
		new window.YT.Player('yt-player-container', playerConfig)
	}

	if (window.YT && window.YT.Player) {
		await nextTick()
		initPlayer()
	} else {
		window.onYouTubeIframeAPIReady = () => {
			nextTick(() => {
				initPlayer()
			})
		}
		await loadYTScript()
	}
})

watch(currentVideo, (video) => {
	if (!player.value || !playerReady.value) return

	if (!video) {
		player.value.stopVideo()
		sliderTime.value = 0
		duration.value = 0
		return
	}

	player.value.loadVideoById(video.videoId)
})

onUnmounted(() => {
	stopTimeUpdate()
	if (player.value) {
		player.value.destroy()
	}
	playerReady.value = false
})

watch(volume, (value) => {
	if (!player.value || !playerReady.value) return

	if (value === 0) {
		isMuted.value = true
		player.value.mute()
	} else {
		isMuted.value = false
		player.value.unMute()
		player.value.setVolume(value)
	}
})

watch(isMuted, (v) => {
	if (!player.value || !playerReady.value) return
	if (v) {
		player.value.mute()
	} else {
		player.value.unMute()
	}
})

const sliderVolume = computed(() => {
	if (isMuted.value) return 0
	return volume.value
})

const formattedTime = computed(() => {
	return `${convertMillisToTime(sliderTime.value * 1000)} / ${convertMillisToTime(duration.value * 1000)}`
})

function handlePlayPause() {
	if (!player.value || !playerReady.value) return
	if (isPlaying.value) {
		player.value.pauseVideo()
	} else {
		player.value.playVideo()
	}
}

function handleSeek(value: number[] | undefined) {
	if (!player.value || !playerReady.value || !value) return
	try {
		player.value.seekTo(value[0], true)
	} catch (e) {
		console.error('Error seeking:', e)
	}
}

function handleVolumeChange(value: number[] | undefined) {
	if (!value) return
	volume.value = value[0]
}

const { data: profile } = useProfile()
const { t } = useI18n()
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
						@click="playerVisible = !playerVisible"
					>
						<EyeOff v-if="playerVisible" class="size-4" />
						<Eye v-else class="size-4" />
					</Button>
					<Button variant="outline" size="icon" class="size-8" @click="openSettingsModal">
						<Settings class="size-4" />
					</Button>
				</div>
			</div>
			<div v-if="profile?.id !== profile?.selectedDashboardId" class="p-6 text-center">
				<p class="text-muted-foreground">{{ t('songRequests.player.noAccess') }}</p>
			</div>

			<div v-else>
				<div v-show="playerVisible" class="h-[300px] overflow-hidden">
					<div
						id="yt-player-container"
						class="w-full h-full [&_iframe]:w-full [&_iframe]:!h-[300px] [&_iframe]:block"
					></div>
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
							<Play v-if="!isPlaying" class="size-4" />
							<Pause v-else class="size-4" />
						</Button>

						<Button
							size="icon"
							class="flex size-8 min-w-8"
							variant="secondary"
							:disabled="currentVideo == null"
							@click="playNext"
						>
							<SkipForward class="size-4" />
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
							@click="isMuted = !isMuted"
						>
							<Volume2 v-if="!isMuted" class="size-4" />
							<VolumeX v-else class="size-4" />
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
						<ListMusic class="size-4 shrink-0" />
						<span class="truncate">{{ currentVideo?.title }}</span>
					</div>

					<div class="flex items-center gap-2">
						<User class="size-4 shrink-0" />
						<span>{{ currentVideo?.orderedByDisplayName || currentVideo?.orderedByName }}</span>
					</div>

					<div class="flex items-center gap-2">
						<Link class="size-4 shrink-0" />
						<a
							class="underline text-sm truncate"
							:href="currentVideo.songLink ?? `https://youtu.be/${currentVideo?.videoId}`"
							target="_blank"
						>
							{{ currentVideo.songLink || `youtu.be/${currentVideo?.videoId}` }}
						</a>
					</div>
				</div>

				<div class="flex gap-2 mb-2 justify-end w-full">
					<AlertDialog>
						<AlertDialogTrigger as-child>
							<Button size="sm" variant="destructive">
								<Ban class="size-4 mr-1" />
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
								<Ban class="size-4 mr-1" />
								{{ t('songRequests.ban.user') }}
							</Button>
						</AlertDialogTrigger>
						<AlertDialogContent>
							<AlertDialogHeader>
								<AlertDialogTitle>{{ t('songRequests.ban.userConfirm') }}</AlertDialogTitle>
							</AlertDialogHeader>
							<AlertDialogFooter>
								<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
								<AlertDialogAction @click="banUser(currentVideo.orderedById)">
									{{ t('deleteConfirmation.confirm') }}
								</AlertDialogAction>
							</AlertDialogFooter>
						</AlertDialogContent>
					</AlertDialog>
				</div>
			</template>

			<div v-else class="flex flex-col items-center justify-center w-full py-4 gap-2">
				<Loader2 class="size-6 animate-spin text-muted-foreground" />
				<p class="text-muted-foreground text-sm">{{ t('songRequests.waiting') }}</p>
			</div>
		</CardFooter>
	</Card>
</template>
