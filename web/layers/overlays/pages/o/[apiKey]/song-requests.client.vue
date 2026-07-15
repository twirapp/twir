<script setup lang="ts">
import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { useRoute } from 'vue-router'

import SongRequestOverlayRenderer from '~/components/song-request-overlay/SongRequestOverlayRenderer.vue'
import {
	type SongRequestOverlayStyle,
	normalizeSongRequestOverlayStyle,
} from '~/components/song-request-overlay/types.js'
import { graphql } from '~/gql/gql.js'

interface OverlaySettings {
	style: SongRequestOverlayStyle
	accentColor: string
	tickerBackgroundColor: string
	tickerTextColor: string
	tickerSpeed: number
	hideOnPause: boolean
}

interface ChannelScope {
	apiKey: string
	channelId: string
	apiKeyGeneration: number
	connectionGeneration: number
}

const DEFAULT_OVERLAY_SETTINGS: OverlaySettings = {
	style: 'CINEMA',
	accentColor: '#8B5CF6',
	tickerBackgroundColor: '#111827E6',
	tickerTextColor: '#FFFFFF',
	tickerSpeed: 35,
	hideOnPause: true,
}

const WIDGET_SNAPSHOT_RETRY_BASE_DELAY = 1000
const WIDGET_SNAPSHOT_RETRY_MAX_DELAY = 10_000
const CHANNEL_RESOLUTION_RETRY_BASE_DELAY = 1000
const CHANNEL_RESOLUTION_RETRY_MAX_DELAY = 10_000

interface YTPlayer {
	loadVideoById: (videoId: string, startSeconds?: number) => void
	seekTo: (seconds: number, allowSeekAhead: boolean) => void
	playVideo: () => void
	pauseVideo: () => void
	stopVideo: () => void
	setVolume: (volume: number) => void
	getCurrentTime: () => number
	getDuration: () => number
	getPlayerState: () => number
	getVideoUrl: () => string
	destroy: () => void
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
		}
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
const overlaySettings = ref<OverlaySettings>({ ...DEFAULT_OVERLAY_SETTINGS })
const apiKeyGeneration = ref(0)
const connectionGeneration = ref(0)
let isUnmounted = false

function handleGraphqlWsReconnect() {
	connectionGeneration.value += 1
}

if (import.meta.client) {
	window.addEventListener('twir:graphql-ws-reconnected', handleGraphqlWsReconnect)
	onScopeDispose(() => {
		window.removeEventListener('twir:graphql-ws-reconnected', handleGraphqlWsReconnect)
	})
}

const settingsSub = useSubscription({
	query: graphql(`
		subscription OverlaySongRequestSettings($apiKey: String!) {
			songRequestOverlaySettings(apiKey: $apiKey) {
				style
				accentColor
				tickerBackgroundColor
				tickerTextColor
				tickerSpeed
				hideOnPause
			}
		}
	`),
	variables: computed(() => ({ apiKey: apiKey.value })),
	pause: computed(() => !apiKey.value),
})

watch(
	() => settingsSub.data.value?.songRequestOverlaySettings,
	(settings) => {
		if (!settings) return
		if (settingsSub.operation.value?.variables.apiKey !== apiKey.value) return

		overlaySettings.value = {
			...settings,
			style: normalizeSongRequestOverlayStyle(settings.style),
		}
	}
)

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
	requestPolicy: 'network-only',
	context: computed(() => ({ apiKeyGeneration: apiKeyGeneration.value })),
})

function getChannelResolutionGeneration(): string {
	return apiKey.value ? `${apiKeyGeneration.value}:${apiKey.value}` : ''
}

function isCurrentChannelResolutionOperation(): boolean {
	if (!apiKey.value) return false
	return (
		channelQuery.operation.value?.variables.apiKey === apiKey.value &&
		channelQuery.operation.value.context.apiKeyGeneration === apiKeyGeneration.value
	)
}

const channelResolutionGeneration = computed(getChannelResolutionGeneration)
const isChannelResolutionComplete = computed(() => {
	return (
		isCurrentChannelResolutionOperation() &&
		!channelQuery.error.value &&
		channelQuery.data.value !== undefined
	)
})
const channelId = computed(() => {
	if (!isChannelResolutionComplete.value) return ''
	return channelQuery.data.value?.channelByApiKey?.id ?? ''
})

let channelResolutionRetryTimer: ReturnType<typeof setTimeout> | null = null
let channelResolutionRetryAttempt = 0

