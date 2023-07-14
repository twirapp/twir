import { useQuery } from '@tanstack/vue-query';
import { TwitchGetUsersResponse } from '@twir/grpc/generated/api/api/twitch';
import { Ref, isRef } from 'vue';

import { unprotectedApiClient } from '@/api/twirp.js';

export const useTwitchGetUsers = (opts: {
	ids?: Ref<string[]> | Ref<string> | string[],
	names?: Ref<string[]> | Ref<string> | string[]
}) => useQuery({
	queryKey: ['twitch', 'search', 'users', opts.ids, opts.names],
	queryFn: async () => {
		const ids = isRef(opts?.ids)
			? Array.isArray(opts.ids.value) ? opts.ids.value : [opts.ids.value]
			: opts?.ids ?? [''];
		const names = isRef(opts?.names)
			? Array.isArray(opts.names.value) ? opts.names.value : [opts.names.value]
			: opts?.names ?? [''];

		if (ids.length === 0 && names.length === 0) {
			return {
				users: [],
			} as TwitchGetUsersResponse;
		}

		const call = await unprotectedApiClient.twitchGetUsers({
			ids,
			names,
		});

		return call.response;
	},
});

export const useTwitchSearchChannels = (query: string | Ref<string>) => useQuery({
	queryKey: ['twitch', 'search', 'channels', query],
	queryFn: async () => {
		const rawQuery = isRef(query) ? query.value : query;

		if (!rawQuery) return {
			channels: [],
		};

		const call = await unprotectedApiClient.twitchSearchChannels({
			query: rawQuery,
		});

		return call.response;
	},
});
