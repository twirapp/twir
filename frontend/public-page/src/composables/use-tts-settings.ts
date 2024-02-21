import { useQuery } from '@tanstack/vue-query';

import { unprotectedClient } from '@/api/twirp.js';
import { useStreamerProfile } from '@/composables/use-streamer-profile';

export const useTTSChannelSettings = () => {
	const { data: profile } = useStreamerProfile();

	return useQuery({
		queryKey: ['channelTTSSettings', profile.value?.id],
		queryFn: async () => {
			const call = await unprotectedClient.getTTSChannelSettings({
				channelId: profile.value!.id,
			});

			return call.response;
		},
		enabled: () => !!profile.value,
	});
};

export const useTTSUsersSettings = () => {
	const { data: profile } = useStreamerProfile();

	return useQuery({
		queryKey: ['usersTTSSettings', profile.value?.id],
		queryFn: async () => {
			const call = await unprotectedClient.getTTSUsersSettings({
				channelId: profile.value!.id,
			});

			return call.response;
		},
		enabled: () => !!profile.value,
	});
};
