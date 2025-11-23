import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

export interface TTSSettings {
	enabled: boolean
	voice: string
	disallowedVoices: string[]
	pitch: number
	rate: number
	volume: number
	doNotReadTwitchEmotes: boolean
	doNotReadEmoji: boolean
	doNotReadLinks: boolean
	allowUsersChooseVoiceInMainCommand: boolean
	maxSymbols: number
	readChatMessages: boolean
	readChatMessagesNicknames: boolean
}

export type TTSSetSettingsFn = (settings: TTSSettings) => void

export const useTTSSettings = createGlobalState(() => {
	const settings = ref<TTSSettings>()

	const setSettings: TTSSetSettingsFn = (newSettings) => {
		settings.value = newSettings
	}

	return {
		settings,
		setSettings,
	}
})

