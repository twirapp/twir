import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport';
import { useQuery } from '@tanstack/vue-query';
import { UnProtectedClient } from '@twir/grpc/generated/api/api.client';
import { type ComputedRef, type Ref, unref } from 'vue';

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
