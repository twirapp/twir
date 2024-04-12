import { useQuery } from '@tanstack/vue-query';
import { useSubscription } from '@urql/vue';
import { computed } from 'vue';

import { protectedApiClient } from './twirp';

import { graphql } from '@/gql';

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


export const useRealtimeDashboardStats = () => {
	const { data, isPaused, fetching } = useSubscription({
		query: graphql(`
			subscription dashboardStats {
				dashboardStats {
					categoryId
					categoryName
					viewers
					startedAt
					title
					chatMessages
					followers
					usedEmotes
					requestedSongs
					subs
				}
			}
		`),
	});

	const stats = computed(() => {
		return data.value?.dashboardStats;
	});

	return { stats, isPaused, fetching };
};
