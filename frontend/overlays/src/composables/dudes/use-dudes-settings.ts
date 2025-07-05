import { useFontSource } from '@twir/fontsource'
import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import type { ChannelData } from '@/types.js'
import type { DudesSprite, DudesUserSettings } from '@twir/types/overlays'
import type { DudesTypes } from '@twirapp/dudes-vue/types'
import type { DudesIgnoreSettings } from '@/gql/graphql'

export interface DudesOverlaySettings {
	maxOnScreen: number
	defaultSprite: keyof typeof DudesSprite
}

export interface DudesConfig {
	ignore: DudesIgnoreSettings
	dudes: {
		dude: DudesTypes.DudeStyles
		sounds: DudesTypes.DudeSounds
		name: DudesTypes.NameBoxStyles
		message: DudesTypes.MessageBoxStyles
		emote: DudesTypes.EmotesStyles
	}
	overlay: DudesOverlaySettings
}

export const useDudesSettings = createGlobalState(() => {
	const fontSource = useFontSource()
	const dudesSettings = ref<DudesConfig | null>(null)
	const dudesUserSettings = new Map<string, DudesUserSettings & { userDisplayName?: string }>()
	const channelData = ref<ChannelData>()

	function updateSettings(settings: DudesConfig): void {
		dudesSettings.value = settings
	}

	function updateChannelData(data: ChannelData): void {
		channelData.value = data
	}

	async function loadFont(
		fontFamily: string,
		fontWeight: number,
		fontStyle: string
	): Promise<string> {
		try {
			await fontSource.loadFont(fontFamily, fontWeight, fontStyle)

			const fontKey = `${fontFamily}-${fontWeight}-${fontStyle}`
			return fontKey
		} catch (err) {
			console.error(err)
			return 'Arial'
		}
	}

	return {
		channelData,
		updateChannelData,
		dudesSettings,
		dudesUserSettings,
		updateSettings,
		loadFont,
	}
})
