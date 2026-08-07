import type { Ref } from 'vue'

interface TwitchPlayerInstance {
	play: () => void
	pause: () => void
	setMuted: (muted: boolean) => void
	isPaused: () => boolean
	addEventListener: (event: string, callback: () => void) => void
}

interface TwitchPlayerOptions {
	width: string | number
	height: string | number
	channel: string
	parent: string[]
	autoplay: boolean
	muted: boolean
	controls: boolean
}

interface TwitchNamespace {
	Player: new (element: HTMLElement, options: TwitchPlayerOptions) => TwitchPlayerInstance
}

declare global {
	interface Window {
		Twitch?: TwitchNamespace
		define?: unknown
		module?: unknown
		exports?: unknown
	}
}

const TWITCH_EMBED_SCRIPT_URL = 'https://player.twitch.tv/js/embed/v1.js'

// https://dev.twitch.tv/docs/embed/video-and-clips — JavaScript Events
const TWITCH_PLAYER_EVENTS = {
	ready: 'ready',
	pause: 'pause',
	playing: 'playing',
	playbackBlocked: 'playbackBlocked',
} as const

let scriptLoadPromise: Promise<void> | null = null

function loadTwitchEmbedScript(): Promise<void> {
	if (window.Twitch?.Player) return Promise.resolve()
	if (scriptLoadPromise) return scriptLoadPromise

	scriptLoadPromise = new Promise((resolve, reject) => {
		const script = document.createElement('script')
		script.src = TWITCH_EMBED_SCRIPT_URL

		// Monaco's AMD loader exposes global define/module, and the Twitch UMD
		// bundle registers itself as an AMD module instead of setting
		// window.Twitch when they are present.
		const masked = { define: window.define, module: window.module, exports: window.exports }
		Reflect.set(window, 'define', undefined)
		Reflect.set(window, 'module', undefined)
		Reflect.set(window, 'exports', undefined)

		const restoreGlobals = () => {
			Reflect.set(window, 'define', masked.define)
			Reflect.set(window, 'module', masked.module)
			Reflect.set(window, 'exports', masked.exports)
		}

		script.onload = () => {
			restoreGlobals()
			if (window.Twitch?.Player) {
				resolve()
			} else {
				scriptLoadPromise = null
				reject(new Error('Twitch embed script loaded without exposing window.Twitch'))
			}
		}
		script.onerror = () => {
			restoreGlobals()
			scriptLoadPromise = null
			reject(new Error('Failed to load Twitch embed script'))
		}
		document.head.appendChild(script)
	})

	return scriptLoadPromise
}

// Twitch refuses autoplay when its own visibility heuristics consider the player
// hidden (the canvas lives under a CSS zoom transform), so playback is driven
// programmatically: play on ready, resume on pause, replay on tab return.
export function useTwitchEmbedPlayer(
	container: Ref<HTMLElement | undefined>,
	channelLogin: Ref<string | null>,
	autoResume: Ref<boolean>
) {
	let player: TwitchPlayerInstance | null = null
	let disposed = false
	let mountAttempts = 0
	const retryTimeouts: ReturnType<typeof setTimeout>[] = []

	// Twitch refuses to START playback while the player is covered by other
	// elements (elementFromPoint-based gate), but keeps playing once started.
	// The background is lifted above the canvas for (re)starts and dropped back
	// behind the layers on the 'playing' event.
	const gateLifted = ref(true)

	function replayIfAllowed() {
		if (disposed || !autoResume.value) return
		gateLifted.value = true
		player?.play()
	}

	function schedulePlayWithRetries(delays: readonly number[] = [0, 300, 1000, 3000, 8000]) {
		for (const timeout of retryTimeouts) clearTimeout(timeout)
		retryTimeouts.length = 0

		for (const delay of delays) {
			retryTimeouts.push(
				setTimeout(() => {
					if (disposed || !autoResume.value || !player) return
					player.setMuted(true)
					if (player.isPaused()) player.play()
				}, delay)
			)
		}
	}

	async function mount() {
		if (!import.meta.client) return

		const login = channelLogin.value
		const element = container.value
		if (!login || !element) return

		try {
			await loadTwitchEmbedScript()
		} catch {
			// handled by the retry below
		}
		if (disposed) return

		if (!window.Twitch) {
			if (mountAttempts++ < 5) {
				retryTimeouts.push(setTimeout(() => void mount(), 1000))
			}
			return
		}
		mountAttempts = 0

		retryTimeouts.push(
			setTimeout(() => {
				if (disposed || !container.value || !window.Twitch) return

				element.replaceChildren()
				player = new window.Twitch.Player(element, {
					width: '100%',
					height: '100%',
					channel: login,
					parent: [window.location.hostname],
					autoplay: true,
					muted: true,
					controls: false,
				})

				player.addEventListener(TWITCH_PLAYER_EVENTS.ready, () => schedulePlayWithRetries())
				player.addEventListener(TWITCH_PLAYER_EVENTS.playing, () => {
					gateLifted.value = false
				})
				player.addEventListener(TWITCH_PLAYER_EVENTS.pause, replayIfAllowed)
				player.addEventListener(TWITCH_PLAYER_EVENTS.playbackBlocked, replayIfAllowed)

				schedulePlayWithRetries()
			}, 500)
		)
	}

	function unmount() {
		disposed = true
		for (const timeout of retryTimeouts) clearTimeout(timeout)
		retryTimeouts.length = 0
		container.value?.replaceChildren()
		player = null
		document.removeEventListener('visibilitychange', handleVisibilityChange)
	}

	function handleVisibilityChange() {
		if (document.visibilityState === 'visible') replayIfAllowed()
	}

	watch(channelLogin, async (login, previousLogin) => {
		// null happens on platform query refetches; only an actual login change
		// justifies a remount.
		if (!login || login === previousLogin) return

		player?.pause()
		container.value?.replaceChildren()
		player = null
		await nextTick()
		void mount()
	})

	watch(autoResume, (enabled) => {
		if (enabled) replayIfAllowed()
		else player?.pause()
	})

	onMounted(() => {
		document.addEventListener('visibilitychange', handleVisibilityChange)
		void mount()
	})

	onUnmounted(unmount)

	return { gateLifted }
}
