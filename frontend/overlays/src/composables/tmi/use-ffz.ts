import { defineStore } from 'pinia';

import { useEmotes } from './use-emotes.js';

import type { FfzChannelResponse, FfzGlobalResponse } from '@/types.js';

export const useFrankerFaceZ = defineStore('ffz', () => {
	const { setFrankerFaceZEmotes } = useEmotes();

	async function fetchFrankerFaceZEmotes(channelId: string): Promise<void> {
		try {
			const [global, channel] = await Promise.all([
				fetch('https://api.frankerfacez.com/v1/set/global'),
				fetch(`https://api.frankerfacez.com/v1/room/id/${channelId}?ts=${Date.now()}`),
			]);

			setFrankerFaceZEmotes((await global.json()) as FfzGlobalResponse);
			setFrankerFaceZEmotes((await channel.json()) as FfzChannelResponse);
		} catch (err) {
			console.error(err);
		}
	}

	return {
		fetchFrankerFaceZEmotes,
	};
});
