<script lang="ts" setup>
import {
	IconBan,
	IconEye,
	IconEyeOff,
	IconLink,
	IconPlayerPauseFilled,
	IconPlayerPlayFilled,
	IconPlayerSkipForwardFilled,
	IconPlaylist,
	IconSettings,
	IconUser,
	IconVolume,
	IconVolume3,
} from '@tabler/icons-vue'
import { useLocalStorage, useScriptTag } from '@vueuse/core'
import { NCard, NEmpty, NList, NListItem, NPopconfirm, NResult, NSlider, NSpin } from 'naive-ui'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api/index.js'
import { useYoutubeSocket } from '@/components/songRequests/hook.js'
import { Button } from '@/components/ui/button'
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

function formatLabelTime(v: number) {
	return `${convertMillisToTime(v * 1000)}/${convertMillisToTime(duration.value * 1000)}`
}

function handlePlayPause() {
	if (!player.value || !playerReady.value) return
	if (isPlaying.value) {
		player.value.pauseVideo()
	} else {
		player.value.playVideo()
	}
}

function handleSeek(time: number) {
	if (!player.value || !playerReady.value) return
	try {
		player.value.seekTo(time, true)
	} catch (e) {
		console.error('Error seeking:', e)
	}
}

const { data: profile } = useProfile()
const { t } = useI18n()
</script>

<template>
	<NCard
		:title="t('songRequests.player.title')"
		content-style="padding: 0;"
		header-style="padding: 10px;"
		segmented
	>
		<div v-if="profile?.id !== profile?.selectedDashboardId" class="p-2.5">
			<NResult status="404" :title="t('songRequests.player.noAccess')" size="small"> </NResult>
		</div>

		<div v-else>
			<div v-show="playerVisible" class="player-wrapper">
				<div id="yt-player-container" class="player-container"></div>
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
						<IconPlayerPlayFilled v-if="!isPlaying" class="size-4" />
						<IconPlayerPauseFilled v-else class="size-4" />
					</Button>

					<Button
						size="icon"
						class="flex size-8 min-w-8"
						variant="secondary"
						:disabled="currentVideo == null"
						@click="playNext"
					>
						<IconPlayerSkipForwardFilled class="size-4" />
					</Button>

					<NSlider
						:value="sliderTime"
						:format-tooltip="formatLabelTime"
						:step="1"
						:max="duration"
						placement="bottom"
						:disabled="!currentVideo"
						@update-value="handleSeek"
					/>
				</div>

				<div class="flex items-center gap-2">
					<Button
						size="icon"
						variant="secondary"
						class="size-8 min-w-8"
						@click="isMuted = !isMuted"
					>
						<IconVolume v-if="!isMuted" class="size-4" />
						<IconVolume3 v-else class="size-4" />
					</Button>
					<NSlider :value="sliderVolume" :step="1" @update-value="(v) => (volume = v)" />
				</div>
			</div>
		</div>

		<template #header-extra>
			<div class="flex gap-2">
				<Button
					variant="secondary"
					size="icon"
					class="size-8"
					@click="playerVisible = !playerVisible"
				>
					<IconEyeOff v-if="playerVisible" class="size-4" />
					<IconEye v-else class="size-4" />
				</Button>
				<Button variant="secondary" size="icon" class="size-8" @click="openSettingsModal">
					<IconSettings class="size-4" />
				</Button>
			</div>
		</template>
		<template #footer>
			<template v-if="currentVideo">
				<NList :show-divider="false">
					<NListItem>
						<template #prefix>
							<IconPlaylist class="flex" />
						</template>

						{{ currentVideo?.title }}
					</NListItem>

					<NListItem>
						<template #prefix>
							<IconUser class="card-icon" />
						</template>

						{{ currentVideo?.orderedByDisplayName || currentVideo?.orderedByName }}
					</NListItem>

					<NListItem>
						<template #prefix>
							<IconLink class="card-icon" />
						</template>

						<a
							class="underline"
							:href="currentVideo.songLink ?? `https://youtu.be/${currentVideo?.videoId}`"
							target="_blank"
						>
							{{ currentVideo.songLink || `youtu.be/${currentVideo?.videoId}` }}
						</a>
					</NListItem>
				</NList>
				<div class="flex gap-2 justify-end">
					<NPopconfirm
						:positive-text="t('deleteConfirmation.confirm')"
						:negative-text="t('deleteConfirmation.cancel')"
						@positive-click="() => banSong(currentVideo.videoId)"
					>
						<template #trigger>
							<Button size="sm" variant="destructive">
								<div class="flex gap-1 items-center">
									<IconBan />
									{{ t('songRequests.ban.song') }}
								</div>
							</Button>
						</template>

						{{ t('songRequests.ban.songConfirm') }}
					</NPopconfirm>
					<NPopconfirm
						:positive-text="t('deleteConfirmation.confirm')"
						:negative-text="t('deleteConfirmation.cancel')"
						@positive-click="() => banUser(currentVideo.orderedById)"
					>
						<template #trigger>
							<Button size="sm" variant="destructive">
								<div class="flex gap-1 items-center">
									<IconBan />
									{{ t('songRequests.ban.user') }}
								</div>
							</Button>
						</template>

						{{ t('songRequests.ban.userConfirm') }}
					</NPopconfirm>
				</div>
			</template>
			<NEmpty v-else :description="t('songRequests.waiting')">
				<template #icon>
					<NSpin size="small" stroke="#959596" />
				</template>
			</NEmpty>
		</template>
	</NCard>
</template>

<style scoped>
.player-wrapper {
	height: 300px;
	overflow: hidden;
}

.player-container {
	position: relative;
	width: 100%;
	height: 100%;
}

.player-container :deep(iframe) {
	width: 100%;
	height: 300px !important;
	display: block;
}
</style>
