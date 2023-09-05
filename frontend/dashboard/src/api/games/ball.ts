import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp';

export const use8ballSettings = () => useQuery({
	queryKey: ['8ballSettings'],
	queryFn: async () => {
		const req = await protectedApiClient.gamesGetEightBallSettings({});
		return req.response;
	},
});

export const use8ballUpdateSettings = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['8ballSettings'],
		mutationFn: async (opts: { answers: string[], enabled: boolean }) => {
			const req = await protectedApiClient.gamesUpdateEightBallSettings(opts);
			return req.response;
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries(['8ballSettings']);
		},
	});
};
