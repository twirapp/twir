import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp.js';

export const useDonateStreamIntegration = () => {
	const queryClient = useQueryClient();
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
			onSuccess: async () => {
				await queryClient.invalidateQueries(queryKey);
			},
		}),
	};
};
