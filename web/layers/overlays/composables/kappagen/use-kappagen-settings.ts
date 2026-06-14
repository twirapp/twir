import { createGlobalState } from '@vueuse/core'
import { readonly, ref } from 'vue'

import type { KappagenEmojiStyle } from '~~/app/gql/graphql'
import type { KappagenConfig } from '@twirapp/kappagen/types'

type Settings = KappagenConfig & { emojiStyle?: KappagenEmojiStyle, excludedEmotes?: string[] }

export const useKappagenSettings = createGlobalState(() => {
	const settings = ref<Settings>({})

	function setSettings(newSettings: Settings) {
		settings.value = newSettings
	}

	return {
		settings: readonly(settings),
		setSettings,
	}
})
