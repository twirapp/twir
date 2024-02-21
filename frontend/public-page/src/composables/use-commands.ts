import { useQuery } from '@tanstack/vue-query';
import { unref } from 'vue';

import { unprotectedClient } from '@/api/twirp.js';
import { useStreamerProfile } from '@/composables/use-streamer-profile';

export const useCommands = () => {
	const { data: profile } = useStreamerProfile();

	return useQuery({
		queryKey: ['commands', profile.value?.id],
		queryFn: async () => {
			const id = unref(profile.value!.id) as string;
			if (!id) return { commands: [] };

			const call = await unprotectedClient.getChannelCommands({
				channelId: id,
			});

			return call.response;
		},
		enabled: () => !!profile.value,
	});
};
