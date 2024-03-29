import { useMutation, useQueryClient } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp';

export const useNightbotIntegrationImporter = () => {
	const queryClient = useQueryClient();

	return {
		useCommandsImporter: () => useMutation({
			mutationKey: ['integrationsNightbotImportCommands'],
			mutationFn: async () => {
				const call = await protectedApiClient.integrationsNightbotImportCommands({});
				return call.response;
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(['commands']);
			},
		}),
		useTimersImporter: () => useMutation({
			mutationKey: ['integrationsNightbotImportTimers'],
			mutationFn: async () => {
				const call = await protectedApiClient.integrationsNightbotImportTimers({});
				return call.response;
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(['timers']);
			},
		}),
	};
};
