import { createGlobalState, useLocalStorage, useScriptTag } from '@vueuse/core'
import { computed, nextTick, ref, watch } from 'vue'

import { useYoutubeSocket } from '@/components/songRequests/hook.js'

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

/**
 * Global YouTube player composable.
 * This manages a single YouTube iframe player instance that persists across route changes.
 * The iframe is always in the DOM and never destroyed.
 */
export const useGlobalYoutubePlayer = createGlobalState(() => {
	const { currentVideo, nextVideo, sendPlaying } = useYoutubeSocket()

	// Player instance and state
	const player = ref<YTPlayer>()
	const playerReady = ref(false)
	const isPlaying = ref(false)
	const sliderTime = ref(0)
	const duration = ref(0)
	const updateTimeInterval = ref<number>()
	const shouldAutoplayNext = ref(false)

	// Settings
	const volume = useLocalStorage('twirPlayerVolume', 10)
	const isMuted = useLocalStorage('twirPlayerIsMuted', false)
	const playerVisible = useLocalStorage('twirPlayerVisible', true)
	const hasEverPlayedSong = useLocalStorage('twirHasEverPlayedSong', false)
	const noCookie = ref(false)

	const { load: loadYTScript } = useScriptTag('https://www.youtube.com/iframe_api', () => {}, {
		manual: true,
	})

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

	function playNext() {
		nextVideo()
	}

	function onPlayerReady(event: { target: YTPlayer }) {
		player.value = event.target
		playerReady.value = true

		console.log('[YouTube Player] Player ready')

		const vol = volume.value
		event.target.setVolume(vol)
		if (isMuted.value) {
			event.target.mute()
		}

		// Cue current video if it exists (without autoplay)
		if (currentVideo.value) {
			event.target.cueVideoById(currentVideo.value.videoId)
		}
	}

	function onPlayerStateChange(event: { target: YTPlayer; data: number }) {
		const stateNames: Record<number, string> = {
			[PlayerState.ENDED]: 'ENDED',
			[PlayerState.PLAYING]: 'PLAYING',
			[PlayerState.PAUSED]: 'PAUSED',
			[PlayerState.BUFFERING]: 'BUFFERING',
		}
		console.log('[YouTube Player] State change:', stateNames[event.data] || event.data)

		if (event.data === PlayerState.PLAYING) {
			console.log('[YouTube Player] ▶️ Playing')
			isPlaying.value = true
			hasEverPlayedSong.value = true
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
			console.log('[YouTube Player] ⏸️ Paused')
			isPlaying.value = false
			stopTimeUpdate()
		} else if (event.data === PlayerState.ENDED) {
			console.log('[YouTube Player] ⏹️ Ended')
			isPlaying.value = false
			stopTimeUpdate()
			shouldAutoplayNext.value = true
			playNext()
		} else if (event.data === PlayerState.BUFFERING) {
			console.log('[YouTube Player] ⏳ Buffering')
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
		console.error('[YouTube Player] Error:', event.data)
		isPlaying.value = false
		stopTimeUpdate()
	}

	async function initPlayer() {
		if (playerReady.value) {
			console.log('[YouTube Player] Already initialized')
			return
		}

		console.log('[YouTube Player] Initializing...')

		const initPlayerInternal = () => {
			if (!window.YT || !window.YT.Player) {
				console.warn('[YouTube Player] YT not available')
				return
			}

			const container = document.getElementById('global-yt-player-container')
			if (!container) {
				console.error('[YouTube Player] Container not found')
				return
			}

			console.log('[YouTube Player] Creating player instance')

			const playerConfig: any = {
				height: '100%',
				width: '100%',
				host: noCookie.value ? 'https://www.youtube-nocookie.com' : undefined,
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
			new window.YT.Player('global-yt-player-container', playerConfig)
		}

		if (window.YT && window.YT.Player) {
			await nextTick()
			initPlayerInternal()
		} else {
			window.onYouTubeIframeAPIReady = () => {
				nextTick(() => {
					initPlayerInternal()
				})
			}
			await loadYTScript()
		}
	}

	function play() {
		if (!player.value || !playerReady.value) {
			console.warn('[YouTube Player] Cannot play: player not ready')
			return
		}
		player.value.playVideo()
	}

	function pause() {
		if (!player.value || !playerReady.value) {
			console.warn('[YouTube Player] Cannot pause: player not ready')
			return
		}
		player.value.pauseVideo()
	}

	function togglePlay() {
		if (isPlaying.value) {
			pause()
		} else {
			play()
		}
	}

	function seek(seconds: number) {
		if (!player.value || !playerReady.value) return
		try {
			player.value.seekTo(seconds, true)
		} catch (e) {
			console.error('Error seeking:', e)
		}
	}

	function setVolume(value: number) {
		volume.value = value
	}

	function toggleMute() {
		isMuted.value = !isMuted.value
	}

	// Watch current video changes
	watch(currentVideo, (video) => {
		console.log('[YouTube Player] Current video changed:', video?.videoId)

		if (!player.value || !playerReady.value) {
			console.log('[YouTube Player] Player not ready, will cue video when ready')
			return
		}

		if (!video) {
			player.value.stopVideo()
			sliderTime.value = 0
			duration.value = 0
			shouldAutoplayNext.value = false
			return
		}

		// If track ended and we're auto-playing next, use loadVideoById
		// Otherwise use cueVideoById to prevent autoplay
		if (shouldAutoplayNext.value) {
			console.log('[YouTube Player] Auto-playing next video')
			player.value.loadVideoById(video.videoId)
			shouldAutoplayNext.value = false
		} else {
			console.log('[YouTube Player] Cueing video without autoplay')
			player.value.cueVideoById(video.videoId)
		}
	})

	// Watch volume changes
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

	// Watch mute changes
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

	return {
		// Player state
		playerReady,
		isPlaying,
		sliderTime,
		duration,
		hasEverPlayedSong,

		// Settings
		volume,
		isMuted,
		playerVisible,
		sliderVolume,
		noCookie,

		// Player actions
		initPlayer,
		play,
		pause,
		togglePlay,
		seek,
		setVolume,
		toggleMute,
		playNext,
	}
})
