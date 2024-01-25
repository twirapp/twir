import { useQuery } from '@tanstack/vue-query';
import type { TwitchGetUsersResponse } from '@twir/api/messages/twitch/twitch';
import { unref, type MaybeRef } from 'vue';

import { unprotectedClient } from '@/api/twirp.js';

export const useTwitchGetUsers = (usersIds: MaybeRef<string[]>) => useQuery({
	queryKey: ['twitchGetUsersByIds', usersIds],
	queryFn: async () => {
		const ids = unref(usersIds);
		if (!ids.length) return {} as TwitchGetUsersResponse;
		const call = await unprotectedClient.twitchGetUsers({
			ids,
			names: [],
		});

		return call.response;
	},
});
