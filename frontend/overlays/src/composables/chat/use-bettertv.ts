import { defineStore } from 'pinia';

import { useEmotes } from './use-emotes.js';

import { BttvChannelResponse, BttvGlobalResponse } from '@/types.js';

export const useBetterTv = defineStore('bettertv', () => {
	const { setBttvEmotes } = useEmotes();

	async function fetchBttvEmotes(channelId: string): Promise<void> {
		try {
			const [global, channel] = await Promise.all([
				fetch('https://api.betterttv.net/3/cached/emotes/global'),
				fetch(`https://api.betterttv.net/3/cached/users/twitch/${channelId}?ts=${Date.now()}`),
			]);

			setBttvEmotes((await global.json()) as BttvGlobalResponse);
			setBttvEmotes((await channel.json()) as BttvChannelResponse);
		} catch (err) {
			console.error(err);
		}
	}

	return {
		fetchBttvEmotes,
	};
});