function resetChannelResolutionRetry() {
	if (channelResolutionRetryTimer) {
		clearTimeout(channelResolutionRetryTimer)
		channelResolutionRetryTimer = null
	}
	channelResolutionRetryAttempt = 0
}

function scheduleChannelResolutionRetry(generation: string) {
	if (
		!generation ||
		channelResolutionRetryTimer ||
		channelQuery.fetching.value ||
		isChannelResolutionComplete.value
	) {
		return
	}

	const delay = Math.min(
		CHANNEL_RESOLUTION_RETRY_BASE_DELAY * 2 ** channelResolutionRetryAttempt,
		CHANNEL_RESOLUTION_RETRY_MAX_DELAY
	)
	channelResolutionRetryAttempt = Math.min(channelResolutionRetryAttempt + 1, 4)

	channelResolutionRetryTimer = setTimeout(() => {
		channelResolutionRetryTimer = null
		if (isUnmounted || generation !== channelResolutionGeneration.value) return
		if (channelQuery.fetching.value || isChannelResolutionComplete.value) return

		channelQuery.executeQuery({ requestPolicy: 'network-only' })
	}, delay)
}

watch(channelResolutionGeneration, resetChannelResolutionRetry, { flush: 'sync' })

watch(
	[
		isChannelResolutionComplete,
		() => channelQuery.fetching.value,
		() => channelQuery.error.value,
		() => channelQuery.operation.value,
	],
	([resolutionComplete, fetching]) => {
		if (resolutionComplete) {
			resetChannelResolutionRetry()
			return
		}
		if (fetching || !isCurrentChannelResolutionOperation()) return

		scheduleChannelResolutionRetry(channelResolutionGeneration.value)
	},
	{ flush: 'sync' }
)

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
	requestPolicy: 'network-only',
	context: computed(() => ({
		apiKeyGeneration: apiKeyGeneration.value,
		connectionGeneration: connectionGeneration.value,
	})),
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
	context: computed(() => ({
		apiKeyGeneration: apiKeyGeneration.value,
		connectionGeneration: connectionGeneration.value,
	})),
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
	context: computed(() => ({
		apiKeyGeneration: apiKeyGeneration.value,
		connectionGeneration: connectionGeneration.value,
	})),
})

function getWidgetSnapshotGeneration(): string {
	const currentChannelId = channelId.value
	if (!apiKey.value || !currentChannelId) return ''

	return `${apiKeyGeneration.value}:${connectionGeneration.value}:${apiKey.value}:${currentChannelId}`
}

function isCurrentWidgetSnapshotOperation(): boolean {
	const currentChannelId = channelId.value
	if (!currentChannelId) return false
	if (widgetDataQuery.operation.value?.variables.channelId !== currentChannelId) return false
	if (widgetDataQuery.operation.value.context.apiKeyGeneration !== apiKeyGeneration.value) {
		return false
	}
	if (widgetDataQuery.operation.value.context.connectionGeneration !== connectionGeneration.value) {
		return false
	}
	return true
}

const widgetSnapshotGeneration = computed(getWidgetSnapshotGeneration)

const currentWidgetData = computed(() => {
	if (!isCurrentWidgetSnapshotOperation() || widgetDataQuery.error.value) return null

	return widgetDataQuery.data.value?.songRequestWidgetData ?? null
})

const isInitialPlaybackReady = computed(() => currentWidgetData.value !== null)

let widgetSnapshotRetryTimer: ReturnType<typeof setTimeout> | null = null
let widgetSnapshotRetryAttempt = 0

function resetWidgetSnapshotRetry() {
	if (widgetSnapshotRetryTimer) {
		clearTimeout(widgetSnapshotRetryTimer)
		widgetSnapshotRetryTimer = null
	}
	widgetSnapshotRetryAttempt = 0
}

function scheduleWidgetSnapshotRetry(generation: string) {
	if (
		!generation ||
		widgetSnapshotRetryTimer ||
		widgetDataQuery.fetching.value ||
		currentWidgetData.value
	) {
		return
	}

	const delay = Math.min(
		WIDGET_SNAPSHOT_RETRY_BASE_DELAY * 2 ** widgetSnapshotRetryAttempt,
		WIDGET_SNAPSHOT_RETRY_MAX_DELAY
	)
	widgetSnapshotRetryAttempt = Math.min(widgetSnapshotRetryAttempt + 1, 4)

	widgetSnapshotRetryTimer = setTimeout(() => {
		widgetSnapshotRetryTimer = null
		if (isUnmounted || generation !== widgetSnapshotGeneration.value) return
		if (widgetDataQuery.fetching.value || currentWidgetData.value) return

		widgetDataQuery.executeQuery({ requestPolicy: 'network-only' })
	}, delay)
}

