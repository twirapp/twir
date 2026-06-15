import { createGlobalState } from '@vueuse/core'

import { useEmotes } from './use-emotes'

import type { BttvChannelResponse, BttvGlobalResponse } from '~~/layers/overlays/app/types'

import { requestWithOutCache } from '~~/layers/overlays/app/helpers'

export const useBetterTv = createGlobalState(() => {
	const { setBttvEmotes } = useEmotes()

	async function fetchBttvEmotes(channelId: string): Promise<void> {
		try {
			const [globalEmotes, channelEmotes] = await Promise.all([
				requestWithOutCache<BttvGlobalResponse>(
					'https://api.betterttv.net/3/cached/emotes/global',
				),
				requestWithOutCache<BttvChannelResponse>(
					`https://api.betterttv.net/3/cached/users/twitch/${channelId}`,
				),
			])

			setBttvEmotes(globalEmotes)
			setBttvEmotes(channelEmotes)
		} catch (err) {
			console.error(err)
		}
	}

	return {
		fetchBttvEmotes,
	}
})
