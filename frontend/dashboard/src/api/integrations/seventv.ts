import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type {
	UpdateDataRequest,
} from '@twir/api/messages/integrations_seventv/integrations_seventv';

import { protectedApiClient } from '@/api/twirp';

export const useSevenTvIntegration = () => {
	const queryClient = useQueryClient();

	return {
		useData: () => useQuery({
			queryKey: ['sevenTvIntegration'],
			refetchInterval: 5000,
			queryFn: async () => {
				const request = await protectedApiClient.integrationsSevenTvGetData({});
				return request.response;
			},
		}),
		useUpdate: () => useMutation({
			mutationKey: ['sevenTvIntegration'],
			mutationFn: async (data: UpdateDataRequest) => {
				const request = await protectedApiClient.integrationsSevenTvUpdate(data);
				return request.response;
			},
			async onSuccess() {
				await queryClient.invalidateQueries(['sevenTvIntegration']);
			},
		}),
	};
};
