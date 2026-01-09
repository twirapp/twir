import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import type { Settings } from '@twir/frontend-faceit-stats'

export const useFaceitStats = createGlobalState(() => {
	const settings = ref<Settings>({
		nickname: '',
		game: 'cs2',
		bgColor: '#1f1f22',
		textColor: '#ffffff',
		borderRadius: 24,
		displayAvarageKdr: false,
		displayWorldRanking: false,
		displayLastTwentyMatches: false,
	})

	function setSettings(data: Settings) {
		settings.value = data
	}

	function copyOverlayLink() {
		const url = new URL(`${window.location.origin}/overlays/faceit-stats`)

		for (const [key, value] of Object.entries(settings.value)) {
			url.searchParams.set(key, value as string)
		}

		navigator.clipboard.writeText(url.toString())
	}

	return {
		settings,
		setSettings,
		copyOverlayLink,
	}
})