watch(widgetSnapshotGeneration, resetWidgetSnapshotRetry, { flush: 'sync' })

watch(
	[
		currentWidgetData,
		() => widgetDataQuery.fetching.value,
		() => widgetDataQuery.error.value,
		() => widgetDataQuery.operation.value,
	],
	([snapshot, fetching]) => {
		if (snapshot) {
			resetWidgetSnapshotRetry()
			return
		}
		if (fetching || !isCurrentWidgetSnapshotOperation()) return

		scheduleWidgetSnapshotRetry(widgetSnapshotGeneration.value)
	},
	{ flush: 'sync' }
)

const playbackState = computed(() => {
	const currentChannelId = channelId.value
	if (!currentChannelId) return null

	if (
		playbackSub.operation.value?.variables.channelId === currentChannelId &&
		playbackSub.operation.value.context.apiKeyGeneration === apiKeyGeneration.value &&
		playbackSub.operation.value.context.connectionGeneration === connectionGeneration.value
	) {
		const subscribedState = playbackSub.data.value?.songRequestPlaybackState
		if (subscribedState) return subscribedState
	}

	return currentWidgetData.value?.playbackState ?? null
})

const queue = computed(() => {
	const currentChannelId = channelId.value
	if (!currentChannelId) return []

	if (
		queueSub.operation.value?.variables.channelId === currentChannelId &&
		queueSub.operation.value.context.apiKeyGeneration === apiKeyGeneration.value &&
		queueSub.operation.value.context.connectionGeneration === connectionGeneration.value
	) {
		const subscribedQueue = queueSub.data.value?.songRequestQueueUpdated
		if (subscribedQueue) return subscribedQueue
	}

	return currentWidgetData.value?.queue ?? []
})

const currentQueueItem = computed(() => {
	if (!playbackState.value?.videoId) return null

	return queue.value.find((item) => item.id === playbackState.value?.videoId) ?? null
})

const requesterName = computed(() => {
	return currentQueueItem.value?.orderedByDisplayName || currentQueueItem.value?.orderedByName || ''
})

const { executeMutation: playMutation } = useMutation(
	graphql(`
		mutation OverlaySongRequestPlay($channelId: UUID!, $videoId: String!) {
			songRequestPlay(channelId: $channelId, videoId: $videoId)
		}
	`)
)

const { executeMutation: skipMutation } = useMutation(
	graphql(`
		mutation OverlaySongRequestSkip($channelId: UUID!) {
			songRequestSkip(channelId: $channelId)
		}
	`)
)

const { executeMutation: updatePositionMutation } = useMutation(
	graphql(`
		mutation OverlaySongRequestUpdatePosition($channelId: UUID!, $position: Float!) {
			songRequestUpdatePosition(channelId: $channelId, position: $position)
		}
	`)
)

const playerReady = ref(false)
const playerIsPlaying = ref(false)
const duration = ref(0)
const currentPosition = ref(0)
const currentVideoId = ref('')
const lastAppliedIsPlaying = ref<boolean | null>(null)
const skippedVideoId = ref('')
const displayedDuration = computed(
	() => duration.value || currentQueueItem.value?.durationSeconds || 0
)

let player: YTPlayer | null = null
let tickInterval: ReturnType<typeof setInterval> | null = null
let playerScope: ChannelScope | null = null
let youtubeReadyCallback: (() => void) | null = null
let playerLoadGeneration = 0
let activePlayerLoadGeneration = -1
let endedPlayerLoadGeneration = -1

function getCurrentChannelScope(): ChannelScope | null {
	const currentApiKey = apiKey.value
	const currentChannelId = channelId.value
	if (!currentApiKey || !currentChannelId) return null

	return {
		apiKey: currentApiKey,
		channelId: currentChannelId,
		apiKeyGeneration: apiKeyGeneration.value,
		connectionGeneration: connectionGeneration.value,
	}
}

function isCurrentChannelScope(scope: ChannelScope): boolean {
	return (
		scope.apiKeyGeneration === apiKeyGeneration.value &&
		scope.connectionGeneration === connectionGeneration.value &&
		scope.apiKey === apiKey.value &&
		scope.channelId === channelId.value
	)
}

