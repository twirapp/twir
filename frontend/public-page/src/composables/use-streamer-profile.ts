import { useQuery } from '@tanstack/vue-query';
import type { GetPublicUserInfoResponse } from '@twir/api/messages/auth/auth';
import type { TwitchUser } from '@twir/api/messages/twitch/twitch';
import { type MaybeRef, unref } from 'vue';

import { unprotectedClient } from '@/api/twirp';

export type StreamerProfile = TwitchUser & GetPublicUserInfoResponse;

export const useStreamerProfile = (userName?: MaybeRef<string | null>) => {
	return useQuery({
		queryKey: ['streamerProfile'],
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
			} as StreamerProfile;
		},
	});
};
