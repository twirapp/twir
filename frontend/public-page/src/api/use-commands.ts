import { useQuery } from '@tanstack/vue-query';
import { storeToRefs } from 'pinia';

import { unprotectedClient } from '@/api/twirp.js';
import { useStreamerProfile } from '@/api/use-streamer-profile';

export const useCommands = () => {
	const { data: profile } = storeToRefs(useStreamerProfile());

	return useQuery({
		queryKey: ['commands', profile],
		queryFn: async () => {
			const id = profile.value?.twitchGetUserByName?.id;
			if (!id) return { commands: [] };

			const call = await unprotectedClient.getChannelCommands({
				channelId: id,
			});

			return call.response;
		},
	});
};