function resetChannelPlaybackState() {
	apiKeyGeneration.value += 1
	playerLoadGeneration += 1
	activePlayerLoadGeneration = -1
	endedPlayerLoadGeneration = -1
	overlaySettings.value = { ...DEFAULT_OVERLAY_SETTINGS }
	playerScope = null
	playerIsPlaying.value = false
	duration.value = 0
	currentPosition.value = 0
	currentVideoId.value = ''
	lastAppliedIsPlaying.value = null
	skippedVideoId.value = ''

	if (player && playerReady.value) {
		player.stopVideo()
	}
}

watch(apiKey, resetChannelPlaybackState, { flush: 'sync' })

function getWindowWithYouTube(): WindowWithYouTube {
	return window as WindowWithYouTube
}

function refreshPlayerDuration() {
	const playerDuration = player?.getDuration() ?? 0
	if (Number.isFinite(playerDuration) && playerDuration > 0) {
		duration.value = playerDuration
	}
}

function getLoadedPlayerVideoId(): string {
	const videoUrl = player?.getVideoUrl()
	if (!videoUrl) return ''

	try {
		return new URL(videoUrl).searchParams.get('v') ?? ''
	} catch {
		return ''
	}
}

function isExpectedPlayerLoad(scope: ChannelScope | null = playerScope): scope is ChannelScope {
	return (
		!!player &&
		!!scope &&
		!!currentVideoId.value &&
		isCurrentChannelScope(scope) &&
		getLoadedPlayerVideoId() === currentVideoId.value
	)
}

function isActivePlayerLoad(scope: ChannelScope | null = playerScope): scope is ChannelScope {
	return (
		isExpectedPlayerLoad(scope) &&
		activePlayerLoadGeneration === playerLoadGeneration &&
		endedPlayerLoadGeneration !== playerLoadGeneration
	)
}

function handlePlayerEnded(win: WindowWithYouTube) {
	const scope = playerScope
	const videoId = currentVideoId.value
	const loadGeneration = playerLoadGeneration
	if (!isActivePlayerLoad(scope) || endedPlayerLoadGeneration === loadGeneration) {
		return
	}

	queueMicrotask(() => {
		if (isUnmounted || !player || !win.YT) return
		if (loadGeneration !== playerLoadGeneration) return
		if (endedPlayerLoadGeneration === loadGeneration) return
		if (!isCurrentChannelScope(scope)) return
		if (player.getPlayerState() !== win.YT.PlayerState.ENDED) return
		if (getLoadedPlayerVideoId() !== videoId) return

		playerIsPlaying.value = false
		endedPlayerLoadGeneration = loadGeneration
		skippedVideoId.value = videoId
		void skipMutation({ channelId: scope.channelId })
	})
}

function handlePlayerStateChange(win: WindowWithYouTube, playerState: number) {
	const scope = playerScope
	const videoId = currentVideoId.value
	const loadGeneration = playerLoadGeneration

	queueMicrotask(() => {
		if (isUnmounted || !player || !win.YT || !scope) return
		if (loadGeneration !== playerLoadGeneration) return
		if (currentVideoId.value !== videoId) return
		if (!isCurrentChannelScope(scope)) return
		if (getLoadedPlayerVideoId() !== videoId) return
		if (player.getPlayerState() !== playerState) return

		if (playerState !== win.YT.PlayerState.ENDED) {
			activePlayerLoadGeneration = loadGeneration
		}
		refreshPlayerDuration()

		if (playerState === win.YT.PlayerState.PLAYING) {
			playerIsPlaying.value = true
			return
		}

		if (playerState === win.YT.PlayerState.PAUSED) {
			playerIsPlaying.value = false
			return
		}

		if (playerState === win.YT.PlayerState.ENDED) {
			handlePlayerEnded(win)
		}
	})
}

function createPlayer() {
	const win = getWindowWithYouTube()
	if (isUnmounted || !win.YT?.Player || player) return

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
				if (isUnmounted) return
				playerReady.value = true
			},
			onStateChange: (event) => {
				if (isUnmounted || !win.YT) return
				handlePlayerStateChange(win, event.data)
			},
		},
	})
}

