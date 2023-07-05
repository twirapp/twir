import { useQuery } from '@tanstack/react-query';
import { UnwrapPromise } from 'next/dist/lib/coalesced-function';

import { protectedApiClient } from '@/services/api/twirp';

export const useRewards = () => useQuery<UnwrapPromise<ReturnType<typeof protectedApiClient.rewardsGet>['response']>>({
	queryKey: ['rewards'],
	queryFn: async () => {
		const call = await protectedApiClient.rewardsGet({});
		return call.response;
	},
});
