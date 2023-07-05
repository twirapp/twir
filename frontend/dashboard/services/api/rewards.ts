import { useQuery } from '@tanstack/react-query';

import { protectedApiClient } from '@/services/api/twirp';

export const useRewards = () => useQuery<ReturnType<typeof protectedApiClient.rewardsGet>['response']>({
	queryKey: ['rewards'],
	queryFn: async () => {
		const call = await protectedApiClient.rewardsGet({});
		return call.response;
	},
});
