import { useQuery } from '@tanstack/vue-query'

import { unprotectedClient } from '@/api/twirp.js'
import { useStreamerProfile } from '@/api/use-streamer-profile'

export function useTTSChannelSettings() {
	const { data: profile } = useStreamerProfile()

	return useQuery({
		queryKey: ['channelTTSSettings', profile],
		queryFn: async () => {
			const call = await unprotectedClient.getTTSChannelSettings({
				channelId: profile.value!.twitchGetUserByName!.id,
			})

			return call.response
		},
		enabled: () => !!profile.value,
	})
}

export function useTTSUsersSettings() {
	const { data: profile } = useStreamerProfile()

	return useQuery({
		queryKey: ['usersTTSSettings', profile],
		queryFn: async () => {
			const call = await unprotectedClient.getTTSUsersSettings({
				channelId: profile.value!.twitchGetUserByName!.id,
			})

			return call.response
		},
		enabled: () => !!profile.value,
	})
}
