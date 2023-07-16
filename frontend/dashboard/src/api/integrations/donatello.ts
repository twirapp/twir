import { useQuery } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp.js';

export const useDonatelloIntegration = () => useQuery({
	queryKey: ['donatello'],
	queryFn: async () => {
		const call = await protectedApiClient.integrationsDonatelloGet({});
		return call.response;
	},
});
