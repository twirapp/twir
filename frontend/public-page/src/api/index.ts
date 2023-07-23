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
		const call = await unprotectedClient.twitchGetUsers({
			names: [unref(userName)],
			ids: [],
		});

		const user = call.response.users[0];
		return user;
	},
	enabled: !!unref(userName),
});

export const useCommands = (channelId: ComputedRef<string | null>) => {
	return useQuery({
		queryKey: ['commands', channelId],
		queryFn: async () => {
			const call = await unprotectedClient.getChannelCommands({
				channelId: unref(channelId) as string,
			});

			return call.response;
		},
		enabled: !channelId.value,
	});
};

export const useSongsQueue = (channelId: ComputedRef<string | null>) => useQuery({
	queryKey: ['songsQueue', channelId],
	queryFn: async () => {
		const call = await unprotectedClient.getSongsQueue({
			channelId: unref(channelId) as string,
		});

		return call.response;
	},
	enabled: !channelId.value,
});
