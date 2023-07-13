import { useQuery } from '@tanstack/vue-query';
import { Ref } from 'vue';

import { unprotectedApiClient } from '@/api/twirp.js';

export const useTwitchSearchUsers = (opts: { ids?: Ref<string[]> | Ref<string>, names?: Ref<string[]> | Ref<string> }) => useQuery({
	queryKey: ['twitch', 'search', 'users', opts.ids, opts.names],
	queryFn: async () => {
		const ids = Array.isArray(opts.ids?.value) ? opts.ids?.value ?? [] : [opts.ids?.value ?? ''];
		const names = Array.isArray(opts.names?.value) ? opts.names?.value ?? [] : [opts.names?.value ?? ''];

		const call = await unprotectedApiClient.twitchSearchUsers({
			ids,
			names,
		});

		return call.response;
	},
});
