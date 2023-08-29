import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { useQuery } from '@tanstack/vue-query';
import { type TwitchGetUsersResponse } from '@twir/grpc/generated/api/api/twitch';
import { UnProtectedClient } from '@twir/grpc/generated/api/api.client';
import { type ComputedRef, MaybeRef, type Ref, unref } from 'vue';

const transport = new TwirpFetchTransport({
	baseUrl: `${window.location.origin}/api/v1`,
	sendJson: import.meta.env.DEV,
});

const unprotectedClient = new UnProtectedClient(transport);

export const useProfile = (userName: string | Ref<string>) => useQuery({
	queryKey: ['profile', userName],
	queryFn: async () => {
		const name = unref(userName);
		if (!name) return;
		const call = await unprotectedClient.twitchGetUsers({
			names: [name],
			ids: [],
		});

		const user = call.response.users[0];
		return user;
	},
});

export const useCommands = (channelId: ComputedRef<string | null>) => {
	return useQuery({
		queryKey: ['commands', channelId],
		queryFn: async () => {
			const id = unref(channelId) as string;
			if (!id) return;

			const call = await unprotectedClient.getChannelCommands({
				channelId: id,
			});

			return call.response;
		},
	});
};

export const useSongsQueue = (channelId: ComputedRef<string | null>) => useQuery({
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
