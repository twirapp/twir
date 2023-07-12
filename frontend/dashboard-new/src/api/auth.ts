import { useQuery } from '@tanstack/vue-query';
import { Profile } from '@twir/grpc/generated/api/api/auth';

import { protectedApiClient } from './twirp.js';

export const useProfile = () =>
	useQuery<Profile>({
		queryKey: [`/api/auth/profile`],
		queryFn: async () => {
			const call = await protectedApiClient.authUserProfile({});
			return call.response;
		},
		retry: false,
	});
