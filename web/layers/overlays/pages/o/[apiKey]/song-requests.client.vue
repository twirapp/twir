<script setup lang="ts">
import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { useRoute } from 'vue-router'
import { graphql } from '~/gql/gql.js'

interface YTPlayer {
	loadVideoById: (videoId: string, startSeconds?: number) => void
	seekTo: (seconds: number, allowSeekAhead: boolean) => void
	playVideo: () => void
	pauseVideo: () => void
	stopVideo: () => void
	setVolume: (volume: number) => void
	getCurrentTime: () => number
	getDuration: () => number
}

interface YTApi {
	Player: new (
		elementId: string,
		options: {
			height: string
			width: string
			playerVars: Record<string, number>
			events: {
				onReady: () => void
				onStateChange: (event: { data: number }) => void
			}
		},
	) => YTPlayer
	PlayerState: {
		PLAYING: number
		PAUSED: number
		ENDED: number
	}
}

type WindowWithYouTube = Window & {
	YT?: YTApi
	onYouTubeIframeAPIReady?: () => void
}

const route = useRoute()
const apiKey = computed(() => route.params.apiKey as string)

const channelQuery = useQuery({
	query: graphql(`
		query OverlaySongRequestsChannelByApiKey($apiKey: String!) {
			channelByApiKey(apiKey: $apiKey) {
				id
			}
		}
	`),
	variables: computed(() => ({ apiKey: apiKey.value })),
	pause: computed(() => !apiKey.value),
})

const channelId = computed(() => channelQuery.data.value?.channelByApiKey?.id ?? '')

const widgetDataQuery = useQuery({
	query: graphql(`
		query OverlaySongRequestWidgetData($channelId: UUID!) {
			songRequestWidgetData(channelId: $channelId) {
				playbackState {
					videoId
					title
					position
					isPlaying
					volume
					updatedAt
				}
				queue {
					id
					title
					songLink
					durationSeconds
					orderedByName
					orderedByDisplayName
					queuePosition
					createdAt
				}
			}
		}
	`),
	variables: computed(() => ({ channelId: channelId.value })),
	pause: computed(() => !channelId.value),
})

const playbackSub = useSubscription({
	query: graphql(`
		subscription OverlaySongRequestPlaybackState($channelId: UUID!) {
			songRequestPlaybackState(channelId: $channelId) {
				videoId
				title
				position
				isPlaying
				volume
				updatedAt
			}
		}
	`),
	variables: computed(() => ({ channelId: channelId.value })),
	pause: computed(() => !channelId.value),
})

const queueSub = useSubscription({
	query: graphql(`
		subscription OverlaySongRequestQueueUpdated($channelId: UUID!) {
			songRequestQueueUpdated(channelId: $channelId) {
				id
				title
				songLink
				durationSeconds
				orderedByName
				orderedByDisplayName
				queuePosition
				createdAt
			}
		}
	`),
	variables: computed(() => ({ channelId: channelId.value })),
	pause: computed(() => !channelId.value),
})

const playbackState = computed(() =>
	playbackSub.data.value?.songRequestPlaybackState
	?? widgetDataQuery.data.value?.songRequestWidgetData?.playbackState
	?? null,
)

const queue = computed(() =>
	queueSub.data.value?.songRequestQueueUpdated
	?? widgetDataQuery.data.value?.songRequestWidgetData?.queue
	?? [],
)

const currentQueueItem = computed(() => {
	if (!playbackState.value?.videoId) return null

	return queue.value.find((item) => item.id === playbackState.value?.videoId) ?? null
})

const requesterName = computed(() => {
	return currentQueueItem.value?.orderedByDisplayName || currentQueueItem.value?.orderedByName || ''
})

const { executeMutation: playMutation } = useMutation(graphql(`
	mutation OverlaySongRequestPlay($channelId: UUID!, $videoId: String!) {
		songRequestPlay(channelId: $channelId, videoId: $videoId)
	}
`))

const { executeMutation: skipMutation } = useMutation(graphql(`
	mutation OverlaySongRequestSkip($channelId: UUID!) {
		songRequestSkip(channelId: $channelId)
	}
`))

const { executeMutation: updatePositionMutation } = useMutation(graphql(`
	mutation OverlaySongRequestUpdatePosition($channelId: UUID!, $position: Float!) {
		songRequestUpdatePosition(channelId: $channelId, position: $position)
	}
`))

const playerReady = ref(false)
const playerIsPlaying = ref(false)
const duration = ref(0)
const currentPosition = ref(0)
const currentVideoId = ref('')
const lastAppliedIsPlaying = ref<boolean | null>(null)
const skippedVideoId = ref('')

let player: YTPlayer | null = null
let tickInterval: ReturnType<typeof setInterval> | null = null

function getWindowWithYouTube(): WindowWithYouTube {
	return window as WindowWithYouTube
}

function createPlayer() {
	const win = getWindowWithYouTube()
	if (!win.YT?.Player || player) return

	player = new win.YT.Player('yt-player', {
		height: '100%',
		width: '100%',
		playerVars: {
			autoplay: 0,
			controls: 0,
			playsinline: 1,
			modestbranding: 1,
			rel: 0,
			disablekb: 1,
			fs: 0,
		},
		events: {
			onReady: () => {
				playerReady.value = true
			},
			onStateChange: (event) => {
				if (!win.YT) return

				if (event.data === win.YT.PlayerState.PLAYING) {
					playerIsPlaying.value = true
					duration.value = player?.getDuration() ?? 0
					return
				}

				if (event.data === win.YT.PlayerState.PAUSED) {
					playerIsPlaying.value = false
					return
				}

				if (event.data === win.YT.PlayerState.ENDED) {
					playerIsPlaying.value = false
					if (currentVideoId.value && channelId.value) {
						skippedVideoId.value = currentVideoId.value
						void skipMutation({ channelId: channelId.value })
					}
				}
			},
		},
	})
}

