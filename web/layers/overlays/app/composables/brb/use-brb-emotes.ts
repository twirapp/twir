import { watch } from 'vue'

import { useBetterTv } from '~~/layers/overlays/app/composables/tmi/use-bettertv'
import { useFrankerFaceZ } from '~~/layers/overlays/app/composables/tmi/use-ffz'
import { useSevenTv } from '~~/layers/overlays/app/composables/tmi/use-seven-tv'

import { useBrbSettings } from './use-brb-settings'

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
