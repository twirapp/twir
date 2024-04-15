import { useMutation } from '@tanstack/vue-query';
import type { UpdateUserRequest } from '@twir/api/messages/users/users';

import { protectedApiClient } from '@/api/twirp';

export const useUser = () => {
	return {
		useRegenerateApiKey: () => useMutation({
			mutationKey: ['userRegenerateApiKey'],
			mutationFn: async () => {
				const call = await protectedApiClient.usersRegenerateApiKey({});
				return call.response;
			},
			async onSuccess() {

			},
		}),
		useUpdate: () => useMutation({
			mutationKey: ['userUpdate'],
			mutationFn: async (data: UpdateUserRequest) => {
				const call = await protectedApiClient.usersUpdate(data);
				return call.response;
			},
			async onSuccess() {

			},
		}),
	};
};
