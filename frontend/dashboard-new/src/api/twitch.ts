import { useQuery } from '@tanstack/vue-query';
import { TwitchSearchUsersResponse } from '@twir/grpc/generated/api/api/twitch';
import { Ref, isRef } from 'vue';

import { unprotectedApiClient } from '@/api/twirp.js';

export const useTwitchSearchUsers = (opts: {
	ids?: Ref<string[]> | Ref<string> | string[],
	names?: Ref<string[]> | Ref<string> | string[]
}) => useQuery({
	queryKey: ['twitch', 'search', 'users', opts.ids, opts.names],
	queryFn: async () => {
		const ids = isRef(opts.ids)
			? Array.isArray(opts.ids?.value) ? opts.ids?.value ?? [] : [opts.ids?.value ?? '']
			: Array.isArray(opts.ids) ? opts.ids ?? [] : [opts.ids ?? ''];

		const names = isRef(opts.names)
			? Array.isArray(opts.names?.value) ? opts.names?.value ?? [] : [opts.names?.value ?? '']
			: Array.isArray(opts.names) ? opts.names ?? [] : [opts.names ?? ''];

		if (ids.length === 0 && names.length === 0) {
			return {
				users: [],
			} as TwitchSearchUsersResponse;
		}

		const call = await unprotectedApiClient.twitchSearchUsers({
			ids,
			names,
		});

		return call.response;
	},
});
