import { useQuery } from '@tanstack/vue-query';

import { protectedApiClient } from './twirp';

export const useDashboardStats = () => useQuery({
	queryKey: ['dashboardStats'],
	queryFn: async () => {
		const call = await protectedApiClient.getDashboardStats({});

		return call.response;
	},
});
