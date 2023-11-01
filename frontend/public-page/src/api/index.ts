import { useQuery } from '@tanstack/vue-query';
import { type GetPublicUserInfoResponse } from '@twir/grpc/generated/api/api/auth';
import type {
	TwitchGetUsersResponse,
	TwitchUser as InternalUser,
} from '@twir/grpc/generated/api/api/twitch';
import { MaybeRef, type Ref, unref } from 'vue';

import { unprotectedClient } from './twirp.js';

export * from './community.js';

type User = InternalUser & GetPublicUserInfoResponse;

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

export const useTTSChannelSettings = (channelId: MaybeRef<string>) => useQuery({
	queryKey: ['channelTTSSettings', channelId],
	queryFn: async () => {
		const id = unref(channelId) as string;
		if (!id) return;
		const call = await unprotectedClient.getTTSChannelSettings({
			channelId: id,
		});

		return call.response;
	},
});

export const useTTSUsersSettings = (channelId: MaybeRef<string>) => useQuery({
	queryKey: ['usersTTSSettings', channelId],
	queryFn: async () => {
		const id = unref(channelId) as string;
		if (!id) return;
		const call = await unprotectedClient.getTTSUsersSettings({
			channelId: id,
		});

		return call.response;
	},
});

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
