import { useQuery } from '@tanstack/vue-query';
import { unref, type MaybeRef } from 'vue';

import { unprotectedClient } from './twirp.js';

export const useSongsQueue = (channelId: MaybeRef<string | null>) => useQuery({
	queryKey: ['songsQueue', channelId],
	queryFn: async () => {
		const id = unref(channelId) as string;
		if (!id) return;

		const call = await unprotectedClient.getSongsQueue({
			channelId: id,
		});

		return call.response;
	},
	refetchInterval: 1000,
});
