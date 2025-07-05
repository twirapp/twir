import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp.js';

export const useDonatepayIntegration = () => {
	const queryClient = useQueryClient();
	const queryKey = ['donatepay'];

	return {
		useGetData: () => useQuery({
			queryKey,
			queryFn: async () => {
				const call = await protectedApiClient.integrationsDonatepayGet({});
				return call.response;
			},
		}),
		usePost: () => useMutation({
			mutationKey: ['donatepay/post'],
			mutationFn: async (apiKey: string) => {
				await protectedApiClient.integrationsDonatepayPut({ apiKey });
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(queryKey);
			},
		}),
	};
};
