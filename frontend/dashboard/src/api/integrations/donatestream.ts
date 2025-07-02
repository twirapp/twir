import { useMutation, useQuery } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp.js';

export const useDonateStreamIntegration = () => {
	const queryKey = ['donatestream'];

	return {
		useGetData: () => useQuery({
			queryKey,
			queryFn: async () => {
				const call = await protectedApiClient.integrationsDonateStreamGet({});
				return call.response;
			},
		}),
		usePost: () => useMutation({
			mutationKey: ['donatestream/post'],
			mutationFn: async (secret: string) => {
				await protectedApiClient.integrationsDonateStreamPostSecret({ secret });
			},
		}),
	};
};
