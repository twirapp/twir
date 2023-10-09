import { useMutation, useQuery } from '@tanstack/vue-query';
import type { UpdateMessage } from '@twir/grpc/generated/api/api/integrations_discord';

import { protectedApiClient } from '@/api/twirp';

export const useDiscordIntegration = () => {
	return {
		getData: () => useQuery({
			queryKey: ['integrationsDiscordGetData'],
			queryFn: async () => {
				const call = await protectedApiClient.integrationsDiscordGetData({});
				return call.response;
			},
		}),
		updater: useMutation({
			mutationKey: ['integrationsDiscordUpdater'],
			mutationFn: async (data: UpdateMessage) => {
				const call = await protectedApiClient.integrationsDiscordUpdate(data);
				return call.response;
			},
		}),
	};
};
