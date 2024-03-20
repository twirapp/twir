import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp.js';

export const usePubgIntegration = () => {
	const queryClient = useQueryClient();
	const queryKey = ['pubg'];

	return {
		useGetData: () => useQuery({
			queryKey,
			queryFn: async () => {
				const call = await protectedApiClient.integrationsPubgGet({});
				return call.response.nickname;
			},
		}),
		usePut: () => useMutation({
			mutationKey: ['pubg/put'],
			mutationFn: async (nickname: string) => {
				await protectedApiClient.integrationsPubgPut({ nickname });
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(queryKey);
			},
		}),
	};
};