onMounted(() => {
	const win = getWindowWithYouTube()
	if (win.YT?.Player) {
		createPlayer()
	} else {
		win.onYouTubeIframeAPIReady = createPlayer

		if (!document.querySelector('script[src="https://www.youtube.com/iframe_api"]')) {
			const tag = document.createElement('script')
			tag.src = 'https://www.youtube.com/iframe_api'
			document.head.appendChild(tag)
		}
	}

	tickInterval = setInterval(() => {
		if (!player || !playerReady.value || !playerIsPlaying.value || !channelId.value) return

		const position = player.getCurrentTime()
		if (Number.isNaN(position)) return

		currentPosition.value = position
		void updatePositionMutation({ channelId: channelId.value, position })
	}, 1000)
})

onUnmounted(() => {
	if (tickInterval) {
		clearInterval(tickInterval)
	}
})

watch(
	[playbackState, playerReady],
	([state, ready]) => {
		if (!ready || !player) return

		if (!state?.videoId) {
			player.stopVideo()
			currentVideoId.value = ''
			lastAppliedIsPlaying.value = null
			currentPosition.value = 0
			return
		}

		if (state.videoId !== currentVideoId.value) {
			currentVideoId.value = state.videoId
			skippedVideoId.value = ''
			currentPosition.value = state.position
			player.loadVideoById(state.videoId, state.position)
			duration.value = player.getDuration() || 0
		} else if (Math.abs(state.position - currentPosition.value) > 2.5) {
			currentPosition.value = state.position
			player.seekTo(state.position, true)
		}

		player.setVolume(state.volume)

		if (state.isPlaying !== lastAppliedIsPlaying.value) {
			if (state.isPlaying) {
				player.playVideo()
			} else {
				player.pauseVideo()
			}
			lastAppliedIsPlaying.value = state.isPlaying
		}
	},
	{ deep: true, immediate: true },
)

watch(
	[queue, playbackState, playerReady],
	([items, state, ready]) => {
		if (!ready || !channelId.value || state?.videoId) return

		const nextVideoId = items[0]?.id ?? ''
		if (!nextVideoId || nextVideoId === skippedVideoId.value) return

		void playMutation({ channelId: channelId.value, videoId: nextVideoId })
	},
	{ immediate: true },
)

const isVisible = computed(() => {
	return !!playbackState.value?.videoId && playbackState.value.isPlaying
})

const progressPercent = computed(() => {
	if (!duration.value || duration.value <= 0) return 0
	return Math.min((currentPosition.value / duration.value) * 100, 100)
})

function formatTime(seconds: number): string {
	const m = Math.floor(seconds / 60)
	const s = Math.floor(seconds % 60)
	return `${m}:${s.toString().padStart(2, '0')}`
}
</script>

<template>
	<div
		class="song-request-overlay"
		:class="{ 'song-request-overlay--hidden': !isVisible }"
	>
		<div class="player-container">
			<div id="yt-player" class="yt-player" />
		</div>
		<div class="track-info">
			<div class="track-title">
				{{ playbackState?.title }}
			</div>
			<div
				v-if="requesterName"
				class="track-requester"
			>
				Requested by {{ requesterName }}
			</div>
			<div class="progress-bar">
				<div
					class="progress-fill"
					:style="{ width: `${progressPercent}%` }"
				/>
			</div>
			<div class="progress-time">
				{{ formatTime(currentPosition) }}
				<template v-if="duration > 0">
					/ {{ formatTime(duration) }}
				</template>
			</div>
		</div>
	</div>
</template>

<style>
html,
body,
#__nuxt {
	background: transparent !important;
}

.song-request-overlay {
	position: fixed;
	inset: 0;
	display: flex;
	flex-direction: column;
	width: 100vw;
	height: 100vh;
	background: transparent;
	overflow: hidden;
}

.song-request-overlay--hidden {
	opacity: 0;
	pointer-events: none;
}

.song-request-overlay .player-container {
	position: relative;
	flex: 1;
	min-height: 0;
	width: 100%;
	background: transparent;
}

.song-request-overlay .yt-player {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
}

.song-request-overlay .track-info {
	flex-shrink: 0;
	padding: 14px 16px 12px;
	background: rgba(0, 0, 0, 0.85);
}

.song-request-overlay .track-title {
	color: white;
	font-size: 16px;
	font-weight: 600;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.song-request-overlay .track-requester {
	color: rgba(255, 255, 255, 0.7);
	font-size: 12px;
	margin-top: 4px;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.song-request-overlay .progress-bar {
	width: 100%;
	height: 4px;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 2px;
	margin-top: 10px;
	margin-bottom: 4px;
}

.song-request-overlay .progress-fill {
	height: 100%;
	background: #8b5cf6;
	border-radius: 2px;
	transition: width 0.5s linear;
}

.song-request-overlay .progress-time {
	color: rgba(255, 255, 255, 0.6);
	font-size: 12px;
}
</style>
