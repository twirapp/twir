import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp.js';

export const useYoutubeModuleSettings = () => {
	const queryClient = useQueryClient();
	return {
		getAll: () => useQuery({
			queryKey: ['youtubeModuleSettings'],
			queryFn: async () => {
				const call = await protectedApiClient.modulesSRGet({});
				return call.response;
			},
		}),
		update: () => useMutation({
			mutationKey: ['youtubeModuleSettingsMutation'],
			mutationFn: async (settings: Parameters<typeof protectedApiClient.modulesSRUpdate>[0]) => {
				await protectedApiClient.modulesSRUpdate(settings);
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(['youtubeModuleSettings']);
			},
		}),
	};
};
