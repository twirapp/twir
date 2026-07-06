<script setup lang="ts">
import { useQuery, useSubscription } from '@urql/vue'
import { useRoute } from 'vue-router'
import { graphql } from '~/gql/gql.js'

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

const playbackState = computed(() => playbackSub.data.value?.songRequestPlaybackState ?? null)

const playerReady = ref(false)
const duration = ref(0)
const currentPosition = ref(0)
let player: any = null
let tickInterval: ReturnType<typeof setInterval> | null = null

onMounted(() => {
	const tag = document.createElement('script')
	tag.src = 'https://www.youtube.com/iframe_api'
	document.head.appendChild(tag)

	;(window as any).onYouTubeIframeAPIReady = () => {
		player = new (window as any).YT.Player('yt-player', {
			height: '100%',
			width: '100%',
			playerVars: {
				autoplay: 0,
				controls: 0,
				modestbranding: 1,
				rel: 0,
				disablekb: 1,
				fs: 0,
			},
			events: {
				onReady: () => {
					playerReady.value = true
				},
				onStateChange: (event: any) => {
					if (event.data === (window as any).YT.PlayerState.PLAYING) {
						duration.value = player.getDuration() ?? 0
					}
				},
			},
		})
	}

	tickInterval = setInterval(() => {
		if (player && playerReady.value && playbackState.value?.isPlaying) {
			const pos = player.getCurrentTime?.()
			if (pos !== undefined) {
				currentPosition.value = pos
			}
		}
	}, 1000)
})

onUnmounted(() => {
	if (tickInterval) {
		clearInterval(tickInterval)
	}
})

const currentVideoId = ref('')

watch(playbackState, (state) => {
	if (!state || !playerReady.value) return

	currentPosition.value = state.position

	if (state.videoId !== currentVideoId.value) {
		currentVideoId.value = state.videoId
		player.loadVideoById(state.videoId, state.position)
	} else {
		player.seekTo(state.position, true)
	}

	player.setVolume(state.volume)

	if (state.isPlaying) {
		player.playVideo()
	} else {
		player.pauseVideo()
	}
}, { deep: true })

const isVisible = computed(() => {
	if (!playbackState.value) return false
	return playbackState.value.isPlaying
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
	<div v-if="isVisible" class="song-request-overlay">
		<div class="player-container">
			<div id="yt-player" class="yt-player" />
		</div>
		<div class="track-info">
			<div class="track-title">
				{{ playbackState?.title }}
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

<style scoped>
.song-request-overlay {
	position: relative;
	width: 100%;
	max-width: 640px;
	background: rgba(0, 0, 0, 0.85);
	border-radius: 8px;
	overflow: hidden;
}

.player-container {
	position: relative;
	width: 100%;
	padding-top: 56.25%;
}

.yt-player {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
}

.track-info {
	padding: 12px 16px;
}

.track-title {
	color: white;
	font-size: 16px;
	font-weight: 600;
	margin-bottom: 8px;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.progress-bar {
	width: 100%;
	height: 4px;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 2px;
	margin-bottom: 4px;
}

.progress-fill {
	height: 100%;
	background: #8b5cf6;
	border-radius: 2px;
	transition: width 0.5s linear;
}

.progress-time {
	color: rgba(255, 255, 255, 0.6);
	font-size: 12px;
}
</style>
