import { useQuery } from '@tanstack/vue-query'

import { unprotectedApiClient } from '../twirp.js'

export function useStreamers() {
	return useQuery({
		queryKey: ['streamers'],
		queryFn: async () => {
			const req = await unprotectedApiClient.getStatsTwirStreamers({})
			return req.response
		}
	})
}
