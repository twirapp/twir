import { useQuery } from '@tanstack/vue-query';

import { unprotectedApiClient } from './twirp';

export const useStreamers = () => useQuery({
	queryKey: ['streamers'],
	queryFn: async () => {
		const req = await unprotectedApiClient.getStatsTwirStreamers({});
		return req.response;
	},
});
