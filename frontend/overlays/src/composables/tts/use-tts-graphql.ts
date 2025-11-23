import { useSubscription } from '@urql/vue'
import { computed, ref, watch } from 'vue'

import type { TTSOnSayFn, TTSOnSkipFn } from '@/types.js'

import { useTTSSettings } from './use-tts-settings.js'

import { graphql } from '@/gql'

type Options = {
	onSay: TTSOnSayFn
	onSkip: TTSOnSkipFn
}

export function useTTSOverlayGraphQL(options: Options) {
	const apiKey = ref<string>('')
	const paused = computed(() => !apiKey.value)
	const { setSettings } = useTTSSettings()

	// Subscribe to settings updates
	const {
		data: settingsData,
		executeSubscription: connectSettings,
		pause: pauseSettings,
	} = useSubscription({
		query: graphql(`
			subscription TTSSettings($apiKey: String!) {
				overlaysTTS(apiKey: $apiKey) {
					id
					enabled
					voice
					disallowedVoices
					pitch
					rate
					volume
					doNotReadTwitchEmotes
					doNotReadEmoji
					doNotReadLinks
					allowUsersChooseVoiceInMainCommand
					maxSymbols
					readChatMessages
					readChatMessagesNicknames
					createdAt
					updatedAt
					channelId
				}
			}
		`),
		get variables() {
			return {
				apiKey: apiKey.value,
			}
		},
		pause: paused,
	})

	// Subscribe to say events
	const {
		data: sayData,
		executeSubscription: connectSay,
		pause: pauseSay,
	} = useSubscription({
		query: graphql(`
			subscription TTSSay($apiKey: String!) {
				overlaysTTSSay(apiKey: $apiKey) {
					text
					voice
					rate
					pitch
					volume
				}
			}
		`),
		get variables() {
			return {
				apiKey: apiKey.value,
			}
		},
		pause: paused,
	})

	// Subscribe to skip events
	const {
		data: skipData,
		executeSubscription: connectSkip,
		pause: pauseSkip,
	} = useSubscription({
		query: graphql(`
			subscription TTSSkip($apiKey: String!) {
				overlaysTTSSkip(apiKey: $apiKey)
			}
		`),
		get variables() {
			return {
				apiKey: apiKey.value,
			}
		},
		pause: paused,
	})

	// Watch for settings updates
	watch(settingsData, (data) => {
		if (!data?.overlaysTTS) return

		const overlay = data.overlaysTTS
		setSettings({
			enabled: overlay.enabled,
			voice: overlay.voice,
			disallowedVoices: overlay.disallowedVoices,
			pitch: overlay.pitch,
			rate: overlay.rate,
			volume: overlay.volume,
			doNotReadTwitchEmotes: overlay.doNotReadTwitchEmotes,
			doNotReadEmoji: overlay.doNotReadEmoji,
			doNotReadLinks: overlay.doNotReadLinks,
			allowUsersChooseVoiceInMainCommand: overlay.allowUsersChooseVoiceInMainCommand,
			maxSymbols: overlay.maxSymbols,
			readChatMessages: overlay.readChatMessages,
			readChatMessagesNicknames: overlay.readChatMessagesNicknames,
		})
	})

	// Watch for say events
	watch(sayData, (data) => {
		if (!data?.overlaysTTSSay) return

		const { text, voice, rate, pitch, volume } = data.overlaysTTSSay
		options.onSay({ text, voice, rate, pitch, volume })
	})

	// Watch for skip events
	watch(skipData, (data) => {
		if (!data?.overlaysTTSSkip) return
		options.onSkip()
	})

	function destroy() {
		pauseSettings()
		pauseSay()
		pauseSkip()
	}

	async function connect(key: string) {
		apiKey.value = key
		connectSettings()
		connectSay()
		connectSkip()
	}

	return {
		connect,
		destroy,
		settingsData,
	}
}

