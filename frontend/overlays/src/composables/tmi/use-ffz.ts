import { defineStore } from 'pinia';

import { useEmotes } from './use-emotes.js';

import { requestWithOutCache } from '@/helpers.js';
import type { FfzChannelResponse, FfzGlobalResponse } from '@/types.js';

export const useFrankerFaceZ = defineStore('ffz', () => {
	const { setFrankerFaceZEmotes } = useEmotes();

	async function fetchFrankerFaceZEmotes(channelId: string): Promise<void> {
		try {
			const [globalEmotes, channelEmotes] = await Promise.all([
				requestWithOutCache<FfzGlobalResponse>(
					'https://api.frankerfacez.com/v1/set/global',
				),
				requestWithOutCache<FfzChannelResponse>(
					`https://api.frankerfacez.com/v1/room/id/${channelId}`,
				),
			]);

			setFrankerFaceZEmotes(globalEmotes);
			setFrankerFaceZEmotes(channelEmotes);
		} catch (err) {
			console.error(err);
		}
	}

	return {
		fetchFrankerFaceZEmotes,
	};
});
