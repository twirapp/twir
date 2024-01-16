import { useQuery } from '@tanstack/vue-query';
import type { GetPublicUserInfoResponse } from '@twir/api/messages/auth/auth';
import type { TwitchUser } from '@twir/api/messages/twitch/twitch';
import { type Ref, unref } from 'vue';

import { unprotectedClient } from './twirp.js';

type User = TwitchUser & GetPublicUserInfoResponse;

export const useProfile = (userName: string | Ref<string>) => useQuery<User | undefined>({
	queryKey: ['profile', userName],
	queryFn: async () => {
		const name = unref(userName);
		if (!name) return;

		const twitchUserCall = await unprotectedClient.twitchGetUsers({
			names: [name],
			ids: [],
		});

		const twitchUser = twitchUserCall.response.users[0];
		if (!twitchUser) return;

		const internalUserCall = await unprotectedClient.getPublicUserInfo({
			userId: twitchUser.id,
		});

		const internalUser = internalUserCall.response;
		if (!internalUser) return;

		return {
			...twitchUser,
			...internalUser,
		};
	},
});
