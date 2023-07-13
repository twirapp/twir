import { useQuery } from '@tanstack/vue-query';

import { unprotectedApiClient } from '@/api/twirp.js';

export const useTwitchSearchUsers = (opts: { ids?: string[], names?: string[] }) => useQuery({
	queryKey: ['twitch', 'search', 'users', opts.ids, opts.names],
	queryFn: async () => {
		console.dir(opts);
		const call = await unprotectedApiClient.twitchSearchUsers({
			ids: opts.ids ?? [],
			names: opts.names ?? [],
		});

		return call.response;
	},
});
