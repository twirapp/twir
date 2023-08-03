import { useQuery } from '@tanstack/vue-query';

import { protectedApiClient } from './twirp';

export const useDashboardStats = () => useQuery({
	queryKey: ['dashboardStats'],
	queryFn: async () => {
		const call = await protectedApiClient.getDashboardStats({});

		return call.response;
	},
});

export const useDashboardEvents = () => useQuery({
	queryKey: ['dashboardEvents'],
	queryFn: async () => {
		const call = await protectedApiClient.getDashboardEventsList({});
		return call.response;
	},
});
