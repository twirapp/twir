import { useQuery } from '@tanstack/vue-query';
import { computed } from 'vue';

import { unprotectedClient } from '@/api/twirp.js';
import { useStreamerProfile } from '@/composables/use-streamer-profile';

export const useSongsQueue = () => {
	const { data: profile } = useStreamerProfile();

	return useQuery({
		queryKey: ['songsQueue', profile.value?.id],
		queryFn: async () => {
			const call = await unprotectedClient.getSongsQueue({
				channelId: profile.value!.id,
			});

			return call.response;
		},
		refetchInterval: 1000,
		enabled: computed(() => !!profile.value),
	});
};
