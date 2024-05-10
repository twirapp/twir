import { useQuery } from '@tanstack/vue-query'

import { unprotectedClient } from '@/api/twirp.js'
import { useStreamerProfile } from '@/api/use-streamer-profile'

export function useSongsQueue() {
	const { data: profile } = useStreamerProfile()

	return useQuery({
		queryKey: ['songsQueue', profile],
		queryFn: async () => {
			const call = await unprotectedClient.getSongsQueue({
				channelId: profile.value!.twitchGetUserByName!.id,
			})

			return call.response
		},
		refetchInterval: 1000,
		enabled: () => !!profile.value,
	})
}
