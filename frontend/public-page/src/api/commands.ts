import { useQuery } from '@tanstack/vue-query';
import { type MaybeRef, unref } from 'vue';

import { unprotectedClient } from './twirp.js';

export const useCommands = (channelId?: MaybeRef<string | null>) => {
	return useQuery({
		queryKey: ['commands', channelId],
		queryFn: async () => {
			const id = unref(channelId) as string;
			if (!id) return { commands: [] };

			const call = await unprotectedClient.getChannelCommands({
				channelId: id,
			});

			return call.response;
		},
	});
};
