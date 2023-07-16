import { useQueryClient, useQuery, useMutation } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp.js';

export const useValorantIntegration = () => {
	const queryClient = useQueryClient();
	const queryKey = ['valorantIntegration'];

	return {
		useGetData: () => useQuery({
			queryKey,
			queryFn: async () => {
				const call = await protectedApiClient.integrationsValorantGet({});
				return call.response;
			},
		}),
		usePost: () => useMutation({
			mutationKey: ['valorantIntegration/post'],
			mutationFn: async (userName: string) => {
				await protectedApiClient.integrationsValorantUpdate({ userName });
				await queryClient.invalidateQueries(queryKey);
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(queryKey);
			},
		}),
	};
};
