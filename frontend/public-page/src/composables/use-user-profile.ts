import { useQuery } from '@tanstack/vue-query';

import { protectedClient } from '@/api/twirp';

export const useUserProfile = () => useQuery({
	queryKey: ['userProfile'],
	queryFn: async () => {
		const request = await protectedClient.authUserProfile({});
		return request.response;
	},
});
