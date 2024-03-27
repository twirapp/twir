import { useMutation } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp';

export const useNightbotIntegrationImporter = () => {
	return {
		useCommandsImporter: () => useMutation({
			mutationKey: ['integrationsNightbotImportCommands'],
			mutationFn: async () => {
				const call = await protectedApiClient.integrationsNightbotImportCommands({});
				return call.response;
			},
	}),
	};
};
