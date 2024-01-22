import { useQuery } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp';

export const useSevenTvIntegration = () => {
	return {
		useData: () => useQuery({
			queryKey: ['sevenTvIntegration'],
			refetchInterval: 5000,
			queryFn: async () => {
				const request = await protectedApiClient.integrationsSevenTvGetData({});
				return request.response;
			},
		}),
	};
};
