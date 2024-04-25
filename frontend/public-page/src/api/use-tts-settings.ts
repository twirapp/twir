import { useQuery } from '@tanstack/vue-query';
import { storeToRefs } from 'pinia';

import { unprotectedClient } from '@/api/twirp.js';
import { useStreamerProfile } from '@/api/use-streamer-profile';

export const useTTSChannelSettings = () => {
	const { data: profile } = storeToRefs(useStreamerProfile());

	return useQuery({
		queryKey: ['channelTTSSettings', profile],
		queryFn: async () => {
			const call = await unprotectedClient.getTTSChannelSettings({
				channelId: profile.value!.twitchGetUserByName!.id,
			});

			return call.response;
		},
		enabled: () => !!profile.value,
	});
};

export const useTTSUsersSettings = () => {
	const { data: profile } = storeToRefs(useStreamerProfile());

	return useQuery({
		queryKey: ['usersTTSSettings', profile],
		queryFn: async () => {
			const call = await unprotectedClient.getTTSUsersSettings({
				channelId: profile.value!.twitchGetUserByName!.id,
			});

			return call.response;
		},
		enabled: () => !!profile.value,
	});
};
