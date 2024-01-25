import { useQuery } from '@tanstack/vue-query';
import { unref, type MaybeRef } from 'vue';

import { unprotectedClient } from '@/api/twirp.js';

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
