import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type { UpdateDuelSettings } from '@twir/api/messages/games/games';

import { protectedApiClient } from '@/api/twirp';

const key = ['duelGame'];

export const useDuelGame = () => {
	const queryClient = useQueryClient();

	return {
		useSettings: () => useQuery({
			queryKey: key,
			queryFn: async () => {
				const req = await protectedApiClient.gamesGetDuelSettings({});
				return req.response;
			},
		}),
		useUpdate: () => useMutation({
			mutationKey: key,
			mutationFn: async (opts: UpdateDuelSettings) => {
				const req = await protectedApiClient.gamesUpdateDuelSettings(opts);
				return req.response;
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(key);
				await queryClient.invalidateQueries(['commands']);
			},
		}),
	};
};
