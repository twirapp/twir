import { useQuery } from '@tanstack/vue-query';

import { protectedApiClient } from './twirp';

export const useModerationAvailableLanguages = () => useQuery({
	queryKey: ['moderationAvailableLanguages'],
	queryFn: async () => {
		const request = await protectedApiClient.moderationAvailableLanguages({});
		return request.response;
	},
});
