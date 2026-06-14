import { createGlobalState } from '@vueuse/core'

import { useEmotes } from './use-emotes'

import type { FfzChannelResponse, FfzGlobalResponse } from '~/layers/overlays/types'

import { requestWithOutCache } from '~/layers/overlays/helpers'

export const useFrankerFaceZ = createGlobalState(() => {
	const { setFrankerFaceZEmotes } = useEmotes()

	async function fetchFrankerFaceZEmotes(channelId: string): Promise<void> {
		try {
			const [globalEmotes, channelEmotes] = await Promise.all([
				requestWithOutCache<FfzGlobalResponse>(
					'https://api.frankerfacez.com/v1/set/global',
				),
				requestWithOutCache<FfzChannelResponse>(
					`https://api.frankerfacez.com/v1/room/id/${channelId}`,
				),
			])

			setFrankerFaceZEmotes(globalEmotes)
			setFrankerFaceZEmotes(channelEmotes)
		} catch (err) {
			console.error(err)
		}
	}

	return {
		fetchFrankerFaceZEmotes,
	}
})
