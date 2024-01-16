import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type { UpdateRussianRouletteSettings } from '@twir/api/messages/games/games';

import { protectedApiClient } from '@/api/twirp';

const key = ['russianRoulette'];

export const useRussianRouletteSettings = () => useQuery({
	queryKey: key,
	queryFn: async () => {
		const req = await protectedApiClient.gamesGetRouletteSettings({});
		return req.response;
	},
});

export const useRussianRouletteUpdateSettings = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: key,
		mutationFn: async (opts: UpdateRussianRouletteSettings) => {
			const req = await protectedApiClient.gamesUpdateRouletteSettings(opts);
			return req.response;
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries(key);
			await queryClient.invalidateQueries(['commands']);
		},
	});
};