onMounted(() => {
	const win = getWindowWithYouTube()

	if (win.YT?.Player) {
		createPlayer()
	} else {
		youtubeReadyCallback = () => createPlayer()
		win.onYouTubeIframeAPIReady = youtubeReadyCallback

		if (!document.querySelector('script[src="https://www.youtube.com/iframe_api"]')) {
			const tag = document.createElement('script')
			tag.src = 'https://www.youtube.com/iframe_api'
			document.head.appendChild(tag)
		}
	}

	tickInterval = setInterval(() => {
		const scope = playerScope
		if (!playerReady.value || !playerIsPlaying.value || !isActivePlayerLoad(scope)) {
			return
		}

		const position = player.getCurrentTime()
		if (Number.isNaN(position)) return
		if (!isActivePlayerLoad(scope)) return

		currentPosition.value = position
		void updatePositionMutation({ channelId: scope.channelId, position })
	}, 1000)
})

onUnmounted(() => {
	isUnmounted = true
	resetChannelResolutionRetry()
	resetWidgetSnapshotRetry()
	playerLoadGeneration += 1
	playerScope = null

	if (tickInterval) {
		clearInterval(tickInterval)
		tickInterval = null
	}

	const win = getWindowWithYouTube()
	if (youtubeReadyCallback && win.onYouTubeIframeAPIReady === youtubeReadyCallback) {
		delete win.onYouTubeIframeAPIReady
	}
	youtubeReadyCallback = null

	player?.destroy()
	player = null
	playerReady.value = false
	playerIsPlaying.value = false
})

watch(
	[playbackState, playerReady, channelId, isInitialPlaybackReady],
	([state, ready, _channelId, initialPlaybackReady]) => {
		if (!ready || !player || !initialPlaybackReady) return
		const scope = getCurrentChannelScope()

		if (!state?.videoId || !scope) {
			player.stopVideo()
			playerLoadGeneration += 1
			activePlayerLoadGeneration = -1
			endedPlayerLoadGeneration = -1
			playerScope = null
			playerIsPlaying.value = false
			currentVideoId.value = ''
			lastAppliedIsPlaying.value = null
			currentPosition.value = 0
			duration.value = 0
			return
		}

		const videoOrScopeChanged =
			state.videoId !== currentVideoId.value || !playerScope || !isCurrentChannelScope(playerScope)

		if (videoOrScopeChanged) {
			playerLoadGeneration += 1
			activePlayerLoadGeneration = -1
			endedPlayerLoadGeneration = -1
			playerScope = scope
			currentVideoId.value = state.videoId
			skippedVideoId.value = ''
			lastAppliedIsPlaying.value = null
			currentPosition.value = state.position
			duration.value = currentQueueItem.value?.durationSeconds ?? 0
			player.loadVideoById(state.videoId, state.position)
			refreshPlayerDuration()
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
	{ deep: true, immediate: true }
)

watch(
	[queue, playbackState, playerReady, isInitialPlaybackReady],
	([items, state, ready, initialPlaybackReady]) => {
		const scope = getCurrentChannelScope()
		if (!ready || !scope || !initialPlaybackReady || state?.videoId) return

		const nextVideoId = items[0]?.id ?? ''
		if (!nextVideoId || nextVideoId === skippedVideoId.value) return
		if (!isCurrentChannelScope(scope)) return

		void playMutation({ channelId: scope.channelId, videoId: nextVideoId })
	},
	{ immediate: true }
)

const isVisible = computed(() => {
	if (!playbackState.value?.videoId) return false
	return !overlaySettings.value.hideOnPause || playbackState.value.isPlaying
})
</script>

<template>
	<div
		class="song-request-overlay"
		:class="{ 'song-request-overlay--hidden': !isVisible }"
	>
		<SongRequestOverlayRenderer
			:style="overlaySettings.style"
			:title="playbackState?.title"
			:requester="requesterName"
			:video-id="playbackState?.videoId"
			:position="currentPosition"
			:duration="displayedDuration"
			:is-playing="playbackState?.isPlaying"
			:accent-color="overlaySettings.accentColor"
			:ticker-background-color="overlaySettings.tickerBackgroundColor"
			:ticker-text-color="overlaySettings.tickerTextColor"
			:ticker-speed="overlaySettings.tickerSpeed"
		>
			<template #media>
				<div class="player-container">
					<div
						id="yt-player"
						class="yt-player"
					/>
				</div>
			</template>
		</SongRequestOverlayRenderer>
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
	width: 100%;
	height: 100%;
	background: transparent;
}

.song-request-overlay .yt-player {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
}
</style>
