import { useMutation, useQueryClient } from '@tanstack/vue-query';
import type { UpdateUserRequest } from '@twir/grpc/generated/api/api/users';

import { profileQueryOptions } from '@/api/auth';
import { protectedApiClient } from '@/api/twirp';

export const useUser = () => {
	const queryClient = useQueryClient();

	return {
		useRegenerateApiKey: () => useMutation({
			mutationKey: ['userRegenerateApiKey'],
			mutationFn: async () => {
				const call = await protectedApiClient.usersRegenerateApiKey({});
				return call.response;
			},
			async onSuccess() {
				await queryClient.invalidateQueries(profileQueryOptions.queryKey);
			},
		}),
		useUpdate: () => useMutation({
			mutationKey: ['userUpdate'],
			mutationFn: async (data: UpdateUserRequest) => {
				const call = await protectedApiClient.usersUpdate(data);
				return call.response;
			},
			async onSuccess() {
				await queryClient.invalidateQueries(profileQueryOptions.queryKey);
			},
		}),
	};
};
