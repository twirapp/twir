import type { GetTwirStreamersResponse_Streamer } from '@twir/api/messages/stats/stats';

import { unProtectedClient } from '@/api/twirp.js';

export async function getStreamers() {
	const streamers: GetTwirStreamersResponse_Streamer[][] = [];
	try {
		const streamersResponse = await unProtectedClient.getStatsTwirStreamers({});
		const sortedStreamers = streamersResponse.response.streamers.sort((a, b) => b.followersCount - a.followersCount);

		if (import.meta.env.DEV) {
			if (sortedStreamers.length) {
				streamers.push(...chunk(Array.from({ length: 100 }).map(() => sortedStreamers.at(0)!), 3));
			}
		} else {
			streamers.push(...chunk(sortedStreamers, 3));
		}

		return streamers;
	} catch (err) {
		console.error(err);
		return streamers;
	}
}

function chunk<T>(arr: T[], size: number): T[][] {
	const result: T[][] = [];

	for (let i = 0; i < arr.length; i += size) {
		result.push(arr.slice(i, i + size));
	}

	return result;
}
