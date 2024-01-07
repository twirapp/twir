import { useQuery } from '@tanstack/vue-query';
import type { GetPublicUserInfoResponse } from '@twir/grpc/generated/api/api/auth';
import type { TwitchUser } from '@twir/grpc/generated/api/api/twitch';
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
