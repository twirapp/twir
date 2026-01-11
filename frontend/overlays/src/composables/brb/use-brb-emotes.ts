import { watch } from 'vue'

import { useBetterTv } from '@/composables/tmi/use-bettertv.js'
import { useFrankerFaceZ } from '@/composables/tmi/use-ffz.js'
import { useSevenTv } from '@/composables/tmi/use-seven-tv.js'

import { useBrbSettings } from './use-brb-settings.js'

export function useBrbEmotes() {
	const { settings } = useBrbSettings()
	const { fetchSevenTvEmotes, destroy: destroySevenTv } = useSevenTv()
	const { fetchBttvEmotes } = useBetterTv()
	const { fetchFrankerFaceZEmotes } = useFrankerFaceZ()

	watch(
		() => settings.value?.channelId,
		(channelId) => {
			if (!channelId) return

			// Fetch all emotes (7TV, BTTV, FFZ)
			fetchSevenTvEmotes(channelId)
			fetchBttvEmotes(channelId)
			fetchFrankerFaceZEmotes(channelId)
		},
		{ immediate: true }
	)

	function destroy() {
		destroySevenTv()
	}

	return {
		destroy,
	}
}
